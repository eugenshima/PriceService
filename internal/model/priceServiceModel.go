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
	Mu         sync.RWMutex
	Subs       map[uuid.UUID][]string
	SubsShares map[uuid.UUID]chan []*Share
	Closed     map[uuid.UUID]bool
}

// NewPubSub creates a new PubSub instance
func NewPubSub() *PubSub {
	return &PubSub{
		Subs:       make(map[uuid.UUID][]string),
		SubsShares: make(map[uuid.UUID]chan []*Share),
		Closed:     make(map[uuid.UUID]bool),
	}
}
