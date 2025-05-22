package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Student struct {
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

// temp data
type StudentList map[int]Student

var students = StudentList{
	0: Student{"Alice", 100},
	1: Student{"Bob", 95},
}

func main() {
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

func fetchStudents() StudentList {
	return students
}

func addStudent(newStudent Student) Student {
	log.Printf("Adding new student, name: '%v', grade: '%v'\n", newStudent.Name, newStudent.Grade)

	students[len(students)+1] = newStudent

	return newStudent
}

func updateStudent(ID int, updateStudent Student) Student {
	existingStudent := students[ID]
	log.Printf("Updating student: %v, name: '%v', grade: '%v'\n", ID, existingStudent.Name, existingStudent.Grade)

	existingStudent.Name = updateStudent.Name
	existingStudent.Grade = updateStudent.Grade

	students[ID] = existingStudent

	log.Printf("Updated student: %v, name: '%v', grade: '%v'\n", ID, existingStudent.Name, existingStudent.Grade)

	return existingStudent
}

func deleteStudent(ID int) Student {
	log.Printf("Deleting student: %v", ID)

	deletedStudent := students[ID]
	delete(students, ID)

	log.Printf("Deleted student: %v, name: '%v', grade: '%v'\n", ID, deletedStudent.Name, deletedStudent.Grade)

	return deletedStudent
}
