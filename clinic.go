package clinic

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Appointment struct {
	ApptID      int
	PatientName string
	DoctorName  string
	Date        string
	Status      string
}

var csvHeader = []string{"ApptID", "PatientName", "DoctorName", "Date", "Status"}

type Store struct {
	filePath     string
	appointments []Appointment
}

func NewStore(filePath string) (*Store, error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, errors.New("appointment storage path cannot be empty")
	}

	store := &Store{filePath: filePath}
	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return store, nil
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return store, nil
		}
		return nil, err
	}
	if !equalFields(header, csvHeader) {
		return nil, errors.New("invalid appointment CSV header")
	}
	reader.FieldsPerRecord = len(csvHeader)

	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, errors.New("invalid appointment ID in CSV file")
		}
		store.appointments = append(store.appointments, Appointment{
			ApptID:      id,
			PatientName: record[1],
			DoctorName:  record[2],
			Date:        record[3],
			Status:      record[4],
		})
	}
	return store, nil
}

func (s *Store) BookAppointment(a Appointment) error {
	a.PatientName = strings.TrimSpace(a.PatientName)
	a.DoctorName = strings.TrimSpace(a.DoctorName)
	a.Date = strings.TrimSpace(a.Date)
	if a.PatientName == "" || a.DoctorName == "" {
		return errors.New("patient name and doctor name are required")
	}
	if err := validateDate(a.Date); err != nil {
		return err
	}
	if a.ApptID < 0 {
		return errors.New("appointment ID cannot be negative")
	}

	maxID := 0
	for _, existing := range s.appointments {
		if existing.ApptID == a.ApptID && a.ApptID != 0 {
			return errors.New("appointment ID already exists")
		}
		if existing.ApptID > maxID {
			maxID = existing.ApptID
		}
	}

	if a.ApptID == 0 {
		if maxID == int(^uint(0)>>1) {
			return errors.New("cannot generate appointment ID: ID range exhausted")
		}
		a.ApptID = maxID + 1
	}
	if a.Status == "" {
		a.Status = "Scheduled"
	}

	s.appointments = append(s.appointments, a)
	if err := s.save(); err != nil {
		s.appointments = s.appointments[:len(s.appointments)-1]
		return err
	}
	return nil
}

func (s *Store) GetAppointment(id int) (*Appointment, error) {
	for i := range s.appointments {
		if s.appointments[i].ApptID == id {
			appointment := s.appointments[i]
			return &appointment, nil
		}
	}
	return nil, errors.New("appointment not found")
}

func (s *Store) RescheduleAppointment(id int, newDate string) error {
	newDate = strings.TrimSpace(newDate)
	if err := validateDate(newDate); err != nil {
		return err
	}

	for i := range s.appointments {
		if s.appointments[i].ApptID == id {
			oldDate := s.appointments[i].Date
			s.appointments[i].Date = newDate
			if err := s.save(); err != nil {
				s.appointments[i].Date = oldDate
				return err
			}
			return nil
		}
	}
	return errors.New("appointment not found")
}

func (s *Store) CancelAppointment(id int) error {
	for i := range s.appointments {
		if s.appointments[i].ApptID == id {
			oldStatus := s.appointments[i].Status
			s.appointments[i].Status = "Canceled"
			if err := s.save(); err != nil {
				s.appointments[i].Status = oldStatus
				return err
			}
			return nil
		}
	}
	return errors.New("appointment not found")
}

func (s *Store) ListAppointments() []Appointment {
	return append([]Appointment(nil), s.appointments...)
}

func validateDate(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return errors.New("date must use YYYY-MM-DD")
	}
	return nil
}

func equalFields(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func (s *Store) save() error {
	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	if err := writer.Write(csvHeader); err != nil {
		file.Close()
		return err
	}
	for _, appointment := range s.appointments {
		record := []string{
			strconv.Itoa(appointment.ApptID),
			appointment.PatientName,
			appointment.DoctorName,
			appointment.Date,
			appointment.Status,
		}
		if err := writer.Write(record); err != nil {
			file.Close()
			return err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}
