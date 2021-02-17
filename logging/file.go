package logging

import (
	"fmt"
	"time"
)

// getLogFilePath get the log file save path
func getLogFilePath(runtimeRootPath, logSavePath string) string {
	return fmt.Sprintf("%s%s", runtimeRootPath, logSavePath)
}

// getLogFileName get the save name of the log file
func getLogFileName(logSaveName, logTimeFormat, logFileExtension string) string {
	return fmt.Sprintf("%s%s.%s",
		logSaveName,
		time.Now().Format(logTimeFormat),
		logFileExtension,
	)
}
