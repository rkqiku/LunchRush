package models

import (
	"time"
)

// Session represents a daily lunch session
type Session struct {
	ID            string       `json:"id"`             // Unique session ID (e.g., "2025-07-01_team1")
	Date          time.Time    `json:"date"`           // Session date
	Status        string       `json:"status"`         // "open" or "locked"
	Restaurants   []Restaurant `json:"restaurants"`    // List of proposed restaurants
	Orders        []Order      `json:"orders"`         // User orders
	NominatedUser string       `json:"nominated_user"` // User responsible for placing the order
	LockTime      time.Time    `json:"lock_time"`      // Time when session locks
	CreatedAt     time.Time    `json:"created_at"`     // Session creation timestamp
	UpdatedAt     time.Time    `json:"updated_at"`     // Last update timestamp
}

// Restaurant represents a proposed restaurant with votes
type Restaurant struct {
	ID        string     `json:"id"`         // Unique restaurant ID
	Name      string     `json:"name"`       // Restaurant name (e.g., "Pizza Place")
	MenuItems []MenuItem `json:"menu_items"` // Available dishes
	Votes     []Vote     `json:"votes"`      // Anonymous votes or reactions
}

// MenuItem represents a dish on a restaurant’s menu
type MenuItem struct {
	ID          string  `json:"id"`          // Unique item ID
	Name        string  `json:"name"`        // Dish name (e.g., "Margherita Pizza")
	Price       float64 `json:"price"`       // Price in dollars
	Description string  `json:"description"` // Optional description
}

// Vote represents an anonymous vote or reaction for a restaurant
type Vote struct {
	VoteID    string    `json:"vote_id"`    // Unique vote ID
	Type      string    `json:"type"`       // "upvote" or "reaction" (e.g., emoji)
	CreatedAt time.Time `json:"created_at"` // Vote timestamp
}

// Order represents a user’s lunch order
type Order struct {
	UserID       string    `json:"user_id"`       // User ID (from Huly)
	RestaurantID string    `json:"restaurant_id"` // Selected restaurant
	MenuItemID   string    `json:"menu_item_id"`  // Selected dish
	Preferences  string    `json:"preferences"`   // Optional dietary preferences
	CreatedAt    time.Time `json:"created_at"`    // Order timestamp
}
