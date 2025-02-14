package missions

import "errors"

var (
	ErrTargetNotFound = errors.New("target not found")
)

type Target struct {
	ID        uint   `json:"id"`
	MissionID uint   `json:"mission_id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Notes     string `json:"notes"`
	Completed bool   `json:"completed"`
}
