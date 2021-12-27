package zap_logger

import "go.uber.org/zap"

func GetLogger() *zap.Logger {
	return zap.NewExample()
}
