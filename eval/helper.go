package eval

import "system-transparency.org/stboot/stlog"

// Helper function to map strings to log.logLevel.
func setLoglevel(level string) {
	switch level {
	case "debug":
		stlog.SetLevel(stlog.DebugLevel)
	case "info":
		stlog.SetLevel(stlog.InfoLevel)
	case "warn":
		stlog.SetLevel(stlog.WarnLevel)
	default:
		stlog.SetLevel(stlog.InfoLevel)
	}
}
