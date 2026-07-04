package termi

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func getLockPath(cfgDir string) string {
	return filepath.Join(cfgDir, "lock")
}

func Lock(cfgDir string) error {
	path := getLockPath(cfgDir)
	for i := 0; i < 8; i++ {
		err := os.Mkdir(path, 0777)
		if err == nil {
			return nil
		}
		d, _ := time.ParseDuration("1s")
		time.Sleep(d)
	}
	return fmt.Errorf("cannot create lock")
}

func Unlock(cfgDir string) error {
	path := getLockPath(cfgDir)
	return os.Remove(path)
}
