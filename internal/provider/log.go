package provider

import "log"

// providerLogger is a logger that logs messages accordingly to the terraform-sdk requirements.
type providerLogger struct {
	logger *log.Logger
}

var logger = providerLogger{
	logger: log.Default(),
}

func (l *providerLogger) Infof(format string, v ...interface{}) {
	l.logger.Printf("[INFO] "+format, v...)
}

func (l *providerLogger) Debugf(format string, v ...interface{}) {
	l.logger.Printf("[DEBUG] "+format, v...)
}

func (l *providerLogger) Errorf(format string, v ...interface{}) {
	l.logger.Printf("[ERROR] "+format, v...)
}

func (l *providerLogger) Warnf(format string, v ...interface{}) {
	l.logger.Printf("[WARN] "+format, v...)
}
