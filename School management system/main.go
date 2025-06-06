package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./school.db")
	if err != nil {
		log.Fatal(err)
	}
	createTable := `CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT,
		last_name TEXT,
		age INTEGER
	);
	CREATE TABLE IF NOT EXISTS teachers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT,
		last_name TEXT,
		subject TEXT
	);
	CREATE TABLE IF NOT EXISTS courses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		teacher_id INTEGER,
		FOREIGN KEY(teacher_id) REFERENCES teachers(id)
	);
	CREATE TABLE IF NOT EXISTS enrollments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		student_id INTEGER,
		course_id INTEGER,
		grade TEXT,
		FOREIGN KEY(student_id) REFERENCES students(id),
		FOREIGN KEY(course_id) REFERENCES courses(id)
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB() // Initialize SQLite DB
	http.HandleFunc("/students", studentsHandler)
	http.HandleFunc("/teachers", teachersHandler)
	http.HandleFunc("/courses", coursesHandler)
	http.HandleFunc("/enroll", enrollHandler)
	http.HandleFunc("/grade", gradeHandler)
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

// --- HTTP Handlers (CRUD + Enrollment/Grading) ---
// Each handler will support GET, POST, PUT, DELETE as appropriate

// Example: Students CRUD
func studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT id, first_name, last_name, age FROM students")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()
		var students []map[string]interface{}
		for rows.Next() {
			var id, age int
			var firstName, lastName string
			rows.Scan(&id, &firstName, &lastName, &age)
			students = append(students, map[string]interface{}{
				"id": id, "first_name": firstName, "last_name": lastName, "age": age,
			})
		}
		json.NewEncoder(w).Encode(students)
	case "POST":
		var s struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Age       int    `json:"age"`
		}
		json.NewDecoder(r.Body).Decode(&s)
		res, err := db.Exec("INSERT INTO students (first_name, last_name, age) VALUES (?, ?, ?)", s.FirstName, s.LastName, s.Age)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		id, _ := res.LastInsertId()
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	case "PUT":
		var s struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Age       int    `json:"age"`
		}
		json.NewDecoder(r.Body).Decode(&s)
		_, err := db.Exec("UPDATE students SET first_name=?, last_name=?, age=? WHERE id=?", s.FirstName, s.LastName, s.Age, s.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(204)
	case "DELETE":
		id := r.URL.Query().Get("id")
		_, err := db.Exec("DELETE FROM students WHERE id=?", id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(204)
	}
}

// Teachers CRUD
func teachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT id, first_name, last_name, subject FROM teachers")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()
		var teachers []map[string]interface{}
		for rows.Next() {
			var id int
			var firstName, lastName, subject string
			rows.Scan(&id, &firstName, &lastName, &subject)
			teachers = append(teachers, map[string]interface{}{
				"id": id, "first_name": firstName, "last_name": lastName, "subject": subject,
			})
		}
		json.NewEncoder(w).Encode(teachers)
	case "POST":
		var t struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Subject   string `json:"subject"`
		}
		json.NewDecoder(r.Body).Decode(&t)
		res, err := db.Exec("INSERT INTO teachers (first_name, last_name, subject) VALUES (?, ?, ?)", t.FirstName, t.LastName, t.Subject)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		id, _ := res.LastInsertId()
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	case "PUT":
		var t struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Subject   string `json:"subject"`
		}
		json.NewDecoder(r.Body).Decode(&t)
		_, err := db.Exec("UPDATE teachers SET first_name=?, last_name=?, subject=? WHERE id=?", t.FirstName, t.LastName, t.Subject, t.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(204)
	case "DELETE":
		id := r.URL.Query().Get("id")
		_, err := db.Exec("DELETE FROM teachers WHERE id=?", id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(204)
	}
}

// Courses CRUD
func coursesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT id, name, code FROM courses")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()
		var courses []map[string]interface{}
		for rows.Next() {
			var id int
			var name, code string
			rows.Scan(&id, &name, &code)
			courses = append(courses, map[string]interface{}{
				"id": id, "name": name, "code": code,
			})
		}
		json.NewEncoder(w).Encode(courses)
	case "POST":
		var c struct {
			Name string `json:"name"`
			Code string `json:"code"`
		}
		json.NewDecoder(r.Body).Decode(&c)
		res, err := db.Exec("INSERT INTO courses (name, code) VALUES (?, ?)", c.Name, c.Code)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		id, _ := res.LastInsertId()
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	case "PUT":
		var c struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		}
		json.NewDecoder(r.Body).Decode(&c)
		_, err := db.Exec("UPDATE courses SET name=?, code=? WHERE id=?", c.Name, c.Code, c.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(204)
	case "DELETE":
		id := r.URL.Query().Get("id")
		_, err := db.Exec("DELETE FROM courses WHERE id=?", id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(204)
	}
}

// --- Enrollment and Grading ---
func enrollHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	var req struct {
		StudentID int `json:"student_id"`
		CourseID  int `json:"course_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	_, err := db.Exec("INSERT INTO enrollments (student_id, course_id, grade) VALUES (?, ?, '')", req.StudentID, req.CourseID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
}

func gradeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	var req struct {
		EnrollmentID int    `json:"enrollment_id"`
		Grade        string `json:"grade"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	_, err := db.Exec("UPDATE enrollments SET grade=? WHERE id=?", req.Grade, req.EnrollmentID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}
