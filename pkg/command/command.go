package command

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

// CustomError is a custom error type that contains the command, stderr, and exit code.
type CMDError struct {
	Command       string
	Stderr        string
	ExitCode      int
	OriginalError error
}

// Error implements the error interface for CustomError.
func (e *CMDError) Error() string {
	return fmt.Sprintf(
		"command '%s' failed with exit code %d: %s (original error: %v)",
		e.Command,
		e.ExitCode,
		e.Stderr,
		e.OriginalError,
	)
}

func WithContext(ctx context.Context, name string, args []string, envVars map[string]string, workingDir string) ([]byte, []byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	// Set the working directory for the command
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	// Set custom environment variables
	if len(envVars) > 0 {
		env := os.Environ()
		for key, value := range envVars {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = env
	}

	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	makeError := func(err error) error {
		exitCode := 0
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
		return &CMDError{
			Command:       fmt.Sprintf("%s %v", name, args),
			Stderr:        stderr.String(),
			ExitCode:      exitCode,
			OriginalError: err,
		}
	}

	if err := cmd.Start(); err != nil {
		err = makeError(err)
		return stdout.Bytes(), stderr.Bytes(), err
	}

	if err := cmd.Wait(); err != nil {
		err = makeError(err)
		return stdout.Bytes(), stderr.Bytes(), err
	}

	return stdout.Bytes(), stderr.Bytes(), nil
}
