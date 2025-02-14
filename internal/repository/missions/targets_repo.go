package missions

import (
	"database/sql"
	"errors"
	"spy-cats/internal/models/missions"
)

type TargetRepository interface {
	Create(target missions.Target) error
	GetByID(id int) (missions.Target, error)
	AddTargetToMission(id, missionId int) error
	CompleteTarget(id int) error
	UpdateNotes(id int, notes string) error
	DeleteMission(id int) error
}

type Targets struct {
	db *sql.DB
}

func NewTargets(db *sql.DB) *Targets {
	return &Targets{db: db}
}

func (r *Targets) Create(target missions.Target) error {
	_, err := r.db.Exec("INSERT INTO targets (mission_id, name, country, notes) VALUES ($1, $2, $3, $4)",
		target.MissionID, target.Name, target.Country, target.Notes)

	return err
}

func (r *Targets) GetByID(id int) (missions.Target, error) {
	var target missions.Target

	err := r.db.QueryRow("SELECT id, mission_id, name, country, notes, completed FROM targets WHERE id = $1", id).
		Scan(&target.ID, &target.MissionID, &target.Name, &target.Country, &target.Notes, &target.Completed)
	if errors.Is(err, sql.ErrNoRows) {
		return missions.Target{}, missions.ErrTargetNotFound
	}

	return target, nil
}

func (r *Targets) AddTargetToMission(id, missionId int) error {
	_, err := r.db.Exec("UPDATE targets SET mission_id=$1 WHERE id=$2", missionId, id)

	return err
}

func (r *Targets) CompleteTarget(id int) error {
	_, err := r.db.Exec("UPDATE targets SET completed=true WHERE id=$1", id)

	return err
}

func (r *Targets) UpdateNotes(id int, notes string) error {
	_, err := r.db.Exec(`UPDATE targets SET notes=$1 WHERE id=$2`, notes, id)

	return err
}

func (r *Targets) DeleteMission(id int) error {
	_, err := r.db.Exec("UPDATE targets SET mission_id=NULL WHERE id=$1", id)

	return err
}
