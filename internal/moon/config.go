package moon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hjson/hjson-go"
	"go.uber.org/zap/zapcore"
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

	if _, err := zapcore.ParseLevel(sampleConfig.LogLevel); err != nil {
		sampleConfig.LogLevel = "info"
	}

	// sampleConfig.Log = "" if not set

	if sampleConfig.CSign == "" {
		err := fmt.Errorf("csign field in config file %s must be set", path) //nolint: err113

		return MyConfig{}, err
	}

	if sampleConfig.DataDir == "" {
		return MyConfig{}, fmt.Errorf("data_dir field in config file %s must be set", path) //nolint: err113
	}

	if strings.Split(sampleConfig.DataDir, "")[0] != "/" {
		executablePath, err := os.Executable()

		if err != nil {
			err := fmt.Errorf("unable to get current executable path: %w", err)

			return MyConfig{}, err
		}

		sampleConfig.DataDir = fmt.Sprintf("%s/%s", filepath.Dir(executablePath), sampleConfig.DataDir)
	}

	return sampleConfig, err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
