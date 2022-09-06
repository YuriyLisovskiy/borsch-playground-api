package app

import (
	"errors"
	"net/http"
	"strconv"

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
	qOffset := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(qOffset)
	if err != nil {
		jObj := map[string]interface{}{
			"message": "Validation Failed",
			"error": map[string]string{
				"resource": "Job Outputs",
				"value":    qOffset,
			},
			"documentation_url": "TODO:",
		}
		return http.StatusUnprocessableEntity, jObj, nil
	}

	err = a.db.Model(&models.JobDbModel{}).First(nil, "ID = ?", jobId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, nil, errors.New("job not found")
		}

		return -1, nil, err
	}

	var jobOutputs []models.JobOutputRowDbModel
	err = a.db.Offset(offset).Find(&jobOutputs, "job_id = ?", jobId).Error
	if err != nil {
		return -1, nil, err
	}

	return http.StatusOK, jobOutputs, nil
}
