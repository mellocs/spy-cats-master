package cats

import (
	"errors"
	"spy-cats/internal/models/cats"
	"spy-cats/internal/repository/cat"
	"spy-cats/internal/validation"
)

type Cats struct {
	repository cat.Repository
}

func NewCats(repo cat.Repository) Cats {
	return Cats{repository: repo}
}

func (c *Cats) GetAll() ([]cats.Cat, error) {
	return c.repository.GetAll()
}

func (c *Cats) GetByID(id int) (cats.Cat, error) {
	return c.repository.GetByID(id)
}

func (c *Cats) Delete(id int) error {
	return c.repository.Delete(id)
}

func (c *Cats) Create(cat cats.Cat) error {
	isValidBreed, err := validation.ValidateBreed(cat.Breed)
	if err != nil {
		return err
	}

	if !isValidBreed {
		return errors.New("invalid breed")
	}
	return c.repository.Create(cat)
}

func (c *Cats) Update(id int, salary uint) error {
	return c.repository.Update(id, salary)
}
