package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const sessionListKey = "session-ids"

var (
	daprClient client.Client
	stateStoreName = "statestore"
	pubsubName     = "pubsub"
	eventTopic     = "LunchRushEventTopic"
	orderBinding   = "orderbinding"

	autoLockCheckInterval = 10 * time.Second
	autoLockStop = make(chan struct{})
	autoLockWg sync.WaitGroup
)

// --- WebSocket support ---
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsHub struct {
	clients map[*websocket.Conn]bool
	lock    sync.Mutex
}

var hub = wsHub{
	clients: make(map[*websocket.Conn]bool),
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	hub.lock.Lock()
	hub.clients[conn] = true
	hub.lock.Unlock()
	defer func() {
		hub.lock.Lock()
		delete(hub.clients, conn)
		hub.lock.Unlock()
		conn.Close()
	}()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func broadcastWS(message []byte) {
	hub.lock.Lock()
	defer hub.lock.Unlock()
	for conn := range hub.clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("WebSocket write error:", err)
			conn.Close()
			delete(hub.clients, conn)
		}
	}
}
// --- End WebSocket support ---

type Event struct {
	Type      string      `json:"type"`
	SessionID string      `json:"sessionId"`
	UserID    string      `json:"userId,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

func initDaprClient() {
	var err error
	daprClient, err = client.NewClient()
	if err != nil {
		log.Fatalf("failed to create Dapr client: %v", err)
	}
}

func publishEvent(ctx context.Context, evt Event) {
	b, err := json.Marshal(evt)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)
		return
	}
	err = daprClient.PublishEvent(ctx, pubsubName, eventTopic, b)
	if err != nil {
		log.Printf("failed to publish event: %v", err)
	}
	broadcastWS(b)
}

func sendOrderToBinding(ctx context.Context, session *LunchSession) {
	b, err := json.Marshal(session)
	if err != nil {
		log.Printf("failed to marshal session for binding: %v", err)
		return
	}
	_, err = daprClient.InvokeBinding(ctx, &client.InvokeBindingRequest{
		Name:      orderBinding,
		Operation: "create",
		Data:      b,
	})
	if err != nil {
		log.Printf("failed to invoke order binding: %v", err)
	}
}

func addSessionID(ctx context.Context, id string) {
	item, _ := daprClient.GetState(ctx, stateStoreName, sessionListKey, nil)
	var ids []string
	if item.Value != nil {
		_ = json.Unmarshal(item.Value, &ids)
	}
	for _, sid := range ids {
		if sid == id {
			return // already present
		}
	}
	ids = append(ids, id)
	b, _ := json.Marshal(ids)
	_ = daprClient.SaveState(ctx, stateStoreName, sessionListKey, b, nil)
}

func removeSessionID(ctx context.Context, id string) {
	item, _ := daprClient.GetState(ctx, stateStoreName, sessionListKey, nil)
	var ids []string
	if item.Value != nil {
		_ = json.Unmarshal(item.Value, &ids)
	}
	newIDs := make([]string, 0, len(ids))
	for _, sid := range ids {
		if sid != id {
			newIDs = append(newIDs, sid)
		}
	}
	b, _ := json.Marshal(newIDs)
	_ = daprClient.SaveState(ctx, stateStoreName, sessionListKey, b, nil)
}

func createSessionHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Date   string  `json:"date"`
		LockAt *string `json:"lockAt,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}
	var lockAt *time.Time
	if req.LockAt != nil {
		t, err := time.Parse(time.RFC3339, *req.LockAt)
		if err == nil {
			lockAt = &t
		}
	}
	session := &LunchSession{
		ID:           uuid.NewString(),
		Date:         parsedDate,
		Participants: []Participant{},
		Locked:       false,
		LockAt:       lockAt,
		Votes:        make(map[string]int),
	}
	ctx := context.Background()
	b, err := json.Marshal(session)
	if err != nil {
		http.Error(w, "failed to marshal session", http.StatusInternalServerError)
		return
	}
	err = daprClient.SaveState(ctx, stateStoreName, "session:"+session.ID, b, nil)
	if err != nil {
		http.Error(w, "failed to save session", http.StatusInternalServerError)
		return
	}
	addSessionID(ctx, session.ID)
	publishEvent(ctx, Event{
		Type:      "SessionCreated",
		SessionID: session.ID,
		Payload:   session,
		Timestamp: time.Now(),
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func getSessionHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := context.Background()
	item, err := daprClient.GetState(ctx, stateStoreName, "session:"+id, nil)
	if err != nil || item.Value == nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}
	var session LunchSession
	if err := json.Unmarshal(item.Value, &session); err != nil {
		http.Error(w, "failed to unmarshal session", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func joinSessionHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := context.Background()
	item, err := daprClient.GetState(ctx, stateStoreName, "session:"+id, nil)
	if err != nil || item.Value == nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}
	var session LunchSession
	if err := json.Unmarshal(item.Value, &session); err != nil {
		http.Error(w, "failed to unmarshal session", http.StatusInternalServerError)
		return
	}
	var req struct {
		UserID string `json:"userId"`
		Name   string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	for _, p := range session.Participants {
		if p.ID == req.UserID {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(session)
			return
		}
	}
	participant := Participant{ID: req.UserID, Name: req.Name}
	session.Participants = append(session.Participants, participant)
	b, err := json.Marshal(session)
	if err != nil {
		http.Error(w, "failed to marshal session", http.StatusInternalServerError)
		return
	}
	err = daprClient.SaveState(ctx, stateStoreName, "session:"+id, b, nil)
	if err != nil {
		http.Error(w, "failed to save session", http.StatusInternalServerError)
		return
	}
	publishEvent(ctx, Event{
		Type:      "UserJoined",
		SessionID: id,
		UserID:    req.UserID,
		Payload:   participant,
		Timestamp: time.Now(),
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := context.Background()
	item, err := daprClient.GetState(ctx, stateStoreName, "session:"+id, nil)
	if err != nil || item.Value == nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}
	var session LunchSession
	if err := json.Unmarshal(item.Value, &session); err != nil {
		http.Error(w, "failed to unmarshal session", http.StatusInternalServerError)
		return
	}
	var req struct {
		UserID     string `json:"userId"`
		Restaurant string `json:"restaurant"`
		Dish       string `json:"dish"`
		Notes      string `json:"notes,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	found := false
	for i, p := range session.Participants {
		if p.ID == req.UserID {
			order := &Order{Restaurant: req.Restaurant, Dish: req.Dish, Notes: req.Notes}
			session.Participants[i].Order = order
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "participant not found", http.StatusNotFound)
		return
	}
	b, err := json.Marshal(session)
	if err != nil {
		http.Error(w, "failed to marshal session", http.StatusInternalServerError)
		return
	}
	err = daprClient.SaveState(ctx, stateStoreName, "session:"+id, b, nil)
	if err != nil {
		http.Error(w, "failed to save session", http.StatusInternalServerError)
		return
	}
	publishEvent(ctx, Event{
		Type:      "OrderPlacedOrUpdated",
		SessionID: id,
		UserID:    req.UserID,
		Payload:   req,
		Timestamp: time.Now(),
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func lockSessionHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := context.Background()
	item, err := daprClient.GetState(ctx, stateStoreName, "session:"+id, nil)
	if err != nil || item.Value == nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}
	var session LunchSession
	if err := json.Unmarshal(item.Value, &session); err != nil {
		http.Error(w, "failed to unmarshal session", http.StatusInternalServerError)
		return
	}
	var req struct {
		UserID string `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if session.Locked {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(session)
		return
	}
	now := time.Now()
	session.Locked = true
	session.LockedBy = req.UserID
	session.LockedAt = &now
	b, err := json.Marshal(session)
	if err != nil {
		http.Error(w, "failed to marshal session", http.StatusInternalServerError)
		return
	}
	err = daprClient.SaveState(ctx, stateStoreName, "session:"+id, b, nil)
	if err != nil {
		http.Error(w, "failed to save session", http.StatusInternalServerError)
		return
	}
	removeSessionID(ctx, id)
	publishEvent(ctx, Event{
		Type:      "SessionLocked",
		SessionID: id,
		UserID:    req.UserID,
		Timestamp: now,
	})
	sendOrderToBinding(ctx, &session)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func nominateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := context.Background()
	item, err := daprClient.GetState(ctx, stateStoreName, "session:"+id, nil)
	if err != nil || item.Value == nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}
	var session LunchSession
	if err := json.Unmarshal(item.Value, &session); err != nil {
		http.Error(w, "failed to unmarshal session", http.StatusInternalServerError)
		return
	}
	var req struct {
		UserID string `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	session.NominatedUserID = req.UserID
	b, err := json.Marshal(session)
	if err != nil {
		http.Error(w, "failed to marshal session", http.StatusInternalServerError)
		return
	}
	err = daprClient.SaveState(ctx, stateStoreName, "session:"+id, b, nil)
	if err != nil {
		http.Error(w, "failed to save session", http.StatusInternalServerError)
		return
	}
	publishEvent(ctx, Event{
		Type:      "UserNominated",
		SessionID: id,
		UserID:    req.UserID,
		Timestamp: time.Now(),
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := context.Background()
	item, err := daprClient.GetState(ctx, stateStoreName, "session:"+id, nil)
	if err != nil || item.Value == nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}
	var session LunchSession
	if err := json.Unmarshal(item.Value, &session); err != nil {
		http.Error(w, "failed to unmarshal session", http.StatusInternalServerError)
		return
	}
	var req struct {
		Restaurant string `json:"restaurant"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if session.Votes == nil {
		session.Votes = make(map[string]int)
	}
	session.Votes[req.Restaurant]++
	b, err := json.Marshal(session)
	if err != nil {
		http.Error(w, "failed to marshal session", http.StatusInternalServerError)
		return
	}
	err = daprClient.SaveState(ctx, stateStoreName, "session:"+id, b, nil)
	if err != nil {
		http.Error(w, "failed to save session", http.StatusInternalServerError)
		return
	}
	publishEvent(ctx, Event{
		Type:      "VoteCast",
		SessionID: id,
		Payload:   req,
		Timestamp: time.Now(),
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func reorderLastWeekHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := context.Background() // removed unused ctx
	// For demo: get all sessions, find the one from 7 days ago
	// (In production, use a better query/indexing strategy)
	// This is a stub for demonstration
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Reorder last week not implemented in demo"}`))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func startAutoLockTimer() {
	autoLockWg.Add(1)
	go func() {
		defer autoLockWg.Done()
		ticker := time.NewTicker(autoLockCheckInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				checkAndAutoLockSessions()
			case <-autoLockStop:
				return
			}
		}
	}()
}

func checkAndAutoLockSessions() {
	ctx := context.Background()
	item, err := daprClient.GetState(ctx, stateStoreName, sessionListKey, nil)
	if err != nil || item.Value == nil {
		return
	}
	var ids []string
	if err := json.Unmarshal(item.Value, &ids); err != nil {
		return
	}
	for _, id := range ids {
		item, err := daprClient.GetState(ctx, stateStoreName, "session:"+id, nil)
		if err != nil || item.Value == nil {
			continue
		}
		var session LunchSession
		if err := json.Unmarshal(item.Value, &session); err != nil {
			continue
		}
		if session.Locked || session.LockAt == nil {
			continue
		}
		if time.Now().After(*session.LockAt) {
			// Auto-lock
			session.Locked = true
			session.LockedBy = "auto-lock"
			now := time.Now()
			session.LockedAt = &now
			b, _ := json.Marshal(session)
			_ = daprClient.SaveState(ctx, stateStoreName, "session:"+id, b, nil)
			removeSessionID(ctx, id)
			publishEvent(ctx, Event{
				Type:      "SessionLocked",
				SessionID: id,
				UserID:    "auto-lock",
				Timestamp: now,
			})
			sendOrderToBinding(ctx, &session)
		}
	}
}

// CORS middleware
func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	initDaprClient()
	startAutoLockTimer()
	r := chi.NewRouter()

	r.Get("/healthz", healthzHandler)
	r.Post("/session", createSessionHandler)
	r.Get("/session/{id}", getSessionHandler)
	r.Post("/session/{id}/join", joinSessionHandler)
	r.Post("/session/{id}/order", orderHandler)
	r.Post("/session/{id}/lock", lockSessionHandler)
	r.Post("/session/{id}/nominate", nominateHandler)
	r.Post("/session/{id}/vote", voteHandler)
	r.Post("/reorder-last-week", reorderLastWeekHandler)
	r.HandleFunc("/ws", wsHandler)

	log.Println("LunchRush microservice running on :8080")
	http.ListenAndServe(":8080", withCORS(r))
}
