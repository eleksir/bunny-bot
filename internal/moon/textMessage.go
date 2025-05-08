package moon

import (
	"fmt"

	"github.com/NicoNex/echotron/v3"
)

// textMessage handles text message events.
func textMessage(msg *echotron.Update) {
	AddChat(&msg.Message.Chat)
	AddUser(msg.Message.From)

	// Not sure if we must handle all kinds of these events.
	if Cfg.LogChats {
		switch {
		// Just log it, it seems later here should be ban/mute command set.
		case msg.Message != nil && msg.Message.ReplyToMessage == nil:
			ChatLog(
				msg.Message.Chat.ID,
				fmt.Sprintf(
					"User %s %s (username = %s, id = %d) in %s %s (%d) says: %s",
					msg.Message.From.FirstName,
					msg.Message.From.LastName,
					msg.Message.From.Username,
					msg.Message.From.ID,
					msg.Message.Chat.Type,
					msg.Message.Chat.Title,
					msg.Message.Chat.ID,
					msg.Message.Text,
				),
			)

		case msg.Message != nil && msg.Message.ReplyToMessage != nil:
			ChatLog(
				msg.Message.Chat.ID,
				fmt.Sprintf(
					"User %s %s (username = %s, id = %d) in %s %s (%d) replies to message %d by %s %s (username = %s, id = %d): %s",
					msg.Message.From.FirstName,
					msg.Message.From.LastName,
					msg.Message.From.Username,
					msg.Message.From.ID,
					msg.Message.Chat.Type,
					msg.Message.Chat.Title,
					msg.Message.Chat.ID,
					msg.Message.ReplyToMessage.ID,
					msg.Message.ReplyToMessage.From.FirstName,
					msg.Message.ReplyToMessage.From.LastName,
					msg.Message.ReplyToMessage.From.Username,
					msg.Message.ReplyToMessage.From.ID,
					msg.Message.Text,
				),
			)
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
