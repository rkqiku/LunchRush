package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ordis/lunchservice/internal/models"
	"github.com/ordis/lunchservice/internal/pubsub"
	"github.com/ordis/lunchservice/internal/store"
)

type Handler struct {
	store  *store.DaprStateStore
	pubsub *pubsub.DaprPubSub
}

func NewHandler(store *store.DaprStateStore, pubsub *pubsub.DaprPubSub) *Handler {
	return &Handler{
		store:  store,
		pubsub: pubsub,
	}
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check if session already exists for today
	existingSession, err := h.store.GetTodaySession(ctx)
	if err == nil && existingSession != nil {
		h.respondJSON(w, http.StatusConflict, map[string]string{"error": "Session already exists for today"})
		return
	}

	// Create new session
	session := &models.LunchSession{
		ID:           uuid.New().String(),
		Date:         time.Now(),
		Status:       "open",
		Participants: []models.Participant{},
		Restaurants:  []models.Restaurant{},
		LockTime:     time.Now().Add(3 * time.Hour), // Default lock time at 3 hours from now
		Locked:       false,
	}

	if err := h.store.SaveSession(ctx, session); err != nil {
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Publish event
	event := &models.Event{
		Type:      pubsub.EventTypeSessionCreated,
		SessionID: session.ID,
		Data:      session,
		Timestamp: time.Now(),
	}
	h.pubsub.PublishEvent(ctx, event)

	h.respondJSON(w, http.StatusCreated, session)
}

func (h *Handler) GetTodaySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, err := h.store.GetTodaySession(ctx)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if session == nil {
		h.respondJSON(w, http.StatusNotFound, map[string]string{"error": "No session found for today"})
		return
	}

	// Ensure arrays are initialized
	if session.Restaurants == nil {
		session.Restaurants = []models.Restaurant{}
	}
	if session.Participants == nil {
		session.Participants = []models.Participant{}
	}

	h.respondJSON(w, http.StatusOK, session)
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := chi.URLParam(r, "id")

	session, err := h.store.GetSession(ctx, sessionID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err)
		return
	}

	// Ensure arrays are initialized
	if session.Restaurants == nil {
		session.Restaurants = []models.Restaurant{}
	}
	if session.Participants == nil {
		session.Participants = []models.Participant{}
	}

	h.respondJSON(w, http.StatusOK, session)
}

func (h *Handler) JoinSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := chi.URLParam(r, "id")

	var req models.JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	session, err := h.store.GetSession(ctx, sessionID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err)
		return
	}

	if session.Status == "locked" {
		h.respondJSON(w, http.StatusForbidden, map[string]string{"error": "Session is locked"})
		return
	}

	// Check if user already joined by username
	for _, p := range session.Participants {
		if p.Username == req.Username {
			h.respondJSON(w, http.StatusConflict, map[string]string{"error": "User already joined"})
			return
		}
	}

	// Add participant
	participant := models.Participant{
		UserID:   req.UserID,
		Username: req.Username,
		JoinedAt: time.Now(),
		LastSeen: time.Now(),
		IsActive: true,
	}
	session.Participants = append(session.Participants, participant)

	if err := h.store.UpdateSession(ctx, session); err != nil {
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Publish event
	event := &models.Event{
		Type:      pubsub.EventTypeParticipantJoined,
		SessionID: session.ID,
		Data:      participant,
		Timestamp: time.Now(),
	}
	h.pubsub.PublishEvent(ctx, event)

	h.respondJSON(w, http.StatusOK, session)
}

func (h *Handler) UpdateMeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := chi.URLParam(r, "id")

	var req models.MealUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	session, err := h.store.GetSession(ctx, sessionID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err)
		return
	}

	if session.Status == "locked" {
		h.respondJSON(w, http.StatusForbidden, map[string]string{"error": "Session is locked"})
		return
	}

	// Update participant's meal choice
	meal := req.Meal
	if meal == "" {
		meal = req.MealChoice
	}
	
	found := false
	for i, p := range session.Participants {
		if p.Username == req.Username {
			session.Participants[i].MealChoice = meal
			session.Participants[i].Meal = meal
			found = true
			break
		}
	}

	if !found {
		h.respondJSON(w, http.StatusNotFound, map[string]string{"error": "User not found in session"})
		return
	}

	if err := h.store.UpdateSession(ctx, session); err != nil {
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Publish event
	event := &models.Event{
		Type:      pubsub.EventTypeMealUpdated,
		SessionID: session.ID,
		Data:      req,
		Timestamp: time.Now(),
	}
	h.pubsub.PublishEvent(ctx, event)

	h.respondJSON(w, http.StatusOK, session)
}

func (h *Handler) LockSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := chi.URLParam(r, "id")

	session, err := h.store.GetSession(ctx, sessionID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err)
		return
	}

	if session.Status == "locked" {
		h.respondJSON(w, http.StatusConflict, map[string]string{"error": "Session already locked"})
		return
	}

	session.Status = "locked"
	session.Locked = true
	
	if err := h.store.UpdateSession(ctx, session); err != nil {
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Publish event
	event := &models.Event{
		Type:      pubsub.EventTypeSessionLocked,
		SessionID: session.ID,
		Data:      session,
		Timestamp: time.Now(),
	}
	h.pubsub.PublishEvent(ctx, event)

	h.respondJSON(w, http.StatusOK, session)
}

func (h *Handler) ProposeRestaurant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := chi.URLParam(r, "id")

	var req models.RestaurantProposal
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	session, err := h.store.GetSession(ctx, sessionID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err)
		return
	}

	if session.Status == "locked" {
		h.respondJSON(w, http.StatusForbidden, map[string]string{"error": "Session is locked"})
		return
	}

	restaurant := models.Restaurant{
		ID:         uuid.New().String(),
		Name:       req.Name,
		Votes:      0,
		ProposedBy: req.ProposedBy,
		Voters:     []string{},
	}

	// Check if restaurant already exists
	for _, r := range session.Restaurants {
		if r.Name == req.Name {
			h.respondJSON(w, http.StatusConflict, map[string]string{"error": "Restaurant already proposed"})
			return
		}
	}

	// Add to restaurants list
	session.Restaurants = append(session.Restaurants, restaurant)
	
	// Keep backward compatibility
	if len(session.Restaurants) == 1 {
		session.Restaurant = restaurant
	}

	if err := h.store.UpdateSession(ctx, session); err != nil {
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Publish event
	event := &models.Event{
		Type:      pubsub.EventTypeRestaurantProposed,
		SessionID: session.ID,
		Data:      restaurant,
		Timestamp: time.Now(),
	}
	h.pubsub.PublishEvent(ctx, event)

	h.respondJSON(w, http.StatusOK, restaurant)
}

func (h *Handler) VoteRestaurant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := chi.URLParam(r, "id")
	restaurantID := chi.URLParam(r, "restaurantId")

	var req models.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	session, err := h.store.GetSession(ctx, sessionID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err)
		return
	}

	if session.Status == "locked" {
		h.respondJSON(w, http.StatusForbidden, map[string]string{"error": "Session is locked"})
		return
	}

	// Find restaurant and update votes
	var foundRestaurant *models.Restaurant
	for i := range session.Restaurants {
		if session.Restaurants[i].ID == restaurantID {
			// Check if user already voted
			for _, voter := range session.Restaurants[i].Voters {
				if voter == req.Username {
					// Remove vote
					newVoters := []string{}
					for _, v := range session.Restaurants[i].Voters {
						if v != req.Username {
							newVoters = append(newVoters, v)
						}
					}
					session.Restaurants[i].Voters = newVoters
					session.Restaurants[i].Votes = len(newVoters)
					foundRestaurant = &session.Restaurants[i]
					break
				}
			}
			
			if foundRestaurant == nil {
				// Add vote
				session.Restaurants[i].Voters = append(session.Restaurants[i].Voters, req.Username)
				session.Restaurants[i].Votes = len(session.Restaurants[i].Voters)
				foundRestaurant = &session.Restaurants[i]
			}
			break
		}
	}
	
	if foundRestaurant == nil {
		h.respondJSON(w, http.StatusNotFound, map[string]string{"error": "Restaurant not found"})
		return
	}
	
	// Update backward compatibility field
	if session.Restaurant.ID == restaurantID {
		session.Restaurant = *foundRestaurant
	}
	
	if err := h.store.UpdateSession(ctx, session); err != nil {
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Publish event
	event := &models.Event{
		Type:      pubsub.EventTypeRestaurantVoted,
		SessionID: session.ID,
		Data:      req,
		Timestamp: time.Now(),
	}
	h.pubsub.PublishEvent(ctx, event)
	
	h.respondJSON(w, http.StatusOK, foundRestaurant)
}

func (h *Handler) SetOrderPlacer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := chi.URLParam(r, "id")

	var req models.OrderPlacerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	session, err := h.store.GetSession(ctx, sessionID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err)
		return
	}

	if session.Status == "locked" {
		h.respondJSON(w, http.StatusForbidden, map[string]string{"error": "Session is locked"})
		return
	}

	// Find participant and set as order placer
	for i := range session.Participants {
		if session.Participants[i].Username == req.Username {
			session.Participants[i].IsOrderPlacer = true
			session.OrderPlacer = &req.Username
		} else {
			session.Participants[i].IsOrderPlacer = false
		}
	}

	if err := h.store.UpdateSession(ctx, session); err != nil {
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Publish event
	event := &models.Event{
		Type:      pubsub.EventTypeOrderPlacerSet,
		SessionID: session.ID,
		Data:      req,
		Timestamp: time.Now(),
	}
	h.pubsub.PublishEvent(ctx, event)

	h.respondJSON(w, http.StatusOK, session)
}

func (h *Handler) GetRestaurants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := chi.URLParam(r, "id")

	session, err := h.store.GetSession(ctx, sessionID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err)
		return
	}

	if session.Restaurants == nil {
		session.Restaurants = []models.Restaurant{}
	}

	h.respondJSON(w, http.StatusOK, session.Restaurants)
}

func (h *Handler) DeleteRestaurant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := chi.URLParam(r, "id")
	restaurantID := chi.URLParam(r, "restaurantId")

	session, err := h.store.GetSession(ctx, sessionID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err)
		return
	}

	if session.Status == "locked" {
		h.respondJSON(w, http.StatusForbidden, map[string]string{"error": "Session is locked"})
		return
	}

	// Remove restaurant from list
	newRestaurants := []models.Restaurant{}
	deleted := false
	for _, r := range session.Restaurants {
		if r.ID != restaurantID {
			newRestaurants = append(newRestaurants, r)
		} else {
			deleted = true
		}
	}

	if !deleted {
		h.respondJSON(w, http.StatusNotFound, map[string]string{"error": "Restaurant not found"})
		return
	}

	session.Restaurants = newRestaurants

	if err := h.store.UpdateSession(ctx, session); err != nil {
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Restaurant deleted"})
}

func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := chi.URLParam(r, "id")

	var req models.HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	session, err := h.store.GetSession(ctx, sessionID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err)
		return
	}

	// Update participant's last seen timestamp
	found := false
	for i, p := range session.Participants {
		if p.Username == req.Username {
			session.Participants[i].LastSeen = time.Now()
			session.Participants[i].IsActive = true
			found = true
			break
		}
	}

	if !found {
		h.respondJSON(w, http.StatusNotFound, map[string]string{"error": "User not found in session"})
		return
	}

	// Clean up inactive users (not seen in last 5 minutes)
	activeParticipants := []models.Participant{}
	inactiveThreshold := time.Now().Add(-5 * time.Minute)
	
	for _, p := range session.Participants {
		if p.LastSeen.After(inactiveThreshold) {
			activeParticipants = append(activeParticipants, p)
		} else {
			// Mark as inactive but keep in list
			p.IsActive = false
			activeParticipants = append(activeParticipants, p)
		}
	}
	
	session.Participants = activeParticipants

	if err := h.store.UpdateSession(ctx, session); err != nil {
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) RemoveParticipant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := chi.URLParam(r, "id")
	username := chi.URLParam(r, "username")

	session, err := h.store.GetSession(ctx, sessionID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err)
		return
	}

	if session.Status == "locked" {
		h.respondJSON(w, http.StatusForbidden, map[string]string{"error": "Session is locked"})
		return
	}

	// Remove participant
	newParticipants := []models.Participant{}
	removed := false
	for _, p := range session.Participants {
		if p.Username != username {
			newParticipants = append(newParticipants, p)
		} else {
			removed = true
		}
	}

	if !removed {
		h.respondJSON(w, http.StatusNotFound, map[string]string{"error": "Participant not found"})
		return
	}

	session.Participants = newParticipants

	// Also remove their votes from restaurants
	for i := range session.Restaurants {
		newVoters := []string{}
		for _, voter := range session.Restaurants[i].Voters {
			if voter != username {
				newVoters = append(newVoters, voter)
			}
		}
		session.Restaurants[i].Voters = newVoters
		session.Restaurants[i].Votes = len(newVoters)
	}

	if err := h.store.UpdateSession(ctx, session); err != nil {
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Publish event
	event := &models.Event{
		Type:      pubsub.EventTypeParticipantLeft,
		SessionID: session.ID,
		Data:      map[string]string{"username": username},
		Timestamp: time.Now(),
	}
	h.pubsub.PublishEvent(ctx, event)

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Participant removed"})
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, err error) {
	if err != nil {
		h.respondJSON(w, status, map[string]string{"error": err.Error()})
	} else {
		h.respondJSON(w, status, map[string]string{"error": "Unknown error"})
	}
}