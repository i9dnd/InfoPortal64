package storage

import (
	"encoding/json"
	"os"
	"sync"

	"event-website/internal/models"
)

type EventStorage struct {
	events    []models.Event
	eventsMap map[string]models.Event
	mu        sync.RWMutex
}

func NewEventStorage() *EventStorage {
	return &EventStorage{
		events:    []models.Event{},
		eventsMap: make(map[string]models.Event),
	}
}

func (s *EventStorage) Load(filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &s.events)
	if err != nil {
		return err
	}
	for _, event := range s.events {
		s.eventsMap[event.ID] = event
	}
	return nil
}

func (s *EventStorage) Save(filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(s.events)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, os.ModePerm)
}

func (s *EventStorage) Add(event models.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = append(s.events, event)
	s.eventsMap[event.ID] = event
}

func (s *EventStorage) Edit(event models.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, e := range s.events {
		if e.ID == event.ID {
			s.events[i] = event
			break
		}
	}
	s.eventsMap[event.ID] = event
}

func (s *EventStorage) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.eventsMap[id]; !exists {
		return false
	}

	delete(s.eventsMap, id)
	for i, event := range s.events {
		if event.ID == id {
			s.events = append(s.events[:i], s.events[i+1:]...)
			break
		}
	}
	return true
}

func (s *EventStorage) GetAll() []models.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.events
}

func (s *EventStorage) GetEventByID(id string) (models.Event, bool) { //Исправить ошибку 7.01
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, exists := s.eventsMap[id]
	return event, exists
}
