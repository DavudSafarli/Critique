package util

import (
	"runtime"
	"testing"
)

// GetDockerHost returns docker host uri based on OS
func GetDockerHost() string {
	if runtime.GOOS == "windows" {
		return "npipe:////./pipe/docker_engine"
	} else {
		return "unix:///var/run/docker.sock"
	}
}

// CreateCleanupWrapper returns a function to be passed to t.Cleanup() function
//
// example usage:
//
// t.Cleanup(util.CreateCleanupWrapper(t, cleanupFunc))
func CreateCleanupWrapper(t *testing.T, cleanupFunc func() error) func() {
	return func() {
		err := cleanupFunc()
		if err != nil {
			t.Errorf("error in cleanup function for %s: %s", t.Name(), err)
		}
	}
}
