package moon

import (
	"os"

	"github.com/NicoNex/echotron/v3"
	"github.com/akyoto/cache"
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

	// SquashedMembers struct with banned members.
	SquashedMembers map[int64]*cache.Cache

	// NewMembers struct with new members appeared via message with new_chat_members field set.
	NewMembers map[int64]*cache.Cache

	// AppearedMembers struct with new members appeared via ChatMemberUpdated event.
	AppearedMembers map[int64]*cache.Cache

	ChatList = make([]int64, 0)
)

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
