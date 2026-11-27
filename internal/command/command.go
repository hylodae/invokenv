package command

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func RunCommand(command []string, env []string) error {
	binary, lookupErr := exec.LookPath(command[0])
	if lookupErr != nil {
		return lookupErr
	}

	runner := exec.Command(binary, command[1:]...)
	runner.Stdin = os.Stdin
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	runner.Env = env

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel)

	if startErr := runner.Start(); startErr != nil {
		return startErr
	}

	// Forward all signals to the child process
	go func() {
		for {
			receivedSignal := <-signalChannel
			runner.Process.Signal(receivedSignal) // nolint:errcheck
		}
	}()

	if waitErr := runner.Wait(); waitErr != nil {
		runner.Process.Signal(syscall.SIGKILL) // nolint:errcheck
		return waitErr
	}

	return nil
}
