/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package server

import (
	"net/http"

	"github.com/YuriyLisovskiy/borsch-playground-api/internal/jobs"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (a *Application) getJobHandler(ctx *gin.Context) {
	form, err := jobs.ParseGetJobForm(ctx)
	if err != nil {
		a.sendJsonError(ctx, http.StatusBadRequest, err)
	}

	result, err := a.jobService.GetJobResult(form)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, result)
	case gorm.ErrRecordNotFound:
		a.sendJsonError(ctx, http.StatusNotFound, err)
	default:
		a.sendJsonError(ctx, http.StatusInternalServerError, err)
	}
}

func (a *Application) getJobOutputHandler(ctx *gin.Context) {
	form, err := jobs.ParseGetJobOutputsForm(ctx)
	if err != nil {
		a.sendJsonError(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := a.jobService.GetJobOutputsAsJsonResult(form)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, result)
	case gorm.ErrRecordNotFound:
		a.sendJsonError(ctx, http.StatusNotFound, err)
	default:
		a.sendJsonError(ctx, http.StatusInternalServerError, err)
	}
}

func (a *Application) getJobOutputAsTxtHandler(ctx *gin.Context) {
	form, err := jobs.ParseGetJobOutputsForm(ctx)
	if err != nil {
		a.sendJsonError(ctx, http.StatusBadRequest, err)
		return
	}

	err = a.jobService.GetJobOutputsAsTextResult(form, ctx.Writer)
	switch err {
	case nil:
		// skip
	case gorm.ErrRecordNotFound:
		a.sendJsonError(ctx, http.StatusNotFound, err)
	default:
		a.sendJsonError(ctx, http.StatusInternalServerError, err)
	}
}

func (a *Application) createJobHandler(ctx *gin.Context) {
	form, err := jobs.ParseCreateJobForm(
		ctx,
		a.settings.AvailableInterpreters,
		ctx.Request.Host,
		ctx.Request.RequestURI,
	)
	if err != nil {
		a.sendJsonError(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := a.jobService.CreateJobResult(form)
	if err != nil {
		a.sendJsonError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, result)
}
