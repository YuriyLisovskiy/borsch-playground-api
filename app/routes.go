/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) addV1Routes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")

	apiV1.GET("/ping", jsonHandler(pingHandler))
	apiV1.GET("/lang/versions", jsonHandler(a.getLanguageVersionsHandler))

	a.addJobRoutes(apiV1)
}

func (a *Application) getLanguageVersionsHandler(c *gin.Context) (int, interface{}, error) {
	return http.StatusOK, a.settings.BorschVersions, nil
}

func pingHandler(*gin.Context) (int, interface{}, error) {
	return http.StatusOK, gin.H{"message": "pong"}, nil
}
