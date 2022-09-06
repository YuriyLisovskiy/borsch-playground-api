package core

import (
	"errors"
	"fmt"
	"strings"
)

type JobInfoHandler interface {
	OnError(err error)
	OnExit(exitCode int, err error)
}

type EvalCodeJob struct {
	container   *DockerContainer
	shell       string
	command     string
	code        string
	infoHandler JobInfoHandler
}

func NewEvalCodeJob(
	image, tag, shell, command, code string,
	outWriter, errWriter StdWriter,
	infoHandler JobInfoHandler,
) *EvalCodeJob {
	return &EvalCodeJob{
		container: &DockerContainer{
			Image:        image,
			Tag:          tag,
			StdoutWriter: outWriter,
			StderrWriter: errWriter,
		},
		shell:       shell,
		command:     command,
		code:        code,
		infoHandler: infoHandler,
	}
}

func (j *EvalCodeJob) Run() {
	if j.container == nil {
		j.infoHandler.OnError(errors.New("docker container is nil"))
		return
	}

	code := strings.ReplaceAll(j.code, "\"", "\\\"")
	shellScript := strings.ReplaceAll(j.command, "<code>", fmt.Sprintf("\"%s\"", code))
	exitCode, err := j.container.Run(j.shell, "-c", shellScript)
	if j.infoHandler != nil {
		j.infoHandler.OnExit(exitCode, err)
	}
}
