package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func Exec(cmdStr string) (string, error) {
	parts := strings.Fields(cmdStr)
	cmd := exec.Command(parts[0], parts[1:]...)

	// 使用 bytes.Buffer 来捕获输出
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// 执行命令
	err := cmd.Run()
	if err != nil {
		fmt.Println("命令执行错误:", err)
		fmt.Println("错误输出:", errBuf.String())
		return errBuf.String(), err
	}

	return outBuf.String(), nil
}
