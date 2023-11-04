package main

import (
	"fmt"
)

func killProcess() error {
	fmt.Println("Windows下不支持自动结束进程")
	return nil
}
