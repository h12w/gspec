// Package atexit add the missing atexit functionality in the testing package.
package atexit

import (
	"os"
	"os/exec"
	"sync"
	"syscall"
)

const withFinalArg = "WITH-FINALIZATION---"

var once sync.Once

// Do starts another process of the same program with the same argument and calls
// f to do clean up job after the process is terminated.
// Do should be called within an init function, and can be called only once in a program.
// A second call to Do will be ignored.
func Do(f func()) {
	once.Do(func() {
		do(f)
	})
}

func do(f func()) {
	for _, arg := range os.Args {
		if arg == withFinalArg {
			return
		}
	}
	exitCode := runSelf()
	f()
	os.Exit(exitCode)
}

func runSelf() int {
	if err := (&exec.Cmd{
		Path:   os.Args[0],
		Args:   append(os.Args[1:], withFinalArg),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}).Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.Sys().(syscall.WaitStatus).ExitStatus()
		}
		return -1
	}
	return 0
}
