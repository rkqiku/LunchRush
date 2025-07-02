package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ordis/lunchservice/internal/handlers"
	"github.com/ordis/lunchservice/internal/pubsub"
	"github.com/ordis/lunchservice/internal/store"
)

const (
	appPort         = "8080"
	daprHTTPPort    = "3500"
	daprGRPCPort    = "50001"
	stateStoreName  = "statestore"
	pubsubName      = "pubsub"
	pubsubTopic     = "lunch-events"
)

func main() {
	// Initialize Dapr client
	stateStore := store.NewDaprStateStore(stateStoreName)
	pubSub := pubsub.NewDaprPubSub(pubsubName, pubsubTopic)

	// Initialize handlers
	h := handlers.NewHandler(stateStore, pubSub)

	// Create Chi router
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/health"))
	
	// Enable CORS for frontend
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Add HTTP routes
	router.Route("/sessions", func(r chi.Router) {
		r.Post("/", h.CreateSession)
		r.Get("/today", h.GetTodaySession)
		r.Get("/{id}", h.GetSession)
		r.Post("/{id}/join", h.JoinSession)
		r.Put("/{id}/meal", h.UpdateMeal)
		r.Post("/{id}/lock", h.LockSession)
		r.Get("/{id}/restaurants", h.GetRestaurants)
		r.Post("/{id}/restaurants", h.ProposeRestaurant)
		r.Delete("/{id}/restaurants/{restaurantId}", h.DeleteRestaurant)
		r.Post("/{id}/restaurants/{restaurantId}/vote", h.VoteRestaurant)
		r.Put("/{id}/order-placer", h.SetOrderPlacer)
		r.Post("/{id}/heartbeat", h.Heartbeat)
		r.Delete("/{id}/participants/{username}", h.RemoveParticipant)
	})

	// Add route for pub/sub events
	router.Post("/events", func(w http.ResponseWriter, r *http.Request) {
		// This will be handled by Dapr SDK
		w.WriteHeader(http.StatusOK)
	})

	// Create Dapr service with Chi router
	s := daprd.NewServiceWithMux(":"+appPort, router)

	// Setup Dapr subscription
	sub := &common.Subscription{
		PubsubName: pubsubName,
		Topic:      pubsubTopic,
		Route:      "/events",
	}

	// Add topic subscription handler
	if err := s.AddTopicEventHandler(sub, eventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	// Start server
	log.Printf("Starting LunchRush service on port %s", appPort)
	
	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	// Start the server
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error starting server: %v", err)
	}
}

// eventHandler handles incoming pub/sub events
func eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	log.Printf("Received event - PubsubName: %s, Topic: %s, ID: %s", 
		e.PubsubName, e.Topic, e.ID)
	
	// Log the event data
	fmt.Printf("Event data: %s\n", e.Data)
	
	// Successfully processed, no retry needed
	return false, nil
}