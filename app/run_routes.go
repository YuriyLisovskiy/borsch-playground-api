package app

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/YuriyLisovskiy/borsch-playground-service/core"
	"github.com/YuriyLisovskiy/borsch-playground-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (a *Application) addRunRoutes(rg *gin.RouterGroup) {
	runRouter := rg.Group("/run")

	runRouter.POST("/", jsonHandler(a.runHandler))
}

func (a *Application) runHandler(c *gin.Context) (int, interface{}, error) {
	var form CreateJobForm
	err := c.ShouldBindJSON(&form)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if !stringArrayContains(a.settings.BorschVersions, form.LangV) {
		return http.StatusBadRequest, nil, errors.New("language version does not exist")
	}

	if len(form.B64Source) == 0 {
		return http.StatusBadRequest, nil, errors.New("code is not provided or empty")
	}

	code, err := base64.StdEncoding.DecodeString(form.B64Source)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	job := &models.JobDbModel{
		Model: models.Model{
			ID: uuid.New().String(),
		},
		Code:     string(code),
		Outputs:  []models.JobOutputRowDbModel{},
		ExitCode: nil,
	}

	err = a.db.Create(job).Error
	if err != nil {
		return -1, nil, err
	}

	err = a.enqueueJob(job, form.LangV)
	if err != nil {
		switch err {
		case core.ErrQueueIsFull:
			return http.StatusServiceUnavailable, nil, errors.New("Service temporary is unavailable, try again later")
		case core.ErrQueueIsUnavailable:
			fallthrough
		default:
			return -1, nil, err
		}
	}

	return http.StatusCreated, gin.H{"job_id": job.ID}, nil
}
