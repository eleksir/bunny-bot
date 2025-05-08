package moon

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ChatLog logs message to separate file, placed in basedir(Cfg.LogFile).
func ChatLog(chatid int64, message string) {
	var (
		filename string
	)

	if Cfg.LogFile != "" {
		filename = fmt.Sprintf("%s/%s", filepath.Dir(Cfg.LogFile), strconv.FormatInt(chatid, 10))
	} else {
		filename = os.DevNull
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		fmt.Printf("Unable to open file %s for logging: %s", filename, err)

		// TODO: Should we try harder?
		return
	}

	defer f.Close()

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:        f,
			NoColor:    true,
			TimeFormat: time.DateTime,
			PartsOrder: []string{"time", "message"},
		},
	)

	log.Info().Msg(message)
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
