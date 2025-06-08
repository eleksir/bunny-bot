package moon

import (
	"math/rand"
	"os"
	"time"

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
	SquashedMembers = cache.New(1 * time.Minute)

	// NewMembers struct with new members appeared via message with new_chat_members field set.
	NewMembers = cache.New(1 * time.Minute)

	// AppearedMembers struct with new members appeared via ChatMemberUpdated event.
	AppearedMembers = cache.New(60 * time.Minute)

	// PendingCASMembers contains just arrived members that pass first cas check, but still waiting second cas check.
	PendingCASMembers = cache.New(60 * time.Minute)

	// GreetChatMembers contains just arrived members that does not look suspicious by other means.
	GreetChatMembers = cache.New(60 * time.Minute)

	ChatList = make([]int64, 0)

	GreetMessages = []string{
		"Дратути",
		"Дарована",
		"Доброе утро, день или вечер",
		"Добро пожаловать в наше скромное коммунити",
		"Наше вам с кисточкой тут, на канальчике",
	}

	// Random number
	Random *rand.Rand
)

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
