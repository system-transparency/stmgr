package eval

import "git.glasklar.is/system-transparency/core/stmgr/log"

// Helper function to map strings to log.logLevel.
func setLoglevel(level string) {
	switch level {
	case "debug":
		log.SetLoglevel(log.DebugLevel)
	case "info":
		log.SetLoglevel(log.InfoLevel)
	case "warn":
		log.SetLoglevel(log.WarnLevel)
	case "panic":
		log.SetLoglevel(log.PanicLevel)
	default:
		log.SetLoglevel(log.ErrorLevel)
	}
}
