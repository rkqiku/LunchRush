package main

import "time"

// Order represents a single meal order by a participant
type Order struct {
	Restaurant string `json:"restaurant"`
	Dish       string `json:"dish"`
	Notes      string `json:"notes,omitempty"`
}

// Participant represents a user in the lunch session
type Participant struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Order *Order `json:"order,omitempty"`
}

// LunchSession represents the lunch session for a given day
type LunchSession struct {
	ID              string         `json:"id"`
	Date            time.Time      `json:"date"`
	Participants    []Participant  `json:"participants"`
	Locked          bool           `json:"locked"`
	LockedBy        string         `json:"lockedBy,omitempty"`
	LockedAt        *time.Time     `json:"lockedAt,omitempty"`
	LockAt          *time.Time     `json:"lockAt,omitempty"`
	NominatedUserID string         `json:"nominatedUserId,omitempty"`
	Votes           map[string]int `json:"votes,omitempty"`
} 