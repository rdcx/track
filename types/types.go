package types

import "time"

type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type HitMap map[string]map[string]time.Time
