package read

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-task/task/v2/internal/taskfile"
)

var (
	ErrIncludedTaskfilesCantHaveIncludes = errors.New(
		`task: Included Taskfiles can't have includes. 
		Please, move the include to the main Taskfile`,
	)

	ErrNoTaskfileFound = errors.New(
		`task: No Taskfile.yml found. Use "task --init" to create a new one`,
	)
)

// Taskfile reads a Taskfile for a given directory
func Taskfile(dir, location string) (*taskfile.Taskfile, error) {
	path := filepath.Join(dir, "Taskfile.yml")
	if location != "" {
		path = filepath.Join(dir, location)
	}
	if location != "" && filepath.IsAbs(location) {
		path = location
	}
	if _, err := os.Stat(path); err != nil {
		return nil, ErrNoTaskfileFound
	}
	t, err := taskfile.LoadFromPath(path)
	if err != nil {
		return nil, err
	}

	if err := t.ProcessIncludes(dir); err != nil {
		return nil, err
	}

	path = filepath.Join(dir, fmt.Sprintf("Taskfile_%s.yml", runtime.GOOS))
	if _, err = os.Stat(path); err == nil {
		osTaskfile, err := taskfile.LoadFromPath(path)
		if err != nil {
			return nil, err
		}
		if err = taskfile.Merge(nil, t, osTaskfile); err != nil {
			return nil, err
		}
	}

	for name, task := range t.Tasks {
		task.Task = name
	}

	return t, nil
}
