package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"go1credit/pkg/clinic"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		os.Exit(2)
	}
	if args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		printUsage()
		return
	}

	store, err := clinic.NewStore("appointments.csv")
	if err != nil {
		fail(err)
	}

	switch args[0] {
	case "book":
		if len(args) < 2 {
			fail(errors.New(`book requires: "Patient Name,Doctor Name,YYYY-MM-DD"`))
		}
		fields := strings.Split(strings.Join(args[1:], " "), ",")
		if len(fields) != 3 {
			fail(errors.New("book requires exactly three comma-separated values: Patient Name,Doctor Name,YYYY-MM-DD"))
		}
		if err := store.BookAppointment(clinic.Appointment{
			PatientName: strings.TrimSpace(fields[0]),
			DoctorName:  strings.TrimSpace(fields[1]),
			Date:        strings.TrimSpace(fields[2]),
		}); err != nil {
			fail(err)
		}
		appointments := store.ListAppointments()
		printAppointment(appointments[len(appointments)-1])

	case "list":
		if len(args) != 1 {
			fail(errors.New("list does not accept additional arguments"))
		}
		appointmentsChannel := make(chan []clinic.Appointment)
		go func() {
			appointmentsChannel <- store.ListAppointments()
		}()
		appointments := <-appointmentsChannel
		if len(appointments) == 0 {
			fmt.Println("No appointments found.")
			return
		}
		for _, appointment := range appointments {
			printAppointment(appointment)
		}

	case "get":
		id := requireID(args[1:])
		appointment, err := store.GetAppointment(id)
		if err != nil {
			fail(err)
		}
		printAppointment(*appointment)

	case "reschedule":
		var id int
		var newDate string
		if len(args) == 1 {
			id, newDate = promptReschedule(store, os.Stdin, os.Stdout)
		} else {
			if len(args) != 3 {
				fail(errors.New("reschedule requires: [ID YYYY-MM-DD]"))
			}
			id = parseID(args[1])
			newDate = args[2]
		}
		if err := store.RescheduleAppointment(id, newDate); err != nil {
			fail(err)
		}
		fmt.Println("Appointment", id, "rescheduled to", newDate)

	case "cancel":
		id := requireID(args[1:])
		if err := store.CancelAppointment(id); err != nil {
			fail(err)
		}
		fmt.Println("Appointment", id, "canceled.")

	default:
		printUsage()
		fail(errors.New("unknown command"))
	}
}

func requireID(args []string) int {
	if len(args) != 1 {
		fail(errors.New("command requires exactly one appointment ID"))
	}
	return parseID(args[0])
}

func promptReschedule(store *clinic.Store, input io.Reader, output io.Writer) (int, string) {
	reader := bufio.NewReader(input)
	fmt.Fprint(output, "Enter appointment ID: ")
	idInput, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		fail(err)
	}
	if strings.TrimSpace(idInput) == "" {
		fail(errors.New("appointment ID is required"))
	}
	id := parseID(strings.TrimSpace(idInput))

	appointment, err := store.GetAppointment(id)
	if err != nil {
		fail(err)
	}
	fmt.Fprintln(output, "Appointment to edit:")
	printAppointmentTo(output, *appointment)

	fmt.Fprint(output, "Enter new date (YYYY-MM-DD): ")
	dateInput, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		fail(err)
	}
	return id, strings.TrimSpace(dateInput)
}

func parseID(value string) int {
	id, err := strconv.Atoi(value)
	if err != nil || id < 1 {
		fail(errors.New("appointment ID must be a positive integer"))
	}
	return id
}

func printAppointment(a clinic.Appointment) {
	printAppointmentTo(os.Stdout, a)
}

func printAppointmentTo(output io.Writer, a clinic.Appointment) {
	fmt.Fprintln(output, "ID:", a.ApptID, "| Patient:", a.PatientName, "| Doctor:", a.DoctorName, "| Date:", a.Date, "| Status:", a.Status)
}

func printUsage() {
	fmt.Println(`Hospital Outpatient Appointment Tracker
	Usage:
	  go run . book "Patient Name,Doctor Name,YYYY-MM-DD"
	  go run . list
	  go run . get ID
	  go run . reschedule [ID YYYY-MM-DD]
	  go run . cancel ID

	Run "go run . reschedule" to enter an appointment ID and new date interactively.
	Appointments are saved in appointments.csv in the current directory.`)
}

func fail(err error) {
	fmt.Println("Error:", err)
	os.Exit(1)
}
