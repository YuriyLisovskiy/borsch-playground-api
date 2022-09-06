package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/YuriyLisovskiy/borsch-playground-service/core"
	"github.com/YuriyLisovskiy/borsch-playground-service/models"
	"github.com/YuriyLisovskiy/borsch-playground-service/settings"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Application struct {
	settings *settings.Settings
	queue    *core.Queue
	db       *gorm.DB
}

func NewApp(s *settings.Settings) (*Application, error) {
	mode := gin.ReleaseMode
	if s.Debug {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)
	db, err := s.Database.Create()
	if err != nil {
		return nil, err
	}

	jobQueue, err := s.Queue.Create()
	if err != nil {
		return nil, err
	}

	app := &Application{
		settings: s,
		queue:    jobQueue,
		db:       db,
	}
	return app, nil
}

func (a *Application) buildRouter() *gin.Engine {
	router := gin.Default()
	a.addV1Routes(router)
	return router
}

func (a *Application) Execute(addr string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := a.buildRouter()
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), a.settings.ShutdownTimeoutSec*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return errors.New(fmt.Sprintf("Server forced to shut down: %v", err))
	}

	log.Println("Server exiting")
	return nil
}

func (a *Application) enqueueJob(job *models.JobDbModel, interpreterVersion string) error {
	jobInfoHandler := &JobInfoHandler{db: a.db, jobId: job.ID}
	r := a.settings.Runner
	tag := fmt.Sprintf("%s%s", interpreterVersion, r.TagSuffix)
	dockerJob := core.NewEvalCodeJob(
		r.Image,
		tag,
		r.Shell,
		r.Command,
		job.Code,
		jobInfoHandler,
		jobInfoHandler,
		jobInfoHandler,
	)
	return a.queue.Enqueue(dockerJob)
}