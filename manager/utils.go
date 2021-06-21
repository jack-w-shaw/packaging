// Copyright 2015 Canonical Ltd.
// Copyright 2015 Cloudbase Solutions SRL
// Licensed under the LGPLv3, see LICENCE file for details.

package manager

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/retry"
)

var (
	logger = loggo.GetLogger("juju.packaging.manager")

	// Override for testing.
	Delay    = 10 * time.Second
	Attempts = 30
)

// CommandOutput is cmd.Output. It was aliased for testing purposes.
var CommandOutput = (*exec.Cmd).CombinedOutput

// ProcessStateSys is ps.Sys. It was aliased for testing purposes.
var ProcessStateSys = (*os.ProcessState).Sys

// RunCommand is helper function to execute the command and gather the output.
var RunCommand = func(command string, args ...string) (output string, err error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// exitStatuser is a mini-interface for the ExitStatus() method.
type exitStatuser interface {
	ExitStatus() int
}

// DNSRetryableError is a high order function that checks the status code from
// both apt and yum, which return 100 on abnormal execution due to outside
// issues (example: momentary dns failure).
func DNSRetryableError(err error, code int, output string) (bool, error) {
	if code != 100 {
		return false, errors.Trace(err)
	}
	return true, nil
}

// FatalError checks to see if a fatal error has occurred in the output of
// the resulting command.
func FatalError(f func(string) error) func(error, int, string) (bool, error) {
	return func(err error, code int, output string) (bool, error) {
		if retryable, retryableErr := DNSRetryableError(err, code, output); retryableErr != nil || !retryable {
			return false, errors.Trace(err)
		}

		err = f(output)
		return err == nil, errors.Annotatef(err, "encountered fatal error")
	}
}

// RunCommandWithRetry is a helper function which tries to execute the given command.
// It tries to do so for 30 times with a 10 second sleep between commands.
// It returns the output of the command, the exit code, and an error, if one occurs,
// logging along the way.
// It was aliased for testing purposes.
var RunCommandWithRetry = func(cmd string, retryError func(error, int, string) (bool, error)) (output string, code int, _ error) {
	// split the command for use with exec
	args := strings.Fields(cmd)
	if len(args) <= 1 {
		return "", 1, errors.New(fmt.Sprintf("too few arguments: expected at least 2, got %d", len(args)))
	}
	if retryError == nil {
		retryError = DNSRetryableError
	}

	logger.Infof("Running: %s", cmd)

	// Retry operation 30 times, sleeping every 10 seconds between attempts.
	// This avoids failure in the case of something else having the dpkg lock
	// (e.g. a charm on the machine we're deploying containers to).
	var out []byte
	tryAgain := false
	err := retry.Call(retry.CallArgs{
		Clock:    clock.WallClock,
		Delay:    Delay,
		Attempts: Attempts,
		NotifyFunc: func(lastError error, attempt int) {
			logger.Infof("Retrying: %s", cmd)
		},
		IsFatalError: func(err error) bool {
			return !tryAgain
		},
		Func: func() error {
			tryAgain = false
			// Create the command for each attempt, because we need to
			// call cmd.CombinedOutput only once. See http://pad.lv/1394524.
			command := exec.Command(args[0], args[1:]...)

			var err error
			out, err = CommandOutput(command)
			if err == nil {
				return nil
			}

			exitError, ok := err.(*exec.ExitError)
			if !ok {
				return errors.Annotatef(err, "unexpected error type %T", err)
			}
			waitStatus, ok := ProcessStateSys(exitError.ProcessState).(exitStatuser)
			if !ok {
				return errors.Annotatef(err, "unexpected process state type %T", exitError.ProcessState.Sys())
			}

			if retryable, retryableErr := retryError(err, waitStatus.ExitStatus(), string(out)); retryableErr != nil {
				return errors.Trace(retryableErr)
			} else if !retryable {
				return errors.Trace(err)
			}

			tryAgain = true
			return err
		},
	})

	if err != nil {
		logger.Errorf("packaging command failed: %v; cmd: %q; output: %s",
			err, cmd, string(out))
		return string(out), code, errors.Errorf("packaging command failed: %v", err)
	}

	return string(out), 0, nil
}
