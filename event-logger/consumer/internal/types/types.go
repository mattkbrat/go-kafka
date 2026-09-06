package types

import "time"

type (
	Event struct {
		EventID   string    `json:"event_id" bigquery:"event_id"`
		UserID    string    `json:"user_id" bigquery:"user_id"`
		EventType string    `json:"event_type" bigquery:"event_type"`
		Message   string    `json:"message" bigquery:"message"`
		TS        time.Time `json:"ts" bigquery:"ts"`
	}
)
