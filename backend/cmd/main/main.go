package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vincenttanjaya/go_worker_shift/internal/controller"
	"github.com/vincenttanjaya/go_worker_shift/internal/handler"
	"github.com/vincenttanjaya/go_worker_shift/internal/middleware"
	"github.com/vincenttanjaya/go_worker_shift/internal/repository"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Connect to SQLite database
	db, err := sql.Open("sqlite3", "./data/roster.db")
	if err != nil {
		log.Fatal("Could not open DB:", err)
	}
	defer db.Close()

	// Create tables if they do not exist can use migration
	schema := `
	CREATE TABLE IF NOT EXISTS workers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS shifts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		start_time TEXT NOT NULL,
		end_time TEXT NOT NULL,
		role TEXT NOT NULL,
		location TEXT,
		assigned_worker_id INTEGER,
		FOREIGN KEY (assigned_worker_id) REFERENCES workers(id)
	);
	CREATE TABLE IF NOT EXISTS shift_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		worker_id INTEGER NOT NULL,
		shift_id INTEGER NOT NULL,
		status TEXT NOT NULL,
		FOREIGN KEY (worker_id) REFERENCES workers(id),
		FOREIGN KEY (shift_id) REFERENCES shifts(id)
	);`

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to run DB schema: %v", err)
	}

	// Repository init
	workerRepo := &repository.WorkerRepository{DB: db}
	shiftRepo := &repository.ShiftRepository{DB: db}
	shiftRequestRepo := &repository.ShiftRequestRepository{DB: db}

	// Controller init
	workerCtrl := &controller.WorkerController{Repo: workerRepo}
	shiftCtrl := &controller.ShiftController{Repo: shiftRepo}
	shiftRequestCtrl := &controller.ShiftRequestController{
		Repo:       shiftRequestRepo,
		ShiftRepo:  shiftRepo,
		WorkerRepo: workerRepo,
	}

	// Handler init
	workerHandler := &handler.WorkerHandler{Ctrl: workerCtrl}
	shiftHandler := &handler.ShiftHandler{Ctrl: shiftCtrl}
	shiftRequestHandler := &handler.ShiftRequestHandler{Ctrl: shiftRequestCtrl}

	// Router setup
	r := mux.NewRouter()

	handler := middleware.CORS(r)

	// === Workers endpoints ===
	r.HandleFunc("/signup", workerHandler.SignupWorker).Methods("POST")
	r.HandleFunc("/workers", workerHandler.ListWorkers).Methods("GET")
	r.HandleFunc("/workers/{id}", workerHandler.GetWorkerByID).Methods("GET")
	r.HandleFunc("/workers/{id}", workerHandler.DeleteWorker).Methods("DELETE")

	// === Shifts endpoints ===
	r.HandleFunc("/shifts", shiftHandler.ListShifts).Methods("GET")
	r.HandleFunc("/shifts/available", shiftHandler.ListAvailableShifts).Methods("GET")
	r.HandleFunc("/shifts/assigned", shiftHandler.ListAssignedShifts).Methods("GET")
	r.HandleFunc("/shifts", shiftHandler.CreateShift).Methods("POST")
	r.HandleFunc("/shifts/{id}", shiftHandler.UpdateShift).Methods("PUT")
	r.HandleFunc("/shifts/{id}", shiftHandler.DeleteShift).Methods("DELETE")

	// === Shift Requests endpoints ===
	r.HandleFunc("/shift-requests", shiftRequestHandler.RequestShift).Methods("POST")
	r.HandleFunc("/shift-requests/{id}/approve", shiftRequestHandler.ApproveRequest).Methods("POST")
	r.HandleFunc("/shift-requests/{id}/reject", shiftRequestHandler.RejectRequest).Methods("POST")
	r.HandleFunc("/shift-requests", shiftRequestHandler.GetShiftRequest).Methods("GET")
	r.HandleFunc("/worker/{id}/requests", shiftRequestHandler.ListRequestsByWorker).Methods("GET")

	log.Println("API running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
