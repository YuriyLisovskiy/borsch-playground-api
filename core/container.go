package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

type StdWriter interface {
	Write(out string)
}

type DockerContainer struct {
	Image        string
	Tag          string
	StdoutWriter StdWriter
	StderrWriter StdWriter

	cmd *exec.Cmd
}

func (dc *DockerContainer) Run(args ...string) (int, error) {
	defer func() {
		dc.cmd = nil
	}()

	fullImage := dc.Image + ":" + dc.Tag
	if dc.cmd != nil {
		return -1, errors.New(fmt.Sprintf("Container %s is already running", fullImage))
	}

	dc.cmd = exec.Command("docker", append([]string{"run", "--rm", fullImage}, args...)...)

	stdoutReader, err := dc.cmd.StdoutPipe()
	if err != nil {
		return -1, err
	}

	stdoutDone, err := dc.makeStdScanner(stdoutReader, dc.StdoutWriter)
	if err != nil {
		return -1, err
	}

	stderrReader, err := dc.cmd.StderrPipe()
	if err != nil {
		return -1, err
	}

	stderrDone, err := dc.makeStdScanner(stderrReader, dc.StderrWriter)
	if err != nil {
		return -1, err
	}

	if err := dc.cmd.Start(); err != nil {
		return -1, err
	}

	<-stdoutDone
	<-stderrDone

	err = dc.cmd.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), err
		}

		return -1, err
	}

	return 0, nil
}

func (dc *DockerContainer) makeStdScanner(reader io.ReadCloser, writer StdWriter) (<-chan bool, error) {
	scanner := bufio.NewScanner(reader)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			writer.Write(scanner.Text())
		}

		done <- true
	}()

	return done, nil
}

func ListDockerContainers() error {
	cmd := exec.Command("docker", "ps", "-a")
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
