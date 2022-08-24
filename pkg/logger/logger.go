package logger

import "go.uber.org/zap"

func New() *zap.Logger {
	return zap.NewExample()
}
