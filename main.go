package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Student struct {
	Id    int
	Name  string
	Grade int
}

func main() {
	port := 8080
	mux := http.NewServeMux()

	mux.HandleFunc("/", baseHandler)
	//mux.HandleFunc("/students", baseHandler)

	log.Printf("Starting server at http://localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%v", port), mux))
	log.Println("Exiting server")
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	var responseBytes []byte

	switch method {
	case "GET":
		log.Printf("GET request for %v\n", r.URL)
		responseBytes = listStudents()
	case "POST":
		log.Printf("POST request for %v\n", r.URL)
		responseBytes = addStudent(r)
	case "PUT":
		log.Printf("PUT request for %v\n", r.URL)
		responseBytes = updateStudent(r)
	case "DELETE":
		log.Printf("DELETE request for %v\n", r.URL)
		responseBytes = deleteStudent(r)
	default:
		log.Printf("Unhandled request for %v\n", r.URL)
		responseBytes = fmt.Appendf(nil, "{error: {msg: 'Unhandled HTTP request', value: '%v'}}", method)
	}

	w.Write(responseBytes)
}

func listStudents() []byte {
	students := fetchStudents()
	fmt.Printf("Students List: %+v\n", students)

	jsonResponse, err := json.Marshal(students)
	if err != nil {
		log.Println("Unable to marshall to JSON")
		jsonResponse = []byte("{}")
	}
	log.Printf("Response: %v\n", string(jsonResponse))
	return jsonResponse
}

func addStudent(r *http.Request) []byte {
	panic("unimplemented")
}

func updateStudent(r *http.Request) []byte {
	panic("unimplemented")
}

func deleteStudent(r *http.Request) []byte {
	panic("unimplemented")
}

func fetchStudents() any {
	var students []Student
	newStudent := Student{1, "Alice", 100}
	students = append(students, newStudent)
	return students
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
