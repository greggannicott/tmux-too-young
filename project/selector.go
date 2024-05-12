package project

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetSelectionFromUser(initialSearchTerm string) project {
	var input string
	for _, choice := range projects {
		input += choice.getFriendlyName() + "\n"
	}
	cmd := exec.Command("fzf-tmux", "-p", "--cycle", "--reverse", "--border", "--info=inline-right", "--header=Select a Project to open in tmux:", "--query="+initialSearchTerm, "-1")
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		// If it exited with code 130, then the user pressed Ctrl-C and we don't want to show an error
		if err.Error() == "exit status 130" {
			os.Exit(0)
		} else {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}
	selectedProjectName := strings.TrimSpace(stdout.String())
	return findProjectDirectoryByFriendlyName(selectedProjectName)
}

func findProjectDirectoryByFriendlyName(name string) project {
	var matchingProjectDirectory project
	for _, p := range projects {
		if p.getFriendlyName() == name {
			matchingProjectDirectory = p
			break
		}
	}
	return matchingProjectDirectory
}
