package moon

import (
	"time"

	"bunny-bot/internal/elog"

	"github.com/NicoNex/echotron/v3"
)

// textMessage handles plain text messages. It does not moderate anything; it
// only records the message so the greeting sweeper can tell that a newcomer
// answered their greeting.
func textMessage(msg *echotron.Update) {
	if msg.Message == nil || msg.Message.Text == "" {
		return
	}

	// Remember the message so the greeting sweeper can find the answer to a
	// greeting.
	GreetAnswers.Set(time.Now().UnixMicro(), msg.Message, greetCacheTTL)

	// And record it in the chat's own log, which outlives the cache.
	elog.LogMessage(msg.Message)
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
