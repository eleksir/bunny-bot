package moon

import (
	"github.com/NicoNex/echotron/v3"
)

func textMessage(msg *echotron.Update) {
	if Cfg.LogChats {
		if msg.Message.ReplyToMessage == nil {
			// Just log it, it seems later here should be ban/mute command set.
			Log.Infof(
				"User %s %s (username = %s, id = %d) in %s %s (%d) says: %s",
				msg.Message.From.FirstName,
				msg.Message.From.LastName,
				msg.Message.From.Username,
				msg.Message.From.ID,
				msg.Message.Chat.Type,
				msg.Message.Chat.Title,
				msg.Message.Chat.ID,
				msg.Message.Text,
			)
		} else {
			// Just log it.
			Log.Infof(
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
			)
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
