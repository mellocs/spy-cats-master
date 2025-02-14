package missions

import "errors"

var (
	ErrMissionNotFound = errors.New("mission not found")
)

type Mission struct {
	ID        uint     `json:"id"`
	CatID     uint     `json:"cat_id,omitempty"`
	Completed bool     `json:"completed"`
	Targets   []Target `json:"targets,omitempty"`
}
