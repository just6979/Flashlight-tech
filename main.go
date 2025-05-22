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

func main() {
	dbConnect()

	host := "localhost:8080"
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", listHandler)
	mux.HandleFunc("GET /students", listHandler)
	mux.HandleFunc("POST /students", addHandler)
	mux.HandleFunc("PUT /students/{id}", updateHandler)
	mux.HandleFunc("DELETE /students/{id}", deleteHandler)

	log.Printf("Starting server at http://%v", host)
	log.Fatal(http.ListenAndServe(host, mux))

	dbClose()
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET request for %v", r.URL)

	students := fetchAllStudents()
	log.Printf("Students List: %+v", students)

	jsonResponse, err := json.Marshal(students)
	if err != nil {
		log.Println("Unable to marshall to JSON")
		jsonResponse = []byte("{}")
	}
	log.Printf("Response: %v", string(jsonResponse))
	w.Write(jsonResponse)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST request for %v", r.URL)

	var newStudent Student
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newStudent)
	if err != nil {
		panic(err)
	}
	newStudent = addStudent(newStudent)

	responseData := map[string]string{"added": fmt.Sprint(newStudent)}
	jsonResponse, _ := json.Marshal(responseData)
	log.Printf("Response: %v", string(jsonResponse))
	w.Write(jsonResponse)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST request for %v", r.URL)

	ID, _ := strconv.Atoi(r.PathValue("id"))

	var updatedStudent Student
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&updatedStudent)
	if err != nil {
		panic(err)
	}
	updatedStudent = updateStudent(ID, updatedStudent)

	responseData := map[string]string{"updated": fmt.Sprint(updatedStudent)}
	jsonResponse, _ := json.Marshal(responseData)
	log.Printf("Response: %v", string(jsonResponse))
	w.Write(jsonResponse)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("DELETE request for %v", r.URL)

	ID, _ := strconv.Atoi(r.PathValue("id"))
	oldStudent := deleteStudent(ID)

	responseData := map[string]string{"deleted": fmt.Sprint(oldStudent)}
	jsonResponse, _ := json.Marshal(responseData)
	log.Printf("Response: %v", string(jsonResponse))
	w.Write(jsonResponse)
}
