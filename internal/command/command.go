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
	runner.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel)
	defer signal.Stop(signalChannel)

	if startErr := runner.Start(); startErr != nil {
		return startErr
	}

	pgid := getProcessGroupID(runner.Process.Pid)
	done := make(chan struct{})

	go func() {
		receivedTermination := false
		for {
			select {
			case sig := <-signalChannel:
				if isTermSignal(sig) {
					if receivedTermination {
						forwardSignalToGroup(runner.Process, pgid, syscall.SIGKILL)
						return
					}
					receivedTermination = true
				}
				forwardSignalToGroup(runner.Process, pgid, sig)
			case <-done:
				return
			}
		}
	}()

	waitErr := runner.Wait()
	close(done)

	return waitErr
}

func forwardSignalToGroup(proc *os.Process, pgid int, sig os.Signal) {
	if s, ok := sig.(syscall.Signal); ok {
		_ = syscall.Kill(-pgid, s)
	}
}

func getProcessGroupID(pid int) int {
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		return pid
	}
	return pgid
}

func isTermSignal(sig os.Signal) bool {
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP:
		return true
	}
	return false
}
