package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vincenttanjaya/go_worker_shift/internal/controller"
)

type ShiftRequestHandler struct {
	Ctrl *controller.ShiftRequestController
}

func (h *ShiftRequestHandler) RequestShift(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WorkerID int `json:"worker_id"`
		ShiftID  int `json:"shift_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.Ctrl.RequestShift(req.WorkerID, req.ShiftID); err != nil {
		http.Error(w, "Could not request shift", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func (h *ShiftRequestHandler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err := h.Ctrl.ApproveRequest(id); err != nil {
		http.Error(w, "Could not approve request", http.StatusBadRequest)
		fmt.Println("Failed to approve request:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (h *ShiftRequestHandler) RejectRequest(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err := h.Ctrl.RejectRequest(id); err != nil {
		http.Error(w, "Could not reject request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (h *ShiftRequestHandler) ListRequestsByWorker(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid worker ID", http.StatusBadRequest)
		return
	}
	reqs, err := h.Ctrl.ListRequestsByWorker(id)
	if err != nil {
		http.Error(w, "Could not fetch requests", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reqs)
}

func (h *ShiftRequestHandler) GetShiftRequest(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	requests, err := h.Ctrl.ListRequestsByStatus(status)
	if err != nil {
		http.Error(w, "Could not fetch requests", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}
