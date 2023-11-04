package internal

import (
	"fmt"
)

func KillProcess() error {
	fmt.Println("Windows下不支持自动结束进程")
	return nil
}
