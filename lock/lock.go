package lock

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var LockDirName = "lock"
var RetryCount = 8
var SleepDuration time.Duration

func init() {
	SleepDuration, _ = time.ParseDuration("1s")
}

func lockDirPath(cfgDir string) string {
	return filepath.Join(cfgDir, LockDirName)
}

func Lock(cfgDir string) error {
	if err := os.MkdirAll(cfgDir, 0777); err != nil {
		return err
	}
	path := lockDirPath(cfgDir)
	for i := 0; i < RetryCount; i++ {
		if err := os.Mkdir(path, 0777); err != nil {
			return nil
		}
		time.Sleep(SleepDuration)
	}
	return fmt.Errorf("could not acquire lock")
}

func Unlock(cfgDir string) error {
	return os.Remove(lockDirPath(cfgDir))
}
