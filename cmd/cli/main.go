// Package main is the entrypoint for the CLI executing test scenarios.
package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/inquiryproj/inquiry/internal/executor"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	wordPtr := flag.String("file", "", "the file name of your test scenario")
	v := flag.Bool("v", false, "verbose logging")
	flag.Parse()
	if *wordPtr == "" {
		logger.Error("file flag is required, provide as --flag <file.yaml>")
		return
	}
	if *v {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	scenarioName := *wordPtr
	f, err := os.Open(scenarioName)
	if err != nil {
		logger.Error("unable to open file", err)
		return
	}

	executorApp, err := executor.New(scenarioName,
		executor.WithReader(f),
		executor.WithLogger(logger),
	)
	if err != nil {
		logger.Error("unable to create test scenario executor", err)
		return
	}
	err = executorApp.Play()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Info("scenario executed successfully")
}
