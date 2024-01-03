package util

import (
	"os"
	"strconv"
)

func WritePid(pidPath string) error {
	pid := os.Getpid()
	pidStr := strconv.Itoa(pid)
	if err := os.WriteFile(pidPath, []byte(pidStr), 0755); err != nil {
		return err
	}
	return nil
}
