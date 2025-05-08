package moon

import (
	"fmt"

	"github.com/NicoNex/echotron/v3"
)

func photo(msg *echotron.Update) {
	AddChat(&msg.Message.Chat)
	AddUser(msg.Message.From)

	if Cfg.LogChats {
		if msg.Message.Text != "" {
			ChatLog(
				msg.Message.Chat.ID,
				fmt.Sprintf(
					"User %s %s (username = %s, id = %d) sent %d photos to %s %s (%d) with text: %s",
					msg.Message.From.FirstName,
					msg.Message.From.LastName,
					msg.Message.From.Username,
					msg.Message.From.ID,
					len(msg.Message.Photo),
					msg.Message.Chat.Type,
					msg.Message.Chat.Title,
					msg.Message.Chat.ID,
					msg.Message.Text,
				),
			)
		} else {
			ChatLog(
				msg.Message.Chat.ID,
				fmt.Sprintf(
					"User %s %s (username = %s, id = %d) sent %d photos to %s %s (%d)",
					msg.Message.From.FirstName,
					msg.Message.From.LastName,
					msg.Message.From.Username,
					msg.Message.From.ID,
					len(msg.Message.Photo),
					msg.Message.Chat.Type,
					msg.Message.Chat.Title,
					msg.Message.Chat.ID,
				),
			)
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
