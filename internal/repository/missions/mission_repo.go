package missions

import (
	"database/sql"
	"errors"
	"spy-cats/internal/models/missions"
)

type MissionRepository interface {
	GetAll() ([]missions.Mission, error)
	GetByID(id int) (missions.Mission, error)
	Create(mission missions.Mission) (int, error)
	Delete(id int) error
	AssignCat(id, catId int) error
	CompleteMission(id int) error
	GetTargetsByMissionID(id uint) ([]missions.Target, error)
}

type Missions struct {
	db *sql.DB
}

func NewMissions(db *sql.DB) *Missions {
	return &Missions{db: db}
}

func (m *Missions) GetAll() ([]missions.Mission, error) {
	rows, err := m.db.Query("SELECT id, cat_id, completed FROM missions")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var allMissions []missions.Mission
	for rows.Next() {
		var mission missions.Mission
		if err := rows.Scan(&mission.ID, &mission.CatID, &mission.Completed); err != nil {
			return nil, err
		}

		allMissions = append(allMissions, mission)
	}

	return allMissions, nil
}

func (m *Missions) GetByID(id int) (missions.Mission, error) {
	var mission missions.Mission

	err := m.db.QueryRow("SELECT id, cat_id, completed FROM missions WHERE id = $1", id).
		Scan(&mission.ID, &mission.CatID, &mission.Completed)
	if errors.Is(err, sql.ErrNoRows) {
		return missions.Mission{}, missions.ErrMissionNotFound
	}

	return mission, nil
}

func (m *Missions) Create(mission missions.Mission) (int, error) {
	var id int
	err := m.db.QueryRow("INSERT INTO missions (cat_id, completed) VALUES ($1, $2) RETURNING id",
		mission.CatID, mission.Completed).Scan(&id)

	return id, err
}

func (m *Missions) Update(id int, mission missions.Mission) error {
	_, err := m.db.Exec("UPDATE missions SET cat_id = $1, completed = $2 WHERE id = $3",
		mission.CatID, mission.Completed, id)

	return err
}

func (m *Missions) AssignCat(id, catId int) error {
	_, err := m.db.Exec("UPDATE missions SET cat_id = $1 WHERE id = $2", catId, id)

	return err
}

func (m *Missions) CompleteMission(id int) error {
	_, err := m.db.Exec("UPDATE missions SET completed = true WHERE id = $1", id)

	return err
}

func (m *Missions) Delete(id int) error {
	_, err := m.db.Exec("DELETE FROM missions WHERE id = $1", id)

	return err
}

func (m *Missions) GetTargetsByMissionID(id uint) ([]missions.Target, error) {
	var targets []missions.Target

	rows, err := m.db.Query("SELECT id, mission_id, name, country, notes, completed FROM targets WHERE mission_id = $1", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var target missions.Target

		if err := rows.
			Scan(&target.ID, &target.MissionID, &target.Name, &target.Country, &target.Notes, &target.Completed); err != nil {
			return nil, err
		}

		targets = append(targets, target)
	}

	return targets, nil
}
