// Package model of our entity
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

type PubSub struct {
	Mu   sync.RWMutex
	Subs map[uuid.UUID][]chan *Share
	//Shares map[uuid.UUID][]string   // shares for concrete client(ID)
	Closed bool
}

func NewPubSub() *PubSub {
	return &PubSub{Subs: make(map[uuid.UUID][]chan *Share)}
}
