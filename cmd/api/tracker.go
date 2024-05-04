package main

import (
	"errors"
	"fmt"
	"net/http"
	"github.com/OGElla/Project-API/internal/data"
	"github.com/OGElla/Project-API/internal/validator"
)

func (app *application) createActivityHandler(w http.ResponseWriter, r *http.Request){
	var input struct{
		Calories data.Calories `json:"calories"`
		Walking data.Walking `json:"walking"`
		Hydrate data.Hydrate `json:"hydrate"`
		Sleep data.Sleep `json:"sleep"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil{
		app.badRequestResponse(w, r, err)
		return
	}

	health := &data.Health{
		Calories: input.Calories,
		Walking: input.Walking,
		Hydrate: input.Hydrate,
		Sleep: input.Sleep,
		UserId: data.CurrentUserID,
	}

	//VALIDATION (99)
	v := validator.New()
	if data.ValidateDaily(v, health); !v.Valid(){
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Health.Insert(health)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return 
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/health/view/%d", health.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"health": health}, headers)
	if err != nil{
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showHealthHandler(w http.ResponseWriter, r *http.Request) {
 	id, err := app.readIDParam(r)
	if err != nil{
		app.notFoundResponse(w, r)
		return
	}

	health, err := app.models.Health.Get(id, data.CurrentUserID)
	if err != nil{
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"health":health}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateHealthHandler(w http.ResponseWriter, r *http.Request){
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return 
	}

	health, err := app.models.Health.Get(id, data.CurrentUserID)
	if err != nil{
		switch{
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default: 
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Calories data.Calories `json:"calories"`
		Walking data.Walking `json:"walking"`
		Hydrate data.Hydrate `json:"hydrate"`
		Sleep data.Sleep `json:"sleep"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	health.Calories = input.Calories
	health.Walking = input.Walking
	health.Hydrate = input.Hydrate
	health.Sleep = input.Sleep

	v := validator.New()
	if data.ValidateDaily(v, health); !v.Valid(){
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Health.Update(health)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"health": health}, nil)
	if err != nil{
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteHealthHandler(w http.ResponseWriter, r *http.Request){
	id, err := app.readIDParam(r)
	if err != nil{
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Health.Delete(id, data.CurrentUserID)
	if err != nil{
		switch{
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "healthtracker successfully deleted"}, nil)
	if err != nil{
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listHealthHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Calories int
		Walking int
		Hydrate int
		Sleep int
		data.Filters
		UserID int
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "calories", "hydrate", "walking", "sleep", "-id", "-calories", "-hydrate", "-walking", "-sleep"}//??


	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Walking = app.readInt(qs, "walking", 0, v)
	input.Hydrate = app.readInt(qs, "hydrate", 0, v)
	input.Sleep = app.readInt(qs, "sleep", 0, v)
	input.Calories = app.readInt(qs, "calories", 1, v)


	if data.ValidateFilters(v, input.Filters); !v.Valid(){
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	healthes, metadata, err := app.models.Health.GetAll(input.Calories, input.Walking, input.Hydrate, input.Sleep, input.Filters, data.CurrentUserID)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"healthes":healthes, "metadata":metadata}, nil)
	if err != nil{
		app.serverErrorResponse(w, r, err)
	}
}

