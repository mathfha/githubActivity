package github

import "time"

type Event struct {
	Type      string    `json:"type"`
	Repo      Repo      `json:"repo"`
	Payload   Payload   `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type Repo struct {
	Name string `json:"name"`
}

type Payload struct {
	Action  string   `json:"action,omitempty"`
	Size    int      `json:"size,omitempty"`
	Commits []Commit `json:"commits,omitempty"`
}

type Commit struct {
	Message string `json:"message"`
}
