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
	apiV1.GET("/lang/versions", a.getLanguageVersionsHandler)

	jobsRouter := apiV1.Group("/jobs")
	jobsRouter.GET("/:id", a.getJobHandler)
	jobsRouter.GET("/:id/output", a.getJobOutputHandler)
	jobsRouter.POST("/", a.createJobHandler)
}

func (a *Application) getLanguageVersionsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, a.settings.BorschVersions)
}
