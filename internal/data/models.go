package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")
	ErrTrackerAlreadyCreated = errors.New("today's tracker has already been created")
	ErrGoalAlreadyCreated = errors.New("today's goal has already been created")
)

type Models struct {
	Health HealthModel
	Goals GoalModel
	Permissions PermissionModel
	Tokens TokenModel
	Users UserModel
}

func NewModels(db *sql.DB) Models{
	return Models{
		Health: HealthModel{DB:db},
		Goals: GoalModel{DB:db},
		Permissions: PermissionModel{DB:db},
		Tokens: TokenModel{DB:db},
		Users: UserModel{DB:db},
	}
}

