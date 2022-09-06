package app

import (
	"errors"
	"net/http"

	"github.com/YuriyLisovskiy/borsch-playground-service/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (a *Application) addJobRoutes(rg *gin.RouterGroup) {
	jobRouter := rg.Group("/job")

	jobRouter.GET("/:id", jsonHandler(a.jobHandler))
	jobRouter.GET("/:id/outputs", jsonHandler(a.jobOutputsHandler))
}

func (a *Application) jobHandler(c *gin.Context) (int, interface{}, error) {
	var job models.JobDbModel
	err := a.db.First(&job, "ID = ?", c.Param("id")).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, nil, errors.New("job not found")
		}

		return -1, nil, err
	}

	return http.StatusOK, job, nil
}

func (a *Application) jobOutputsHandler(c *gin.Context) (int, interface{}, error) {
	jobId := c.Param("id")
	var job models.JobDbModel
	err := a.db.Model(&job).
		Preload("Outputs").
		First(&job, "ID = ?", jobId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, nil, errors.New("job not found")
		}

		return -1, nil, err
	}

	return http.StatusOK, job.Outputs, nil
}
