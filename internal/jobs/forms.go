/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package jobs

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const MaxLimit = 100

type GetJobForm struct {
	JobId       string
	RequestHost string
	RequestURI  string
}

func ParseGetJobForm(ctx *gin.Context) (*GetJobForm, error) {
	return &GetJobForm{
		JobId:       ctx.Param("id"),
		RequestHost: ctx.Request.Host,
		RequestURI:  ctx.Request.RequestURI,
	}, nil
}

type CreateJobForm struct {
	LangVersion string `json:"lang_version"`
	SourceCode  string `json:"source_code"`
	RequestHost string `json:"-"`
	RequestURI  string `json:"-"`
}

func ParseCreateJobForm(ctx *gin.Context, availableInterpreters []string, rHost, rURI string) (*CreateJobForm, error) {
	form := &CreateJobForm{
		RequestHost: rHost,
		RequestURI:  rURI,
	}
	err := ctx.ShouldBindJSON(form)
	if err != nil {
		return nil, err
	}

	if !stringArrayHasItem(availableInterpreters, form.LangVersion) {
		return nil, errors.New("language version does not exist")
	}

	if len(form.SourceCode) == 0 {
		return nil, errors.New("source code is not provided")
	}

	return form, nil
}

type GetJobOutputsForm struct {
	JobId  string
	Offset int
	Limit  int
}

func ParseGetJobOutputsForm(ctx *gin.Context) (*GetJobOutputsForm, error) {
	jobId := ctx.Param("id")
	if jobId == "" {
		return nil, errors.New("job id is empty")
	}

	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "-1"))
	if err != nil {
		return nil, errors.New("offset is invalid integer value")
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "-1"))
	if err != nil {
		return nil, errors.New("limit is invalid integer value")
	}

	if limit > MaxLimit {
		return nil, errors.New("limit is too large, 100 is maximum value")
	}

	return &GetJobOutputsForm{
		JobId:  jobId,
		Offset: offset,
		Limit:  limit,
	}, nil
}

func stringArrayHasItem(array []string, item string) bool {
	for _, elem := range array {
		if elem == item {
			return true
		}
	}

	return false
}
