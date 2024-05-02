package main

import (
	"errors"
	"fmt"
	"net/http"
	"github.com/OGElla/Project-API/internal/data"
	"github.com/OGElla/Project-API/internal/validator"
)

func (app *application) createGoalHandler(w http.ResponseWriter, r *http.Request){
	var input struct{
		Walking data.Walking `json:"walking"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil{
		app.badRequestResponse(w, r, err)
		return
	}

	goal := &data.Goals{
		Walking: input.Walking,
		UserId: data.CurrentUserID,
	}

	//VALIDATION (99)
	v := validator.New()
	if data.ValidateGoal(v, goal); !v.Valid(){
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Goals.Insert(goal)
	if err!=nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/goals/view/%d", goal.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"goal": goal}, headers)
	if err != nil{
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showGoalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
   if err != nil{
	   app.notFoundResponse(w, r)
	   return
   }

   goal, err := app.models.Goals.Get(id, data.CurrentUserID)
   if err != nil{
	   switch {
	   case errors.Is(err, data.ErrRecordNotFound):
		   app.notFoundResponse(w, r)
	   default:
		   app.serverErrorResponse(w, r, err)
	   }
	   return
   }

   err = app.writeJSON(w, http.StatusOK, envelope{"goal":goal}, nil)
   if err != nil {
	   app.serverErrorResponse(w, r, err)
   }
}

func (app *application) updateGoalHandler(w http.ResponseWriter, r *http.Request){
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return 
	}

	goal, err := app.models.Goals.Get(id, data.CurrentUserID)
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
		Walking data.Walking `json:"walking"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	goal.Walking = input.Walking

	v := validator.New()
	if data.ValidateGoal(v, goal); !v.Valid(){
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Goals.Update(goal)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"goal": goal}, nil)
	if err != nil{
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteGoalHandler(w http.ResponseWriter, r *http.Request){
	id, err := app.readIDParam(r)
	if err != nil{
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Goals.Delete(id, data.CurrentUserID)
	if err != nil{
		switch{
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "goal successfully deleted"}, nil)
	if err != nil{
		app.serverErrorResponse(w, r, err)
	}
}