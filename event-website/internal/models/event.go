package models

import "github.com/google/uuid"

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"` //Сделать кастомную категорию 9.01
}

func NewEvent(title, description, category string) Event {
	return Event{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Category:    category,
	}
}
