// Package model of our entity
package model

// Share struct represents a model for shares
type Share struct {
	ShareName  string      `json:"share_name"`
	SharePrice interface{} `json:"price"`
}
