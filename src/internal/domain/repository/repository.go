package repository

import "github.com/inflop/splitty/internal/domain/model"

// EventRepository definiuje interfejs dla repozytorium wydarzeń
type EventRepository interface {
	Save(event *model.Event) error
	FindByID(id int) (*model.Event, error)
	Delete(id int) error
	FindAll() ([]*model.Event, error)
}
