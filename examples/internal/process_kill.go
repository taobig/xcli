//go:build !windows

package internal

import "syscall"

func KillProcess() error {
	err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM) //syscall.Kill不能在windows下使用
	return err
}
