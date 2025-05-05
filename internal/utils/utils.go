package utils

import (
	"log"
)

var debugMode = false

func Initialize(isDebug bool) {
	debugMode = isDebug

	DebugLog("DebugMode:", debugMode)
}

func DebugLog(format string, v ...interface{}) {
	if debugMode {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func InfoLog(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}
