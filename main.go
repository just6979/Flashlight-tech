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

	mux.HandleFunc("/{$}", listHandler)
	mux.HandleFunc("GET /students", listHandler)
	mux.HandleFunc("POST /students", addHandler)
	mux.HandleFunc("PUT /students/{id}", updateHandler)
	mux.HandleFunc("DELETE /students/{id}", deleteHandler)

	log.Printf("Starting server at http://localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%v", port), mux))
	log.Println("Exiting server")
}

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
	panic("unimplemented")
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("PUT request for %v\n", r.URL)
	panic("unimplemented")
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("DELETE request for %v\n", r.URL)
	panic("unimplemented")
}

func fetchStudents() []Student {
	var students []Student
	students = append(students, Student{1, "Alice", 100})
	students = append(students, Student{2, "Bob", 95})
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
