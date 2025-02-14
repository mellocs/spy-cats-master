package missions

import (
	"errors"
	"spy-cats/internal/models/missions"
	missionsRep "spy-cats/internal/repository/missions"
)

type Missions struct {
	missionsRepository missionsRep.MissionRepository
	targetsRepository  missionsRep.TargetRepository
}

func NewMissions(missionRepo missionsRep.MissionRepository, targetsRepo missionsRep.TargetRepository) Missions {
	return Missions{missionsRepository: missionRepo, targetsRepository: targetsRepo}
}

func (m *Missions) GetAll() ([]missions.Mission, error) {
	allMissions, err := m.missionsRepository.GetAll()
	if err != nil {
		return nil, err
	}

	for i := range allMissions {
		targets, err := m.missionsRepository.GetTargetsByMissionID(allMissions[i].ID)
		if err != nil {
			return nil, err
		}
		allMissions[i].Targets = targets
	}

	return allMissions, nil
}

func (m *Missions) GetById(id int) (missions.Mission, error) {
	mission, err := m.missionsRepository.GetByID(id)
	if err != nil {
		return missions.Mission{}, err
	}

	mission.Targets, err = m.missionsRepository.GetTargetsByMissionID(mission.ID)
	if err != nil {
		return mission, err
	}

	return mission, nil
}

func (m *Missions) Create(mission missions.Mission) error {
	missionId, err := m.missionsRepository.Create(mission)
	if err != nil {
		return err
	}

	if len(mission.Targets) > 0 {
		for _, target := range mission.Targets {
			target.MissionID = uint(missionId)
			err := m.targetsRepository.Create(target)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Missions) Delete(id int) error {
	mission, err := m.missionsRepository.GetByID(id)
	if err != nil {
		return err
	}

	if mission.CatID != 0 {
		return errors.New("mission has an executor")
	}
	return m.missionsRepository.Delete(id)
}

func (m *Missions) AssignCat(id, catId int) error {
	return m.missionsRepository.AssignCat(id, catId)
}

func (m *Missions) CompleteMission(id int) error {
	targets, err := m.missionsRepository.GetTargetsByMissionID(uint(id))
	if err != nil {
		return err
	}

	isAllowedToComplete := true
	for _, target := range targets {
		if target.Completed == false {
			isAllowedToComplete = false
		}
	}
	if !isAllowedToComplete {
		return errors.New("targets isn't completed")
	}

	err = m.missionsRepository.CompleteMission(id)

	return err
}
