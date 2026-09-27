
package students

import "fmt"

type Student struct {
	RoleNo  string
	Name    string
	phoneno int
}

func (self *Student) UpdateName(name string) {
	self.Name = name
}
func (s Student) Display() {
	fmt.Println(s.Name, s.RoleNo, s.phoneno)
}
func NewStudent(name string, roleno string, phoneno int) Student {
	return Student{
		RoleNo:  roleno,
		Name:    name,
		phoneno: phoneno,
	}
}

type StudentList struct {
	Students []Student
}

func (s *StudentList) Create(student Student) {
	s.Students = append(s.Students, student)
}
func (s StudentList) Search(rollNo string) (Student, bool) {
	for _, student := range s.Students {
		if student.RoleNo == rollNo {
			return student, true
		}
	}

	return Student{}, false
}
func (s *StudentList) Update(rollNo string, name string) bool {
	for i := 0; i < len(s.Students); i++ {
		if s.Students[i].RoleNo == rollNo {
			s.Students[i].UpdateName(name)
			return true
		}
	}

	return false
}
func (s *StudentList) Delete(rollNo string) bool {
	for i := 0; i < len(s.Students); i++ {
		if s.Students[i].RoleNo == rollNo {
			s.Students = append(s.Students[:i], s.Students[i+1:]...)
			return true
		}
	}

	return false
}
func (s StudentList) DisplayAll() {
	if len(s.Students) == 0 {
		fmt.Println("No students available.")
		return
	}

	fmt.Println("\n===== ALL STUDENTS =====")

	for _, student := range s.Students {
		student.Display()
		fmt.Println("--------------------")
	}
}
