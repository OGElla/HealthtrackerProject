package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"fmt"
	"github.com/OGElla/Project-API/internal/validator"
)

var MaxWalking int64

type Health struct {
	ID int64 `json:"-"` //unique integer ID
	CreatedAt time.Time `json:"created at"` //timestamp 
	Calories Calories `json:"calories,omitempty"` //calories
	Walking Walking `json:"walking,omitempty"` //steps
	Hydrate Hydrate `json:"hydrate,omitempty"` //water
	Sleep Sleep `json:"sleep,omitempty"`//time
	UserId int64 `json:"-"` // id of user
	Version int32 `json:"version"`//the version number
}
//reusable
func ValidateDaily(v *validator.Validator, health *Health) {
	v.Check(health.Walking > 0, "walking", "must be a positive integer")
}

type HealthModel struct{
	DB *sql.DB
}

func (m HealthModel) Insert(health *Health) error{

	var currentDay string

	now := time.Now()
	day := now.Format("2006-01-02")

	firstQuery := `SELECT walking FROM goals WHERE date(created_at) = $1 and user_id = $2`

	err := m.DB.QueryRow(firstQuery, string(day), CurrentUserID).Scan(
		&MaxWalking,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):	
		default:
			return err
		}
	}	

	tempQuery := `SELECT date(created_at) FROM healthtracker WHERE user_id = $1 ORDER BY created_at DESC
	LIMIT 1`

	err = m.DB.QueryRow(tempQuery, CurrentUserID).Scan(
		&currentDay,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
		default:
			return err
		}
	}

	t, _ := time.Parse(time.RFC3339, currentDay)
	
	dateString := t.Format("2006-01-02")

	if dateString ==  day{
		return ErrTrackerAlreadyCreated
	}

	query := `INSERT INTO healthtracker (calories, walking, hydrate, sleep, user_id) VALUES($1, $2, $3, $4, $5) RETURNING id, created_at, version` 

	args := []interface{}{health.Calories, health.Walking, health.Hydrate, health.Sleep, health.UserId}

	if int64(health.Walking) >= MaxWalking {
		midQuery := `UPDATE goals
		SET achieved = true
		WHERE user_id = $1 and date(created_at) = $2	`

		_, err := m.DB.Exec(midQuery, CurrentUserID, string(day))
		if err != nil{
			return err
		}
		
	}

	return  m.DB.QueryRow(query, args...).Scan(&health.ID, &health.CreatedAt, &health.Version)
}

func (m HealthModel) Get(id int64, userID int64) (*Health, error) {
	if id < 1{
		return nil, ErrRecordNotFound
	}

	var exID int64;

	tempQuery := `
	select id from (select row_number() over(order by user_id, created_at) as row_id, id from healthtracker where user_id = $1) as foo where row_id = $2; 
	`

	err := m.DB.QueryRow(tempQuery, userID, id).Scan(
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

	query := `SELECT id, created_at,calories, walking, hydrate, sleep, version FROM healthtracker WHERE id = $1 and user_id = $2`

	var health Health

	err = m.DB.QueryRow(query, exID, userID).Scan(
		&health.ID,
		&health.CreatedAt,
		&health.Calories,
		&health.Walking,
		&health.Hydrate,
		&health.Sleep,
		&health.Version,
	)

	if err != nil{
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &health, nil
}

func (m HealthModel) Update(health *Health) error {


	query := `UPDATE healthtracker SET calories = $1, walking = $2, hydrate = $3, sleep =$4, version = version + 1 WHERE id = $5 and user_id = $6 RETURNING version`

	args := []interface{}{
		health.Calories,
		health.Walking,
		health.Hydrate, 
		health.Sleep,
		health.ID,
		CurrentUserID,
	}

	qoalQuery := `SELECT user_id FROM goals WHERE $1 > walking AND date(created_at) = $2`

	var id int64

	now := time.Now()
	day := now.Format("2006-01-02")

	err := m.DB.QueryRow(qoalQuery, &health.Walking,day).Scan(
		&id,
	)

	fmt.Println(id)


	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):	
		default:
			return err
		}
	}	

	fmt.Println(id)
	fmt.Println(CurrentUserID)
	if id == CurrentUserID{
		changeQuery := `UPDATE goals SET achieved = true WHERE user_id = $1 AND date(created_at) = $2`
		_, err := m.DB.Exec(changeQuery, CurrentUserID, day)
		if err != nil{
			return err
		}
	}

	if id != CurrentUserID{
		changeQuery := `UPDATE goals SET achieved = false WHERE user_id = $1 AND date(created_at) = $2`
		_, err := m.DB.Exec(changeQuery, CurrentUserID, day)
		if err != nil{
			return err
		}
	}


	return m.DB.QueryRow(query, args...).Scan(&health.Version)
}

func (m HealthModel) Delete(id int64, userID int64) error {
	if id < 1{
		return ErrRecordNotFound
	}

	var exID int64;

	tempQuery := `
	select id from (select row_number() over(order by user_id, created_at) as row_id, id from healthtracker where user_id = $1) as foo where row_id = $2; 
	`

	err := m.DB.QueryRow(tempQuery, userID, id).Scan(
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

	query := `DELETE FROM healthtracker WHERE id = $1 and user_id = $2`
	result, err := m.DB.Exec(query, exID, userID)
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


func (m HealthModel) GetAll(calories int, walking int, hydrate int, sleep int, filters Filters, userID int64) ([]*Health, Metadata, error) {

	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, created_at, calories, walking, hydrate, sleep, version
	FROM healthtracker 
	WHERE user_id = $1
	AND (calories > $2 OR $2 = 0)
	AND (walking > $3 OR $3 = 0)
	AND (hydrate > $4 OR $4 = 0)
	AND (sleep > $5 OR $5 = 0)
	ORDER BY %s %s, id ASC
	LIMIT $6 OFFSET $7
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{userID, calories, walking, hydrate, sleep, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()	

	totalRecords := 0
	healthes := []*Health{}

	for rows.Next() {
		var health Health

		err := rows.Scan(
			&totalRecords,
			&health.ID,
			&health.CreatedAt,
			&health.Calories,
			&health.Walking,
			&health.Hydrate,
			&health.Sleep,
			&health.Version,
		)

		if err != nil{
			return nil, Metadata{}, err
		}

		healthes = append(healthes, &health)
	}

	if err = rows.Err(); err != nil{
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return healthes, metadata, nil
}

