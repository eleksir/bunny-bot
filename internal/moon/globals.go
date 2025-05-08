package moon

import (
	"os"

	"github.com/NicoNex/echotron/v3"
	"github.com/cockroachdb/pebble"
	"go.uber.org/zap"
)

var (
	Cfg     MyConfig
	SigChan = make(chan os.Signal, 1)
	Log     *zap.SugaredLogger
	Tg      echotron.API

	// Мапка с открытыми дескрипторами баз настроек.
	DB = make(map[string]*pebble.DB)
)

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
