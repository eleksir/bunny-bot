package moon

import (
	"fmt"
	"log"
	"os"
)

func ChatLog(file, text string) {
	var (
		isOpened = false
	)

	for _, l := range chatlogs {
		if l.Name == file {
			isOpened = true
			break
		}
	}

	if !isOpened {
		// open log file
		logFile, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			fmt.Printf("Unable to open log file %s: %s", file, err)
		}

		l := Chatlog{Name: file, File: logFile}

		chatlogs = append(chatlogs, l)
	}

	for _, l := range chatlogs {
		if l.Name == file {
			log.SetOutput(l.File)
			// optional: log date-time, filename, and line number
			log.SetFlags(log.Lshortfile | log.LstdFlags)
			log.Println(text)
			break
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
