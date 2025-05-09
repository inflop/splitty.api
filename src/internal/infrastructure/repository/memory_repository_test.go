package repository_test

import (
	"testing"

	"github.com/inflop/splitty.api/internal/domain/model"
	"github.com/inflop/splitty.api/internal/infrastructure/repository"
)

func TestSaveAndFindByID(t *testing.T) {
	// Utworzenie repozytorium
	repo := repository.NewInMemoryEventRepository()

	// Utworzenie testowego wydarzenia
	event := &model.Event{
		Name: "Test Event",
		Participants: []model.Participant{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
		},
		Expenses: []model.Expense{
			{
				ID:          1,
				Category:    "Food",
				TotalAmount: 100,
				Payments: []model.Payment{
					{ParticipantID: 1, Amount: 100},
				},
				SharedWith: []int{1, 2},
			},
		},
	}

	// Zapisanie wydarzenia
	err := repo.Save(event)
	if err != nil {
		t.Fatalf("Failed to save event: %v", err)
	}

	// Sprawdzenie czy ID zostało przypisane
	if event.ID == 0 {
		t.Fatal("Event ID was not set")
	}

	// Pobranie wydarzenia
	savedEvent, err := repo.FindByID(event.ID)
	if err != nil {
		t.Fatalf("Failed to find event: %v", err)
	}

	// Sprawdzenie danych
	if savedEvent.Name != event.Name {
		t.Errorf("Expected event name %s, got %s", event.Name, savedEvent.Name)
	}

	if len(savedEvent.Participants) != len(event.Participants) {
		t.Errorf("Expected %d participants, got %d", len(event.Participants), len(savedEvent.Participants))
	}

	if len(savedEvent.Expenses) != len(event.Expenses) {
		t.Errorf("Expected %d expenses, got %d", len(event.Expenses), len(savedEvent.Expenses))
	}
}

func TestFindAll(t *testing.T) {
	// Utworzenie repozytorium
	repo := repository.NewInMemoryEventRepository()

	// Dodanie kilku wydarzeń
	for i := 0; i < 3; i++ {
		event := &model.Event{
			Name: "Event " + string(rune('A'+i)),
		}
		err := repo.Save(event)
		if err != nil {
			t.Fatalf("Failed to save event: %v", err)
		}
	}

	// Pobranie wszystkich wydarzeń
	events, err := repo.FindAll()
	if err != nil {
		t.Fatalf("Failed to find all events: %v", err)
	}

	// Sprawdzenie liczby wydarzeń
	if len(events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(events))
	}
}

func TestDelete(t *testing.T) {
	// Utworzenie repozytorium
	repo := repository.NewInMemoryEventRepository()

	// Dodanie wydarzenia
	event := &model.Event{
		Name: "Test Event",
	}
	err := repo.Save(event)
	if err != nil {
		t.Fatalf("Failed to save event: %v", err)
	}

	// Usunięcie wydarzenia
	err = repo.Delete(event.ID)
	if err != nil {
		t.Fatalf("Failed to delete event: %v", err)
	}

	// Próba pobrania usuniętego wydarzenia
	_, err = repo.FindByID(event.ID)
	if err == nil {
		t.Fatal("Expected error when finding deleted event")
	}
}
