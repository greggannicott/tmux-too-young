package main

import (
	"fmt"
	"os"
)

type launchableDir struct {
	fullPath string
}

func main() {
	launchableDirs := getDirs()
	for _, dir := range launchableDirs {
		fmt.Println(dir.getFriendlyName())
	}
}

func getDirs() []launchableDir {
	launchableDirs := []launchableDir{}
	// The following is hard coded but eventually will be obtained via a loop over a config entry.
	currentDir := "/Users/greggannicott/code/"
	dirs, _ := os.ReadDir(currentDir)
	for _, dir := range dirs {
		fullPath := currentDir + dir.Name()
		_, err := os.Stat(fullPath + "/.git/")
		if err == nil {
			launchableDir := launchableDir{
				fullPath: fullPath,
			}
			launchableDirs = append(launchableDirs, launchableDir)
		}
	}
	return launchableDirs
}

func (l launchableDir) getFriendlyName() string {
	return l.fullPath
}
