package router

import (
	"github.com/gorilla/mux"
	"github.com/inflop/splitty/internal/infrastructure/api/handler"
)

// SetupRoutes konfiguruje ścieżki API
func SetupRoutes(eventHandler *handler.EventHandler) *mux.Router {
	router := mux.NewRouter()

	// Definiowanie endpointów API
	router.HandleFunc("/api/events", eventHandler.CreateEvent).Methods("POST")
	router.HandleFunc("/api/events", eventHandler.GetAllEvents).Methods("GET")
	router.HandleFunc("/api/events/{id}", eventHandler.GetEvent).Methods("GET")
	router.HandleFunc("/api/events/{id}", eventHandler.UpdateEvent).Methods("PUT")
	router.HandleFunc("/api/events/{id}", eventHandler.DeleteEvent).Methods("DELETE")
	router.HandleFunc("/api/events/{id}/summary", eventHandler.GetEventSummary).Methods("GET")

	return router
}
