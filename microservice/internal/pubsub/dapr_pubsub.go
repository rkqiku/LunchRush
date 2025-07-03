package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/ordis/lunchservice/internal/models"
)

type DaprPubSub struct {
	client      dapr.Client
	pubsubName  string
	topicName   string
}

func NewDaprPubSub(pubsubName, topicName string) *DaprPubSub {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}

	return &DaprPubSub{
		client:     client,
		pubsubName: pubsubName,
		topicName:  topicName,
	}
}

func (p *DaprPubSub) PublishEvent(ctx context.Context, event *models.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := p.client.PublishEvent(ctx, p.pubsubName, p.topicName, data); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (p *DaprPubSub) Close() error {
	if p.client != nil {
		p.client.Close()
	}
	return nil
}

// Event types
const (
	EventTypeSessionCreated     = "session.created"
	EventTypeSessionLocked      = "session.locked"
	EventTypeParticipantJoined  = "participant.joined"
	EventTypeParticipantLeft    = "participant.left"
	EventTypeMealUpdated        = "meal.updated"
	EventTypeRestaurantProposed = "restaurant.proposed"
	EventTypeRestaurantVoted    = "restaurant.voted"
	EventTypeOrderPlacerSet     = "orderplacer.set"
)