package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/ordis/lunchservice/internal/models"
)

type DaprStateStore struct {
	client    dapr.Client
	storeName string
}

func NewDaprStateStore(storeName string) *DaprStateStore {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	
	return &DaprStateStore{
		client:    client,
		storeName: storeName,
	}
}

func (s *DaprStateStore) SaveSession(ctx context.Context, session *models.LunchSession) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Save session by ID
	if err := s.client.SaveState(ctx, s.storeName, session.ID, data, nil); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	// Also save reference for today's session
	today := time.Now().Format("2006-01-02")
	if err := s.client.SaveState(ctx, s.storeName, "session:today:"+today, []byte(session.ID), nil); err != nil {
		return fmt.Errorf("failed to save today's session reference: %w", err)
	}

	return nil
}

func (s *DaprStateStore) GetSession(ctx context.Context, id string) (*models.LunchSession, error) {
	item, err := s.client.GetState(ctx, s.storeName, id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if item.Value == nil {
		return nil, fmt.Errorf("session not found")
	}

	var session models.LunchSession
	if err := json.Unmarshal(item.Value, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

func (s *DaprStateStore) GetTodaySession(ctx context.Context) (*models.LunchSession, error) {
	today := time.Now().Format("2006-01-02")
	
	// Get today's session ID
	item, err := s.client.GetState(ctx, s.storeName, "session:today:"+today, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's session reference: %w", err)
	}

	if item.Value == nil {
		return nil, nil // No session for today
	}

	sessionID := string(item.Value)
	return s.GetSession(ctx, sessionID)
}

func (s *DaprStateStore) UpdateSession(ctx context.Context, session *models.LunchSession) error {
	return s.SaveSession(ctx, session)
}

func (s *DaprStateStore) Close() error {
	if s.client != nil {
		s.client.Close()
	}
	return nil
}