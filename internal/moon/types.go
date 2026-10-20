package moon

import (
	"github.com/NicoNex/echotron/v3"
	"github.com/akyoto/cache"
)

// MyConfig describes bot config file format.
type MyConfig struct {
	Token    string `json:"token,omitempty"`
	LogLevel string `json:"loglevel,omitempty"`
	LogFile  string `json:"logfile,omitempty"`

	// ChatLogs is the directory for per-chat, per-day chat event logs. Empty
	// disables chat logging. If set, the directory must already exist: a path
	// that does not is a typo, and silently logging nowhere is the worst
	// possible response to one.
	ChatLogs string `json:"chat_logs,omitempty"`
}

// tgClient is the subset of the Telegram API that bunny-bot actually calls.
// It exists so tests can substitute a fake; echotron.API satisfies it.
type tgClient interface {
	SendMessage(text string, chatID int64, opts *echotron.MessageOptions) (echotron.APIResponseMessage, error)
	BanChatMember(chatID int64, userID int64, opts *echotron.BanOptions) (echotron.APIResponseBool, error)
	DeleteMessage(chatID int64, messageID int) (echotron.APIResponseBase, error)
}

// chatUser identifies a user's membership of one specific chat. The same user is
// a different identity in every chat: being silent in one says nothing about
// another, and answering in one must not rescue them in another.
type chatUser struct {
	chatID int64
	userID int64
}

// banOutcome is what a squash attempt actually did. The three cases are distinct
// on purpose: a duplicate ban is a success, not a failure, and reporting it as
// one makes a working ban look broken in the logs.
type banOutcome int

// eventCache is a typed wrapper over akyoto/cache.Cache that stores values of
// type T keyed by a Unix timestamp in microseconds. It hides the underlying
// sync.Map so callers never deal with interface{} keys/values or per-access type
// conversions.
type eventCache[T any] struct {
	cache *cache.Cache
}

// joinRef is one stale join, kept together with everything the verdict and the
// retraction need so neither has to look anything up.
type joinRef struct {
	micros int64

	// chat and user come from the join event itself, so a ban can be logged with
	// the newcomer's real name and the chat's real title instead of a cache
	// lookup that may have expired or never been populated.
	chat echotron.Chat
	user *echotron.User

	// greetingID is the question we asked them, deleted once they are banned for
	// never answering it. Zero means we never managed to ask.
	greetingID int
}

// pendingGreeting is one recorded join: the message that announced it, and the
// id of the greeting we sent in reply.
//
// The greeting id lives here rather than in a second cache, so that recording,
// forgetting and retracting are all one lookup on one key. A zero greeting id
// means "no greeting was ever delivered for this join" — never a failure.
type pendingGreeting struct {
	join       *echotron.Message
	greetingID int
}

// retractOutcome is what a retraction pass actually achieved. The three cases are
// distinct because "nothing to delete" and "tried and failed" call for different
// log lines, and only the second is worth an operator's attention.
type retractOutcome int

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
