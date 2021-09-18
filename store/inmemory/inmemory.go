package inmemory

import (
	"context"
	"sync"

	"github.com/eugeneradionov/froxy/models"
	"github.com/google/uuid"
)

type Store struct {
	mu    *sync.RWMutex
	store map[uuid.UUID]*models.Stream
}

func NewStore() *Store {
	return &Store{
		mu:    &sync.RWMutex{},
		store: make(map[uuid.UUID]*models.Stream),
	}
}

func (s *Store) CreatStream(ctx context.Context) (*models.Stream, error) {
	stream := &models.Stream{ID: uuid.New(), Chunks: make(map[uint]*models.Chunk)}
	s.mu.Lock()
	s.store[stream.ID] = stream
	s.mu.Unlock()

	return stream, nil
}

func (s *Store) CreateChunk(ctx context.Context, streamID uuid.UUID, chunk *models.Chunk, position uint) error {
	s.mu.RLock()
	stream, ok := s.store[streamID]
	s.mu.RUnlock()

	if !ok {
		return models.ErrRecordNotFound
	}

	s.mu.Lock()
	stream.Chunks[position] = chunk
	s.mu.Unlock()

	return nil
}

func (s *Store) GetStream(ctx context.Context, streamID uuid.UUID) (*models.Stream, error) {
	s.mu.RLock()
	stream, ok := s.store[streamID]
	s.mu.RUnlock()

	if !ok {
		return nil, models.ErrRecordNotFound
	}

	return stream, nil
}
