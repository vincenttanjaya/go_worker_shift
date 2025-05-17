package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vincenttanjaya/go_worker_shift/internal/controller"
)

type WorkerHandler struct {
	Ctrl *controller.WorkerController
}

func (h *WorkerHandler) SignupWorker(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	worker, err := h.Ctrl.SignupWorker(req.Name)
	if err != nil {
		http.Error(w, "Could not sign up worker", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(worker)
}

func (h *WorkerHandler) ListWorkers(w http.ResponseWriter, r *http.Request) {
	workers, err := h.Ctrl.ListWorkers()
	if err != nil {
		http.Error(w, "Could not fetch workers", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workers)
}

func (h *WorkerHandler) GetWorkerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid worker ID", http.StatusBadRequest)
		return
	}
	worker, err := h.Ctrl.GetWorkerByID(id)
	if err != nil {
		http.Error(w, "Worker not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(worker)
}

func (h *WorkerHandler) DeleteWorker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid worker ID", http.StatusBadRequest)
		return
	}
	err = h.Ctrl.DeleteWorker(id)
	if err != nil {
		http.Error(w, "Could not delete worker", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
