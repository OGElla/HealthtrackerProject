package data

import(
	"database/sql"
	"time"
	"errors"
	"github.com/OGElla/Project-API/internal/validator"
)

type Goals struct{
	ID int64 `json:"-"` //unique integer ID
	CreatedAt time.Time `json:"created at"` //timestamp 
	Walking Walking `json:"walking,omitempty"` //steps
	UserId int64 `json:"-"` // id of user
	Achieved bool `json:"achieved"` // achieve of goal
	Version int32 `json:"version"`//the version number
}

type GoalModel struct{
	DB *sql.DB
}

func ValidateGoal(v *validator.Validator, goal *Goals) {
	v.Check(goal.Walking > 0, "walking", "must be a positive integer")
}

func (g GoalModel) Insert(goal *Goals) error{

	var id int64

	now := time.Now()
	day := now.Format("2006-01-02")
	tempQuery := `SELECT user_id FROM goals WHERE date(created_at) = $1`

	err := g.DB.QueryRow(tempQuery, string(day)).Scan(
		&id,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			if goal.UserId == id {
				return ErrGoalAlreadyCreated
			}
		default:
			return err
		}
	}

	if goal.UserId == id {
		return ErrGoalAlreadyCreated
	}
	
	query := `INSERT INTO goals (walking, user_id) VALUES($1, $2) RETURNING id, created_at, version, achieved` 

	args := []interface{}{goal.Walking, goal.UserId}

	return g.DB.QueryRow(query, args...).Scan(&goal.ID, &goal.CreatedAt, &goal.Version, &goal.Achieved)
}

func (g GoalModel) Get(id int64, userID int64) (*Goals, error) {
	if id < 1{
		return nil, ErrRecordNotFound
	}

	var exID int64;

	tempQuery := `
	select id from (select row_number() over(order by user_id, created_at) as row_id, id from goals where user_id = $1) as foo where row_id = $2; 
	`

	err := g.DB.QueryRow(tempQuery, userID, id).Scan(
		&exID,
	)

	if err != nil{
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	query := `SELECT id, created_at, walking, version, achieved FROM goals WHERE id = $1 and user_id = $2`

	var goal Goals

	err = g.DB.QueryRow(query, exID, userID).Scan(
		&goal.ID,
		&goal.CreatedAt,
		&goal.Walking,
		&goal.Version,
		&goal.Achieved,
	)

	if err != nil{
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &goal, nil
}

func (g GoalModel) Update(goal *Goals) error {


	query := `UPDATE goals SET walking = $1, version = version + 1 WHERE id = $2 and user_id = $3 RETURNING version`

	args := []interface{}{
		&goal.Walking,
		&goal.ID,
		CurrentUserID,
	}

	return g.DB.QueryRow(query, args...).Scan(&goal.Version)
}

func (g GoalModel) Delete(id int64, userID int64) error {
	if id < 1{
		return ErrRecordNotFound
	}

	var exID int64;

	tempQuery := `
	select id from (select row_number() over(order by user_id, created_at) as row_id, id from goals where user_id = $1) as foo where row_id = $2; 
	`

	err := g.DB.QueryRow(tempQuery, userID, id).Scan(
		&exID,
	)

	if err != nil{
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	query := `DELETE FROM goals WHERE id = $1 and user_id = $2`
	result, err := g.DB.Exec(query, exID, userID)
	if err != nil{
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
