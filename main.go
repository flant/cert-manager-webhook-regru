package main

import (
	"log/slog"
	"os"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/flant/cert-manager-webhook-regru/regru"
	"github.com/go-logr/logr"
	logsapi "k8s.io/component-base/logs/api/v1"
)

// appVersion is injected at build time via -ldflags "-X main.appVersion=<version>".
var appVersion string

// slogJSONFactory implements logsapi.LogFormatFactory to provide a structured JSON logger
// based on the standard library slog package.
type slogJSONFactory struct{}

// Create returns a logr.Logger backed by slog's JSON handler, writing to the configured error stream.
func (slogJSONFactory) Create(_ logsapi.LoggingConfiguration, o logsapi.LoggingOptions) (logr.Logger, logsapi.RuntimeControl) {
	w := o.ErrorStream
	if w == nil {
		w = os.Stderr
	}
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return logr.FromSlogHandler(handler), logsapi.RuntimeControl{
		Flush: func() {},
	}
}

func init() {
	if err := logsapi.RegisterLogFormat("slog-json", slogJSONFactory{}, logsapi.LoggingStableOptions); err != nil {
		panic(err)
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	version := appVersion
	if version == "" {
		version = "dev"
	}
	logger.Info("starting cert-manager-webhook-regru", "version", version)

	groupName := os.Getenv("GROUP_NAME")
	if groupName == "" {
		logger.Error("GROUP_NAME environment variable is required")
		os.Exit(1)
	}

	username := os.Getenv("REGRU_USERNAME")
	if username == "" {
		logger.Error("REGRU_USERNAME environment variable is required")
		os.Exit(1)
	}

	password := os.Getenv("REGRU_PASSWORD")
	if password == "" {
		logger.Error("REGRU_PASSWORD environment variable is required")
		os.Exit(1)
	}

	os.Args = append(os.Args, "--logging-format=slog-json")

	regru.InitClient(username, password)

	logger.Info("running webhook server", "groupName", groupName)

	cmd.RunWebhookServer(groupName, &regru.Solver{})
}
