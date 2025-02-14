package missions

import (
	"errors"
	"spy-cats/internal/models/missions"
	missionsRep "spy-cats/internal/repository/missions"
)

type Targets struct {
	missionsRepository missionsRep.MissionRepository
	targetsRepository  missionsRep.TargetRepository
}

func NewTargets(missionRepo missionsRep.MissionRepository, targetsRepo missionsRep.TargetRepository) Targets {
	return Targets{missionsRepository: missionRepo, targetsRepository: targetsRepo}
}

func (t *Targets) Create(target missions.Target) error {
	return t.targetsRepository.Create(target)
}

func (t *Targets) AddTargetToMission(targetId, missionId int) error {
	mission, err := t.missionsRepository.GetByID(missionId)
	if err != nil {
		return err
	}
	if !mission.Completed {
		err := t.targetsRepository.AddTargetToMission(targetId, missionId)
		if err != nil {
			return err
		}
	}

	return err
}

func (t *Targets) CompleteTarget(targetId int) error {
	return t.targetsRepository.CompleteTarget(targetId)
}

func (t *Targets) UpdateNotes(targetId int, notes string) error {
	target, err := t.targetsRepository.GetByID(targetId)
	if err != nil {
		return err
	}
	if target.Completed {
		return errors.New("target already completed")
	}

	return t.targetsRepository.UpdateNotes(targetId, notes)
}

func (t *Targets) DeleteTargetFromMission(targetId int) error {
	target, err := t.targetsRepository.GetByID(targetId)
	if err != nil {
		return err
	}
	if target.Completed {
		return errors.New("target already completed")
	}

	return t.targetsRepository.DeleteMission(targetId)
}
