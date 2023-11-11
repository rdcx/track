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

type DomainMap map[Domain]KeyMap
type KeyMap map[Key]UrlMap
type UrlMap map[Url]HitSlice
type HitSlice []Hit

type Hit struct {
	Loc  string    `json:"loc"`
	Time time.Time `json:"time"`
}
