package models

import (
	"time"
)

type LunchSession struct {
	ID           string         `json:"id"`
	Date         time.Time      `json:"date"`
	Status       string         `json:"status"` // "open" | "locked"
	Restaurant   Restaurant     `json:"restaurant"` // Deprecated: kept for backward compatibility
	Restaurants  []Restaurant   `json:"restaurants"` // New: list of proposed restaurants
	Participants []Participant  `json:"participants"`
	OrderPlacer  *string        `json:"orderPlacer,omitempty"`
	LockTime     time.Time      `json:"lockTime"`
	Locked       bool           `json:"locked"`
}

type Restaurant struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Votes      int      `json:"votes"`
	ProposedBy string   `json:"proposedBy"`
	Voters     []string `json:"voters"` // List of usernames who voted
}

type Participant struct {
	UserID       string    `json:"userId"`
	Username     string    `json:"username"`
	MealChoice   string    `json:"mealChoice"`
	Meal         string    `json:"meal"` // Alias for mealChoice
	JoinedAt     time.Time `json:"joinedAt"`
	IsOrderPlacer bool     `json:"isOrderPlacer"`
	LastSeen     time.Time `json:"lastSeen"`
	IsActive     bool      `json:"isActive"`
}

type VoteRequest struct {
	UserID       string `json:"userId"`
	Username     string `json:"username"`
	RestaurantID string `json:"restaurantId"`
}

type JoinRequest struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
}

type MealUpdateRequest struct {
	UserID     string `json:"userId"`
	Username   string `json:"username"`
	MealChoice string `json:"mealChoice,omitempty"`
	Meal       string `json:"meal,omitempty"`
}

type RestaurantProposal struct {
	Name       string `json:"name"`
	ProposedBy string `json:"proposedBy"`
	UserID     string `json:"userId,omitempty"`
}

type Event struct {
	Type      string      `json:"type"`
	SessionID string      `json:"sessionId"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type OrderPlacerRequest struct {
	UserID   string `json:"userId,omitempty"`
	Username string `json:"username"`
}

type HeartbeatRequest struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
}