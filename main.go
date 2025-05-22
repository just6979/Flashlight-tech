package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

var db *sql.DB

func main() {
	var err error
	dbConnStr := "postgres://postgres:@localhost/postgres?sslmode=disable"
	log.Printf("Connecting to database %v", dbConnStr)
	db, err = sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	port := 8080
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", listHandler)
	mux.HandleFunc("GET /students", listHandler)
	mux.HandleFunc("POST /students", addHandler)
	mux.HandleFunc("PUT /students/{id}", updateHandler)
	mux.HandleFunc("DELETE /students/{id}", deleteHandler)

	log.Printf("Starting server at http://localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%v", port), mux))
}

// request handlers

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET request for %v\n", r.URL)

	students := fetchAllStudents()
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
	newStudent = addStudent(newStudent)

	responseData := map[string]string{"added": fmt.Sprint(newStudent)}
	jsonResponse, _ := json.Marshal(responseData)
	log.Printf("Response: %v\n", string(jsonResponse))
	w.Write(jsonResponse)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST request for %v\n", r.URL)

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
	log.Printf("Response: %v\n", string(jsonResponse))
	w.Write(jsonResponse)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("DELETE request for %v\n", r.URL)

	ID, _ := strconv.Atoi(r.PathValue("id"))
	oldStudent := deleteStudent(ID)

	responseData := map[string]string{"deleted": fmt.Sprint(oldStudent)}
	jsonResponse, _ := json.Marshal(responseData)
	log.Printf("Response: %v\n", string(jsonResponse))
	w.Write(jsonResponse)
}

// data handlers

func fetchStudent(ID int) (bool, Student) {
	student := Student{}
	isFound := true

	log.Printf("Fetching student: id %v", ID)

	const findStatement = "SELECT id, name, grade FROM students WHERE id = $1;"
	row := db.QueryRow(findStatement, ID)
	err := row.Scan(&student.ID, &student.Name, &student.Grade)
	if err == sql.ErrNoRows {
		log.Printf("Unable to find Student: id %v", ID)
		isFound = false
	} else if err != nil {
		log.Printf("Unable to map Row to Student: %v", err)
		isFound = false
	}

	log.Printf("Found student: id %v, name: '%v', grade: '%v'",
		student.ID, student.Name, student.Grade)

	return isFound, student
}

func fetchAllStudents() []Student {
	students := []Student{}

	rows, err := db.Query("select * from students")
	if err != nil {
		log.Printf("Query failed: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var nextStudent Student
		err := rows.Scan(&nextStudent.ID, &nextStudent.Name, &nextStudent.Grade)
		if err != nil {
			log.Printf("Unable to map Row to Student: %v\n", err)
			continue
		}
		log.Printf("id %d, name is %s, grade is %d\n", nextStudent.ID, nextStudent.Name, nextStudent.Grade)
		students = append(students, nextStudent)
	}

	return students
}

func addStudent(newStudent Student) Student {
	log.Printf("Adding new student, name: '%v', grade: '%v'\n", newStudent.Name, newStudent.Grade)

	addStatement := `insert into students (name, grade) values ($1, $2) returning id`
	id := 0
	err := db.QueryRow(addStatement, newStudent.Name, newStudent.Grade).Scan(&id)
	if err != nil {
		log.Printf("Error adding Student: %v", err)
	}
	newStudent.ID = id
	return newStudent
}

func updateStudent(ID int, updatedStudent Student) Student {
	log.Printf("Trying to update student: id %v", ID)
	isFound, existingStudent := fetchStudent(ID)
	if !isFound {
		log.Printf("Unable to find student: id %v", ID)
		return Student{}
	}

	updatedStudent.ID = ID
	log.Printf("Existing student: id %v, name: '%v', grade: '%v'",
		existingStudent.ID, existingStudent.Name, existingStudent.Grade)
	log.Printf("Updated student info: name: '%v', grade: '%v'",
		existingStudent.Name, existingStudent.Grade)

	const addStatement = `UPDATE students SET name = $2, grade = $3 WHERE id = $1`
	res, err := db.Exec(addStatement, ID, updatedStudent.Name, updatedStudent.Grade)
	if err != nil {
		log.Printf("Error updating Student: id %v: %v", ID, err)
		return Student{}
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error checking update: %v", err)
		return Student{}
	}
	if count == 0 {
		log.Println("Error updating student, no rows changed.")
		return Student{}
	}

	log.Printf("Updated student: %v, name: '%v', grade: '%v'",
		updatedStudent.ID, updatedStudent.Name, updatedStudent.Grade)

	return updatedStudent
}

func deleteStudent(ID int) Student {
	log.Printf("Trying to delete student: id %v", ID)
	isFound, deletedStudent := fetchStudent(ID)
	if !isFound {
		log.Printf("Unable to find student: id %v", ID)
		return Student{}
	}

	log.Printf("Deleting student: id %v, name %v, grade %v",
		deletedStudent.ID, deletedStudent.Name, deletedStudent.Grade)

	const deleteStatement = `DELETE FROM students WHERE id = $1;`
	res, err := db.Exec(deleteStatement, ID)
	if err != nil {
		log.Printf("Error deleting Student: id %v: %v", ID, err)
		return Student{}
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error checking deletion: %v", err)
		return Student{}
	}
	if count == 0 {
		log.Println("Error deleting student, no rows deleted.")
		return Student{}
	}

	log.Printf("Deleted student: id %v, name '%v', grade '%v'\n",
		deletedStudent.ID, deletedStudent.Name, deletedStudent.Grade)

	return deletedStudent
}
