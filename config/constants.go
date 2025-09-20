package config

import (
	"errors"

	"gorm.io/gorm"
)

type ServiceCode int

const (
	CodeOK                  ServiceCode = 0
	CodeBadRequest          ServiceCode = 1
	CodeUnprocessableEntity ServiceCode = 2
	CodeDatabaseError       ServiceCode = 3
	CodeNotFound            ServiceCode = 4
	CodeUnauthorized        ServiceCode = 5
	CodeForbidden           ServiceCode = 6
	CodeConflict            ServiceCode = 7
	CodeExternalRequestFail ServiceCode = 8
)

var ( // Errors
	ErrRecordNotFound = gorm.ErrRecordNotFound

	ErrCatNotFound = errors.New("cat not found")

	ErrMissionNotFound             = errors.New("mission not found")
	ErrMissionAlreadyAssigned      = errors.New("mission has cat already assigned")
	ErrMissionAlreadyComplete      = errors.New("mission already complete")
	ErrMissionHasMaxTargets        = errors.New("mission already has max number of targets")
	ErrMissionHasInvalidTargetsLen = errors.New("mission has invalid number of targets")

	ErrTargetNotFound        = errors.New("target not found")
	ErrTargetAlreadyComplete = errors.New("target already complete")
)

const (
	MinMissionTargets = 1
	MaxMissionTargets = 3
)
