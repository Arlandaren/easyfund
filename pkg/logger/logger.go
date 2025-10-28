package logger

import (
	"go.uber.org/zap"
)

func New() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	return config.Build()
}
