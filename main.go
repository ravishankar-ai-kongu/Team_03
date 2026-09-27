package main

import (
	"fmt"
	"go1credit/pkg/students"
)

func main() {
	studentList := students.StudentList{}
	for {
		fmt.Println("\n===== STUDENT MANAGEMENT SYSTEM =====")
		fmt.Println("1. Create Student")
		fmt.Println("2. Update Student")
		fmt.Println("3. Search Student")
		fmt.Println("4. Delete Student")
		fmt.Println("5. Display All Students")
		fmt.Println("6. Exit")

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {

		case 1:
			var rollNo string
			var name string
			var phoneNo int

			fmt.Print("Enter Roll No: ")
			fmt.Scan(&rollNo)

			fmt.Print("Enter Name: ")
			fmt.Scan(&name)

			fmt.Print("Enter Phone No: ")
			fmt.Scan(&phoneNo)

			student := students.NewStudent(name, rollNo, phoneNo)
			studentList.Create(student)

			fmt.Println("Student created successfully.")

		case 2:
			var rollNo string
			var name string

			fmt.Print("Enter Roll No: ")
			fmt.Scan(&rollNo)

			fmt.Print("Enter New Name: ")
			fmt.Scan(&name)

			if studentList.Update(rollNo, name) {
				fmt.Println("Student updated successfully.")
			} else {
				fmt.Println("Student not found.")
			}

		case 3:
			var rollNo string

			fmt.Print("Enter Roll No: ")
			fmt.Scan(&rollNo)

			student, found := studentList.Search(rollNo)

			if found {
				student.Display()
			} else {
				fmt.Println("Student not found.")
			}

		case 4:
			var rollNo string

			fmt.Print("Enter Roll No: ")
			fmt.Scan(&rollNo)

			if studentList.Delete(rollNo) {
				fmt.Println("Student deleted successfully.")
			} else {
				fmt.Println("Student not found.")
			}

		case 5:
			studentList.DisplayAll()

		case 6:
			fmt.Println("Thank you!")
			return

		default:
			fmt.Println("Invalid choice.")
		}
	}
}
