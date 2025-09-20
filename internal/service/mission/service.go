package cat

import (
	"backend/config"
	request "backend/internal/controller/http/request/cat"
	response "backend/internal/controller/http/response/cat"
	"context"
	"fmt"
)

func (s service) GetMissions(ctx context.Context) ([]response.Mission, config.ServiceCode, error) {
	missions, err := s.repo.GetMissions(ctx)
	return response.MissionsToResponse(missions), config.DBErrToServiceCode(err), err
}

func (s service) GetMission(ctx context.Context, missionID uint) (response.Mission, config.ServiceCode, error) {
	mission, err := s.repo.GetMission(ctx, missionID)
	if mission.ID == 0 {
		return response.Mission{}, config.CodeNotFound, config.ErrMissionNotFound
	}
	return response.MissionToResponse(mission), config.DBErrToServiceCode(err), err
}

func (s service) CreateMission(ctx context.Context, mission request.Mission) (response.Mission, config.ServiceCode, error) {
	tx := s.repo.NewTransaction(ctx)
	defer tx.Rollback()

	missionEntity := mission.ToEntity()
	createdMission, err := s.repo.CreateMission(ctx, tx, missionEntity)
	if err != nil {
		return response.Mission{}, config.DBErrToServiceCode(err), fmt.Errorf("create mission err: %v", err)
	}

	if mission.CatID != nil {
		cat, err := s.catRepo.GetCatByID(ctx, *mission.CatID)
		if err != nil {
			return response.Mission{}, config.DBErrToServiceCode(err), fmt.Errorf("get cat err: %v", err)
		}
		if cat.ID == 0 {
			return response.Mission{}, config.CodeNotFound, config.ErrCatNotFound
		}
	}

	targets := mission.Targets.ToEntity(createdMission.ID)
	createdTargets, err := s.targetRepo.CreateTargets(ctx, tx, targets)
	if err != nil {
		return response.Mission{}, config.DBErrToServiceCode(err), fmt.Errorf("create mission targets err: %v", err)
	}

	createdMission.Targets = createdTargets

	err = tx.Commit().Error
	return response.MissionToResponse(createdMission), config.DBErrToServiceCode(err), err
}

func (s service) AssignCat(ctx context.Context, missionID, catID uint) (config.ServiceCode, error) {
	mission, err := s.repo.GetMission(ctx, missionID)
	if err != nil {
		return config.DBErrToServiceCode(err), fmt.Errorf("get mission err: %v", err)
	}
	if mission.ID == 0 {
		return config.CodeNotFound, config.ErrMissionNotFound
	}
	if mission.CatID != nil {
		return config.CodeForbidden, config.ErrMissionAlreadyAssigned
	}

	cat, err := s.catRepo.GetCatByID(ctx, catID)
	if err != nil {
		return config.DBErrToServiceCode(err), fmt.Errorf("get cat err: %v", err)
	}
	if cat.ID == 0 {
		return config.CodeNotFound, config.ErrCatNotFound
	}

	err = s.repo.AssignCat(ctx, missionID, catID)
	return config.DBErrToServiceCode(err), err
}

func (s service) CompleteMission(ctx context.Context, missionID uint) (config.ServiceCode, error) {
	err := s.repo.CompleteMission(ctx, missionID)
	return config.DBErrToServiceCode(err), err
}

func (s service) DeleteMission(ctx context.Context, missionID uint) (config.ServiceCode, error) {
	mission, err := s.repo.GetMission(ctx, missionID)
	if err != nil {
		return config.DBErrToServiceCode(err), fmt.Errorf("get mission err: %v", err)
	}
	if mission.ID == 0 {
		return config.CodeNotFound, config.ErrMissionNotFound
	}
	if mission.CatID != nil {
		return config.CodeForbidden, config.ErrMissionAlreadyAssigned
	}

	err = s.repo.DeleteMission(ctx, missionID)
	return config.DBErrToServiceCode(err), err
}

func (s service) CreateTarget(ctx context.Context, body request.Target, missionID uint) (config.ServiceCode, error) {
	mission, err := s.repo.GetMission(ctx, missionID)
	if err != nil {
		return config.DBErrToServiceCode(err), fmt.Errorf("get mission err: %v", err)
	}
	if mission.ID == 0 {
		return config.CodeNotFound, config.ErrMissionNotFound
	}
	if mission.IsCompleted {
		return config.CodeForbidden, config.ErrMissionAlreadyAssigned
	}
	if len(mission.Targets) == config.MaxMissionTargets {
		return config.CodeForbidden, config.ErrMissionHasMaxTargets
	}

	target := body.ToEntity(missionID)
	err = s.targetRepo.CreateTarget(ctx, target)
	return config.DBErrToServiceCode(err), err
}

func (s service) UpdateTarget(ctx context.Context, body request.UpdateTarget, targetID, missionID uint) (config.ServiceCode, error) {
	mission, err := s.repo.GetMission(ctx, missionID)
	if err != nil {
		return config.DBErrToServiceCode(err), fmt.Errorf("get mission err: %v", err)
	}
	if mission.ID == 0 {
		return config.CodeNotFound, config.ErrMissionNotFound
	}
	if mission.IsCompleted {
		return config.CodeForbidden, config.ErrMissionAlreadyComplete
	}

	target, err := s.targetRepo.GetTargetByID(ctx, targetID, missionID)
	if err != nil {
		return config.DBErrToServiceCode(err), fmt.Errorf("get target err: %v", err)
	}
	if target.ID == 0 {
		return config.CodeNotFound, config.ErrTargetNotFound
	}
	if target.IsCompleted {
		return config.CodeForbidden, config.ErrTargetAlreadyComplete
	}

	err = s.targetRepo.UpdateTarget(ctx, body.ToEntity(targetID, missionID))
	return config.DBErrToServiceCode(err), err
}

func (s service) CompleteTarget(ctx context.Context, targetID, missionID uint) (config.ServiceCode, error) {
	err := s.targetRepo.CompleteTarget(ctx, targetID, missionID)
	return config.DBErrToServiceCode(err), err
}

func (s service) DeleteTarget(ctx context.Context, targetID, missionID uint) (config.ServiceCode, error) {
	target, err := s.targetRepo.GetTargetByID(ctx, targetID, missionID)
	if err != nil {
		return config.DBErrToServiceCode(err), fmt.Errorf("get target err: %v", err)
	}
	if target.ID == 0 {
		return config.CodeNotFound, config.ErrTargetNotFound
	}
	if target.IsCompleted {
		return config.CodeForbidden, config.ErrTargetAlreadyComplete
	}

	err = s.targetRepo.DeleteTarget(ctx, targetID, missionID)
	return config.DBErrToServiceCode(err), err
}
