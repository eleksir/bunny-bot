package moon

import (
	"fmt"

	"github.com/NicoNex/echotron/v3"
)

// animation parses message with animation.
func animation(msg *echotron.Update) {
	if Cfg.LogChats {
		ChatLog(
			msg.Message.Chat.ID,
			fmt.Sprintf(
				"User %s %s (username = %s, id = %d) sent animation to chat %s %s (%d)",
				msg.Message.From.FirstName,
				msg.Message.From.LastName,
				msg.Message.From.Username,
				msg.Message.From.ID,
				msg.Message.Chat.Type,
				msg.Message.Chat.Title,
				msg.Message.Chat.ID,
			),
		)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
