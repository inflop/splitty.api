package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/inflop/splitty.api/internal/domain/model"
	"github.com/inflop/splitty.api/internal/domain/repository"
	"github.com/inflop/splitty.api/internal/domain/service"
)

// EventHandler obsługuje zapytania HTTP związane z wydarzeniami
type EventHandler struct {
	eventRepository repository.EventRepository
	expenseService  *service.ExpenseService
}

// NewEventHandler tworzy nowy handler wydarzeń
func NewEventHandler(
	eventRepository repository.EventRepository,
	expenseService *service.ExpenseService,
) *EventHandler {
	return &EventHandler{
		eventRepository: eventRepository,
		expenseService:  expenseService,
	}
}

// CreateEvent tworzy nowe wydarzenie
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event model.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Walidacja danych
	if event.Name == "" {
		http.Error(w, "Event name is required", http.StatusBadRequest)
		return
	}

	if err := h.eventRepository.Save(&event); err != nil {
		http.Error(w, "Failed to save event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// GetEvent pobiera wydarzenie po ID
func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.eventRepository.FindByID(id)
	if err != nil {
		http.Error(w, "Event not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// GetEventSummary oblicza i zwraca podsumowanie wydarzenia
func (h *EventHandler) GetEventSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.eventRepository.FindByID(id)
	if err != nil {
		http.Error(w, "Event not found: "+err.Error(), http.StatusNotFound)
		return
	}

	summary := h.expenseService.CalculateSummary(event)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// UpdateEvent aktualizuje wydarzenie
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Sprawdzamy czy wydarzenie istnieje
	_, err = h.eventRepository.FindByID(id)
	if err != nil {
		http.Error(w, "Event not found: "+err.Error(), http.StatusNotFound)
		return
	}

	var event model.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Walidacja danych
	if event.Name == "" {
		http.Error(w, "Event name is required", http.StatusBadRequest)
		return
	}

	// Ustawiamy ID z URL
	event.ID = id

	if err := h.eventRepository.Save(&event); err != nil {
		http.Error(w, "Failed to update event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// DeleteEvent usuwa wydarzenie
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.eventRepository.Delete(id); err != nil {
		http.Error(w, "Failed to delete event: "+err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllEvents pobiera wszystkie wydarzenia
func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.eventRepository.FindAll()
	if err != nil {
		http.Error(w, "Failed to get events: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
