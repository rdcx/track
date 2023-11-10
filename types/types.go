package types

import "time"

type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type TrackResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Key     Key    `json:"key"`
}

type Domain string
type Key string
type Url string

type HitMap map[Domain]map[Key]map[Url][]time.Time
