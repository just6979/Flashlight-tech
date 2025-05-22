package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

// temp data
var students []Student

func main() {
	// temp data
	students = append(students, Student{1, "Alice", 100})
	students = append(students, Student{2, "Bob", 95})

	port := 8080
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", listHandler)
	mux.HandleFunc("GET /students", listHandler)
	mux.HandleFunc("POST /students", addHandler)
	mux.HandleFunc("PUT /students/{id}", updateHandler)
	mux.HandleFunc("DELETE /students/{id}", deleteHandler)

	log.Printf("Starting server at http://localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%v", port), mux))
	log.Println("Exiting server")
}

// request handlers

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET request for %v\n", r.URL)

	students := fetchStudents()
	log.Printf("Students List: %+v\n", students)

	jsonResponse, err := json.Marshal(students)
	if err != nil {
		log.Println("Unable to marshall to JSON")
		jsonResponse = []byte("{}")
	}
	log.Printf("Response: %v\n", string(jsonResponse))
	w.Write(jsonResponse)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST request for %v\n", r.URL)

	var newStudent Student
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newStudent)
	if err != nil {
		panic(err)
	}
	addStudent(newStudent)

	responseData := map[string]string{"added": fmt.Sprint(newStudent)}
	jsonResponse, _ := json.Marshal(responseData)
	log.Printf("Response: %v\n", string(jsonResponse))
	w.Write(jsonResponse)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST request for %v\n", r.URL)

	ID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	var updatedStudent Student
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&updatedStudent)
	if err != nil {
		panic(err)
	}
	updateStudent(ID, updatedStudent)

	responseData := map[string]string{"updated": fmt.Sprint(updatedStudent)}
	jsonResponse, _ := json.Marshal(responseData)
	log.Printf("Response: %v\n", string(jsonResponse))
	w.Write(jsonResponse)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("DELETE request for %v\n", r.URL)

	ID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	oldStudent := deleteStudent(ID)

	responseData := map[string]string{"deleted": fmt.Sprint(oldStudent)}
	jsonResponse, _ := json.Marshal(responseData)
	log.Printf("Response: %v\n", string(jsonResponse))
	w.Write(jsonResponse)
}

// data handlers

func fetchStudents() []Student {
	return students
}

func addStudent(newStudent Student) Student {
	log.Printf("Adding new student, name: '%v', grade: '%v'\n", newStudent.Name, newStudent.Grade)
	students = append(students, newStudent)
	return newStudent
}

func updateStudent(ID int, updateStudent Student) Student {
	log.Printf("Updating student: %v, name: '%v', grade: '%v'\n", ID, updateStudent.Name, updateStudent.Grade)
	var existingIndex int
	for i := 0; i < len(students); i++ {
		if students[i].ID == ID {
			existingIndex = i
			break
		}
	}

	existingStudent := students[existingIndex]
	existingStudent.Name = updateStudent.Name
	existingStudent.Grade = updateStudent.Grade

	students = append(students[:existingIndex], existingStudent, students[existingIndex+1])

	log.Printf("Updated student: %v, name: '%v', grade: '%v'\n", ID, existingStudent.Name, existingStudent.Grade)
	return existingStudent
}

func deleteStudent(ID int) Student {
	var deletedIndex int
	var deletedStudent Student
	for i := 0; i < len(students); i++ {
		if students[i].ID == ID {
			deletedIndex = i
			break
		}
	}
	deletedStudent = students[deletedIndex]

	students = append(students[:deletedIndex], students[deletedIndex+1])

	log.Printf("Deleting student: %v, name: '%v', grade: '%v'\n", ID, deletedStudent.Name, deletedStudent.Grade)
	return deletedStudent
}

/*
Back-end: Implement a RESTful API with CRUD operations for managing students.
    API
        GET /students - Retrieve a list of all students.
        POST /students - Add a new student.
        PUT /students/:id - Update an existing student.
        DELETE /students/:id - Delete a student.
    Data Model
        id (integer, primary key)
        name (string)
        grade (integer)
Front-end: Create a functional UI that interacts with the backend API. It does not need to be polished or aesthetically pleasing; it just needs to be able to connect with the back-end and perform the following functions:
    View a list of all students.
    Add a new student.
    Update an existing student.
    Delete a student.
*/
