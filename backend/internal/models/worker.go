package models

type Worker struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Possibly add Email, etc.
}
