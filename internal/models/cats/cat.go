package cats

import "errors"

var (
	ErrCatNotFound = errors.New("cat not found")
)

type Cat struct {
	ID                uint   `json:"id"`
	Name              string `json:"name"`
	YearsOfExperience uint8  `json:"years_of_experience"`
	Breed             string `json:"breed"`
	Salary            uint   `json:"salary"`
}
