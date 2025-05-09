package repository

import (
	"errors"
	"sync"

	"github.com/inflop/splitty.api/internal/domain/model"
	"github.com/inflop/splitty.api/internal/domain/repository"
)

// Sprawdzenie czy implementacja spełnia interfejs
var _ repository.EventRepository = (*InMemoryEventRepository)(nil)

// InMemoryEventRepository implementacja repozytorium w pamięci
type InMemoryEventRepository struct {
	events map[int]*model.Event
	nextID int
	mutex  sync.RWMutex
}

// NewInMemoryEventRepository tworzy nowe repozytorium w pamięci
func NewInMemoryEventRepository() *InMemoryEventRepository {
	return &InMemoryEventRepository{
		events: make(map[int]*model.Event),
		nextID: 1,
	}
}

// Save zapisuje wydarzenie
func (r *InMemoryEventRepository) Save(event *model.Event) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if event.ID == 0 {
		event.ID = r.nextID
		r.nextID++
	}

	// Głębokie kopiowanie obiektu aby uniknąć problemów z współdzieleniem referencji
	eventCopy := copyEvent(event)
	r.events[event.ID] = eventCopy

	// Aktualizujemy oryginał, aby otrzymał ID jeśli było 0
	event.ID = eventCopy.ID

	return nil
}

// FindByID znajduje wydarzenie po ID
func (r *InMemoryEventRepository) FindByID(id int) (*model.Event, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	event, exists := r.events[id]
	if !exists {
		return nil, errors.New("event not found")
	}

	// Zwracamy kopię aby uniknąć problemów z współdzieleniem referencji
	return copyEvent(event), nil
}

// Delete usuwa wydarzenie
func (r *InMemoryEventRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.events[id]; !exists {
		return errors.New("event not found")
	}

	delete(r.events, id)
	return nil
}

// FindAll zwraca wszystkie wydarzenia
func (r *InMemoryEventRepository) FindAll() ([]*model.Event, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	events := make([]*model.Event, 0, len(r.events))
	for _, event := range r.events {
		events = append(events, copyEvent(event))
	}

	return events, nil
}

// Funkcja pomocnicza do głębokiego kopiowania obiektów Event
func copyEvent(event *model.Event) *model.Event {
	if event == nil {
		return nil
	}

	newEvent := &model.Event{
		ID:   event.ID,
		Name: event.Name,
	}

	// Kopiowanie uczestników
	if len(event.Participants) > 0 {
		newEvent.Participants = make([]model.Participant, len(event.Participants))
		for i, p := range event.Participants {
			newEvent.Participants[i] = model.Participant{
				ID:    p.ID,
				Name:  p.Name,
				Email: p.Email,
			}
		}
	}

	// Kopiowanie wydatków
	if len(event.Expenses) > 0 {
		newEvent.Expenses = make([]model.Expense, len(event.Expenses))
		for i, e := range event.Expenses {
			expense := model.Expense{
				ID:          e.ID,
				Category:    e.Category,
				TotalAmount: e.TotalAmount,
			}

			// Kopiowanie płatności
			if len(e.Payments) > 0 {
				expense.Payments = make([]model.Payment, len(e.Payments))
				for j, p := range e.Payments {
					expense.Payments[j] = model.Payment{
						ParticipantID: p.ParticipantID,
						Amount:        p.Amount,
					}
				}
			}

			// Kopiowanie sharedWith
			if len(e.SharedWith) > 0 {
				expense.SharedWith = make([]int, len(e.SharedWith))
				copy(expense.SharedWith, e.SharedWith)
			}

			newEvent.Expenses[i] = expense
		}
	}

	return newEvent
}
