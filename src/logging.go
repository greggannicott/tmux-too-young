package main

import (
	"fmt"
	"log/slog"
	"os"
)

func setupLogging() {
	homeDir, _ := os.UserHomeDir()
	dir := homeDir + "/Library/logs/"
	fileName := "tmux-too-young.log"
	path := dir + fileName
	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Print("Unable to create log file '"+path+"':", err)
		os.Exit(1)
	}
	logger = slog.New(slog.NewTextHandler(logFile, nil))
}
