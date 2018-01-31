package main

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Message struct {
	PersistentLink uuid.UUID `json:"id,omitempty"`
	Message        string    `json:"message,omitempty"`
	CipherMessage  string    `json:"cipher-message,omitempty"`
	OneTimeLink    uuid.UUID `json:"link,omitempty"`
	SeenDate       time.Time `json:"seen-date,omitempty"`
}
