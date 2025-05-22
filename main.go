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

	fmt.Printf("Starting server at http://localhost:%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%v", port), mux))
	fmt.Println("Exiting server")
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	switch method {
	case "GET":
		fmt.Printf("Handling GET request for %v\n", r.URL)
		listStudents(w)
	case "POST":
		println("Handling POST request")
		addStudent(w, r)
	case "PUT":
		println("Handling PUT request")
		updateStudent(w, r)
	case "DELETE":
		println("Handling DELETE request")
		deleteStudent(w, r)
	default:
		println("Handling unknown request")
		w.Write([]byte(fmt.Sprintf("{error: {msg: 'Unhandled HTTP request', value: '%v'}}", method)))
	}
}

func listStudents(w http.ResponseWriter) {

	students := fetchStudents()
	fmt.Printf("Students List: %+v\n", students)

	jsonResponse, err := json.Marshal(students)
	if err != nil {
		println("Unable to marshall to JSON")
		response := "{}"
		println(response)
		w.Write([]byte(response))
		return
	}
	fmt.Printf("Response: %v\n", string(jsonResponse))
	w.Write(jsonResponse)
}

func addStudent(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
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
