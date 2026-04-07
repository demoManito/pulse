package main

import (
	"flag"
	"os"

	"github.com/demoManito/pulse/cmd/http"
	"github.com/demoManito/pulse/config"
	"github.com/demoManito/pulse/internal/service"
	"github.com/demoManito/pulse/pkg/logger"
	"github.com/demoManito/pulse/pkg/logger/logrus"
)

// Check Command interface
var (
	_ Command = (*http.Command)(nil)
)

// Command interface
type Command interface {
	Name() string
	Run(conf *config.Config) error
}

func commandMap() map[string]Command {
	cmds := []Command{
		&http.Command{},
	}
	commandMap := make(map[string]Command, len(cmds))
	for _, cmd := range cmds {
		commandMap[cmd.Name()] = cmd
	}
	return commandMap
}

func main() {
	cfg, err := config.LoadConfig(*flag.String("config", "./config/config.test.yaml", "config file"))
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	logger.SetDefault(logrus.New(logrus.WithLevel(logger.InfoLevel)))

	err = service.Init(cfg)
	if err != nil {
		logger.Fatalf("Failed to init service: %v", err)
	}
	defer service.Close()

	cmd, ok := commandMap()[os.Args[1]]
	if !ok {
		logger.Fatalf("Unknown command: %s", os.Args[1])
	}
	err = cmd.Run(cfg)
	if err != nil {
		logger.Fatal(err)
	}
}
