package model

type CursorToken struct {
	LastKey   string `json:"lastKey"`   // "<created_at RFC3339Nano>|<id>"
	Direction string `json:"direction"` // "forward"
	Limit     int    `json:"limit"`
	Exp       int64  `json:"exp,omitempty"`
}
