package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler{
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/", app.homeHandler)

	router.HandlerFunc(http.MethodGet, "/health/daily", app.requirePermission("health:read", app.listHealthHandler))
	//TRACKING
	router.HandlerFunc(http.MethodGet, "/health/view/:id", app.requirePermission("health:read", app.showHealthHandler))
	router.HandlerFunc(http.MethodPost, "/health/daily", app.requirePermission("health:read", app.createActivityHandler))
	router.HandlerFunc(http.MethodPut, "/health/view/:id", app.requirePermission("health:read", app.updateHealthHandler))
	router.HandlerFunc(http.MethodDelete, "/health/view/:id", app.requirePermission("health:write", app.deleteHealthHandler))
	//GOAL
	router.HandlerFunc(http.MethodGet, "/goals/view/:id", app.requirePermission("health:read", app.showGoalHandler))
	router.HandlerFunc(http.MethodPost, "/goals/daily", app.requirePermission("health:read", app.createGoalHandler))
	router.HandlerFunc(http.MethodPut, "/goals/view/:id", app.requirePermission("health:read", app.updateGoalHandler))
	router.HandlerFunc(http.MethodDelete, "/goals/view/:id", app.requirePermission("health:write", app.deleteGoalHandler))

	router.HandlerFunc(http.MethodPost, "/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.authenticate(router)
}