package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"envrionment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJson(w, http.StatusOK, env, nil)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
	}

}
