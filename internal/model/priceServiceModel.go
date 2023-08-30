// Package model of our entity and subscription
package model

import (
	"sync"

	"github.com/google/uuid"
)

// Share struct represents a model for shares
type Share struct {
	ShareName  string      `json:"share_name"`
	SharePrice interface{} `json:"price"`
}

// PubSub struct represents a model for subscriptions
type PubSub struct {
	Mu     sync.RWMutex
	Subs   map[uuid.UUID][]chan map[string]float64
	Closed bool
}

// NewPubSub creates a new PubSub instance
func NewPubSub() *PubSub {
	return &PubSub{Subs: make(map[uuid.UUID][]chan map[string]float64)}
}
