package moon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"bunny-bot/internal/elog"
	"bunny-bot/internal/mylog"

	"github.com/hjson/hjson-go"
)

// parseConfig разбирает и валидирует даденный конфиг.
func ParseConfig(path string) (MyConfig, error) {
	fileInfo, err := os.Stat(path)

	// Assume that no such file or directory error occurs.
	if err != nil {
		return MyConfig{}, fmt.Errorf("unable to find config: %w", err)
	}

	// Config file looks too long for config.
	if fileInfo.Size() > 65535 {
		err := fmt.Errorf("config file %s is too long for config, skipping", path) //nolint: err113

		return MyConfig{}, err
	}

	buf, err := os.ReadFile(path)

	// Noluck reading file.
	if err != nil {
		return MyConfig{}, err
	}

	/*
	 * According to docs hjson kinda can patse quirky json, but unmarshalls it to map.                            *
	 * Struct on output look more interesting: so we unmarshall config to map, then mashall back to json and then *
	 * unmarshall it via ecodinf/json to struct. Not very effective way but we do this not too often.             *
	 */
	var (
		sampleConfig MyConfig
		tmp          map[string]any
	)

	err = hjson.Unmarshal(buf, &tmp)

	// Не удалось распарсить.
	if err != nil {
		return MyConfig{}, err
	}

	tmpjson, err := json.Marshal(tmp)

	// Не удалось преобразовать map-ку в json.
	if err != nil {
		return MyConfig{}, err
	}

	if err := json.Unmarshal(tmpjson, &sampleConfig); err != nil {
		return MyConfig{}, err
	}

	if sampleConfig.Token == "" {
		return MyConfig{}, fmt.Errorf("token field in config file %s must be set", path) //nolint: err113
	}

	if _, err := mylog.ParseLevel(sampleConfig.LogLevel); err != nil {
		sampleConfig.LogLevel = "info"
	}

	// Chat logs are optional, but a configured path that does not exist is a
	// typo and must stop startup rather than silently disable logging.
	//
	// A relative chat_logs is resolved against the directory holding the config
	// file, not the process working directory: a daemon's cwd is wherever its
	// supervisor happened to start it, so a relative path there would move the
	// log directory on every restart under a different supervisor.
	if sampleConfig.ChatLogs != "" && !filepath.IsAbs(sampleConfig.ChatLogs) {
		sampleConfig.ChatLogs = filepath.Join(filepath.Dir(path), sampleConfig.ChatLogs)
	}

	if err := elog.Init(sampleConfig.ChatLogs); err != nil {
		return MyConfig{}, fmt.Errorf("chat_logs field in config file %s is invalid: %w", path, err)
	}

	// sampleConfig.Log = "" if not set

	return sampleConfig, err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
