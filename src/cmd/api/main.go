package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/inflop/splitty/internal/domain/service"
	"github.com/inflop/splitty/internal/infrastructure/api/handler"
	"github.com/inflop/splitty/internal/infrastructure/api/router"
	repo "github.com/inflop/splitty/internal/infrastructure/repository"
	"github.com/rs/cors"
)

func main() {
	// Konfiguracja loggera
	logger := log.New(os.Stdout, "[SPLITTY] ", log.LstdFlags)
	logger.Println("Starting Splitty API...")

	// Inicjalizacja repozytoriów
	eventRepository := repo.NewInMemoryEventRepository()

	// Inicjalizacja usług
	expenseService := service.NewExpenseService()

	// Inicjalizacja handlerów
	eventHandler := handler.NewEventHandler(eventRepository, expenseService)

	// Konfiguracja routera
	r := router.SetupRoutes(eventHandler)

	// Konfiguracja CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // W produkcji należy ograniczyć do konkretnych domen
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400, // 24h w sekundach
	})

	// Konfiguracja serwera
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Domyślny port
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      c.Handler(r),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Uruchomienie serwera
	logger.Printf("Server listening on port %s\n", port)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
