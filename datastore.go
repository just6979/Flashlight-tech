package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func dbConnect() {
	var err error
	dbConnStr := "postgres://postgres:@localhost/postgres?sslmode=disable"
	log.Printf("Connecting to database %v", dbConnStr)
	db, err = sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
}

func dbClose() {
	db.Close()
}

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

	log.Println("Fetching all students")

	const selectStatement = "SELECT id, name, grade FROM students;"
	rows, err := db.Query(selectStatement)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var nextStudent Student
		err := rows.Scan(&nextStudent.ID, &nextStudent.Name, &nextStudent.Grade)
		if err != nil {
			log.Printf("Unable to map Row to Student: %v", err)
			continue
		}
		log.Printf("id %d, name is %s, grade is %d",
			nextStudent.ID, nextStudent.Name, nextStudent.Grade)
		students = append(students, nextStudent)
	}

	return students
}

func addStudent(newStudent Student) Student {
	log.Printf("Adding new student, name: '%v', grade: '%v'",
		newStudent.Name, newStudent.Grade)

	const addStatement = `INSERT INTO students (name, grade) VALUES ($1, $2) RETURNING id`
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

	log.Printf("Deleted student: id %v, name '%v', grade '%v'",
		deletedStudent.ID, deletedStudent.Name, deletedStudent.Grade)

	return deletedStudent
}
