package moon

import (
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

	// PendingCASMembers contains newly comed members that pass first cas check, but still waiting second cas check.
	PendingCASMembers = cache.New(60 * time.Minute)

	// GreetPending contains people that should answer to bot greet.
	GreetPending = cache.New(1 * time.Minute)

	// NoMessagePeople contains people that are never say something in chat.
	NoMessagePeople = cache.New(1 * time.Minute)

	// People cache with users info.
	People = cache.New(8 * time.Hour)

	// Groups cache with groups info.
	Groups = cache.New(8 * time.Hour)

	ChatList = make([]int64, 0)
)

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
