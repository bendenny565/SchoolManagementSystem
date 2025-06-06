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
	db, err = sql.Open("sqlite3", "school.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Create tables if not exist
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS students (
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
		code TEXT
	);
	CREATE TABLE IF NOT EXISTS enrollments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		student_id INTEGER,
		course_id INTEGER,
		grade TEXT,
		FOREIGN KEY(student_id) REFERENCES students(id),
		FOREIGN KEY(course_id) REFERENCES courses(id)
	);
	`)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
