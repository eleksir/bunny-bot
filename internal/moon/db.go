package moon

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/cockroachdb/errors/oserror"
	"github.com/cockroachdb/pebble"
)

// StoreKV сохраняет в указанной бд ключ и значение.
func StoreKV(db *pebble.DB, key string, value string) error {
	var (
		kArray = []byte(key)
		vArray = []byte(value)
	)

	err := db.Set(kArray, vArray, pebble.Sync)

	if err != nil {
		return err
	}

	err = db.Flush()

	return err
}

// FetchV достаёт значение по ключу.
func FetchV(db *pebble.DB, key string) (string, error) {
	var (
		kArray = []byte(key)

		vArray      []byte
		valueString = ""
	)

	vArray, closer, err := db.Get(kArray)

	if err != nil {
		return valueString, err
	}

	valueString = string(vArray)
	err = closer.Close()

	return valueString, err
}

// GetValue достаёт настройку из БД с настройками.
// TODO: should return string, error and should not log error internally.
func (cfg *MyConfig) GetValue(dataBaseName string, chatID string, setting string) string {
	var err error

	chatHash := sha256.Sum256([]byte(chatID))
	dataBaseDir := fmt.Sprintf("db/%s/%x", dataBaseName, chatHash)

	// Если БД не открыта, откроем её
	if _, ok := DB[dataBaseDir]; !ok {
		var options pebble.Options
		// По дефолту ограничение ставится на мегабайты данных, а не на количество файлов, поэтому с дефолтными
		// настройками порождается огромное количество файлов. Умолчальное ограничение на количество файлов - 500 штук,
		// что нас не устраивает, поэтому немного снизим эту цифру до более приемлемых значений.
		options.L0CompactionFileThreshold = 8

		if err := os.MkdirAll(dataBaseDir, os.ModePerm); err != nil {
			Log.Errorf("Unable to create db dir %s: %s", dataBaseDir, err)

			return ""
		}

		DB[dataBaseDir], err = pebble.Open(dataBaseDir, &options)

		if err != nil {
			Log.Errorf("Unable to open db, %s: %s\n", dataBaseDir, err)

			return ""
		}
	}

	value, err := FetchV(DB[dataBaseDir], setting)

	// Если из базы ничего не вынулось, по каким-то причинам, то просто вернём пустую строку.
	if err != nil {
		switch {
		case errors.Is(err, pebble.ErrNotFound):
			Log.Debugf("Unable to get value for %s: no record found in db %s", setting, dataBaseDir)
		case errors.Is(err, fs.ErrNotExist):
			Log.Debugf("Unable to get value for %s: db dir %s does not exist", setting, dataBaseDir)
		case errors.Is(err, oserror.ErrNotExist):
			Log.Debugf("Unable to get value for %s: db dir %s does not exist", setting, dataBaseDir)
		default:
			Log.Errorf("Unable to get value for %s in db dir %s: %s", setting, dataBaseDir, err)
		}

		return ""
	}

	return value
}

// SaveSetting сохраняет настройку в БД с настройками.
func (cfg *MyConfig) SaveKeyValue(dataBaseName string, chatID string, setting string, value string) error {
	var (
		chatHash    = sha256.Sum256([]byte(chatID))
		dataBaseDir = fmt.Sprintf("%s/db/%s/%x", cfg.DataDir, dataBaseName, chatHash)
		err         error
	)

	// Если БД не открыта, откроем её.
	if _, ok := DB[dataBaseDir]; !ok {
		var options pebble.Options
		// По дефолту ограничение ставится на мегабайты данных, а не на количество файлов, поэтому с дефолтными
		// настройками порождается огромное количество файлов. Умолчальное ограничение на количество файлов - 500 штук,
		// что нас не устраивает, поэтому немного снизим эту цифру до более приемлемых значений.
		options.L0CompactionFileThreshold = 8

		if err := os.MkdirAll(dataBaseDir, os.ModePerm); err != nil {
			return fmt.Errorf("unable to create db dir %s: %w", dataBaseDir, err)
		}

		DB[dataBaseDir], err = pebble.Open(dataBaseDir, &options)

		if err != nil {
			return fmt.Errorf("unable to open settings db, %s: %w", dataBaseDir, err)
		}
	}

	if err := StoreKV(DB[dataBaseDir], setting, value); err != nil {
		return fmt.Errorf("unable to save setting %s in database %s: %w", setting, dataBaseDir, err)
	}

	return nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
