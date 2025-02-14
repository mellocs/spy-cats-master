package cat

import (
	"database/sql"
	"errors"
	"spy-cats/internal/models/cats"
)

type Repository interface {
	GetAll() ([]cats.Cat, error)
	GetByID(id int) (cats.Cat, error)
	Create(cat cats.Cat) error
	Update(id int, salary uint) error
	Delete(id int) error
}

type Cats struct {
	db *sql.DB
}

func NewCats(db *sql.DB) *Cats {
	return &Cats{db}
}

func (c *Cats) GetAll() ([]cats.Cat, error) {
	rows, err := c.db.Query("SELECT id, name, years_of_experience, breed, salary FROM cats")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	catSlice := make([]cats.Cat, 0)
	for rows.Next() {
		var cat cats.Cat
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary); err != nil {
			return nil, err
		}

		catSlice = append(catSlice, cat)
	}

	return catSlice, rows.Err()
}

func (c *Cats) GetByID(id int) (cats.Cat, error) {
	var cat cats.Cat
	err := c.db.QueryRow("SELECT id, name, years_of_experience, breed, salary FROM cats WHERE id=$1", id).
		Scan(&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary)
	if errors.Is(err, sql.ErrNoRows) {
		return cat, cats.ErrCatNotFound
	}

	return cat, err
}

func (c *Cats) Create(cat cats.Cat) error {
	_, err := c.db.Exec(
		"INSERT INTO cats (name, years_of_experience, breed, salary) values ($1, $2, $3, $4)",
		cat.Name, cat.YearsOfExperience, cat.Breed, cat.Salary)

	return err
}

func (c *Cats) Delete(id int) error {
	_, err := c.db.Exec("DELETE FROM cats WHERE id=$1", id)

	return err
}

func (c *Cats) Update(id int, salary uint) error {
	_, err := c.db.Exec("UPDATE cats SET salary=$1 WHERE id=$2", salary, id)

	return err
}
