package internal

import (
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

func KillProcessItself() error {
	//fmt.Println("Windows下不支持自动结束进程")
	//return nil

	//p, err := os.FindProcess(syscall.Getpid())
	//if err != nil {
	//	logrus.Error(err)
	//	return err
	//}
	//if err = p.Signal(syscall.SIGTERM); err != nil { // ERROR: not supported by windows. (Sending Interrupt on Windows is not implemented.)
	//if err = p.Signal(syscall.SIGKILL); err != nil { //如果用syscall.SIGKILL，进程会直接结束，不会触发stopAction回调。
	//	logrus.Error(err)
	//	return err
	//}

	{ // from: github.com/kardianos/service@v1.2.2/servicetest_windows_test.go
		dll, err := syscall.LoadDLL("kernel32.dll")
		if err != nil {
			logrus.Errorf("LoadDLL(\"kernel32.dll\") err: %s", err)
			return err
		}
		p, err := dll.FindProc("GenerateConsoleCtrlEvent")
		if err != nil {
			logrus.Errorf("FindProc(\"GenerateConsoleCtrlEvent\") err: %s", err)
			return err
		}
		// Send the CTRL_BREAK_EVENT to a console process group that shares
		// the console associated with the calling process.
		// https://msdn.microsoft.com/en-us/library/windows/desktop/ms683155(v=vs.85).aspx
		pid := os.Getpid()
		r1, _, err := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
		if r1 == 0 {
			logrus.Errorf("Call(CTRL_BREAK_EVENT, %d) err: %s", pid, err)
			return err
		}
	}
	return nil
}
