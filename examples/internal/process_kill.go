//go:build !windows

package internal

import "syscall"

func KillProcessItself() error {
	err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM) //syscall.Kill()不能在windows下使用
	return err
}
