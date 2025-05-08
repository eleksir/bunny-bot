package moon

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hjson/hjson-go"
	"go.uber.org/zap/zapcore"
)

// parseConfig разбирает и валидирует даденный конфиг.
func ParseConfig(path string) (MyConfig, error) {
	fileInfo, err := os.Stat(path)

	// Предполагаем, что файла либо нет, либо мы не можем его прочитать, второе надо бы логгировать, но пока забьём.
	if err != nil {
		return MyConfig{}, err
	}

	// Конфиг-файл длинноват для конфига, попробуем следующего кандидата.
	if fileInfo.Size() > 65535 {
		err := fmt.Errorf("config file %s is too long for config, skipping", path) //nolint: err113

		return MyConfig{}, err
	}

	buf, err := os.ReadFile(path)

	// Не удалось прочитать.
	if err != nil {
		return MyConfig{}, err
	}

	// Исходя из документации, hjson какбы умеет парсить "кривой" json, но парсит его в map-ку.
	// Интереснее на выходе получить структурку: то есть мы вначале конфиг преобразуем в map-ку, затем эту map-ку
	// сериализуем в json, а потом json превращааем в структурку. Не очень эффективно, но он и не часто требуется.
	var (
		sampleConfig MyConfig
		tmp          map[string]interface{}
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

	return sampleConfig, err
}
