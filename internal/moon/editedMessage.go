package moon

import (
	"fmt"

	"github.com/NicoNex/echotron/v3"
)

// editedMessage handles editedMessage type of update.
func editedMessage(msg *echotron.Update) {
	AddUser(msg.EditedMessage.From)
	AddChat(&msg.EditedMessage.Chat)

	// Edited message can contain all kinds of edited attributes that telegram allows to change in message.
	// Not sure is we should handle all of them.
	if Cfg.LogChats {
		if msg.EditedMessage.Text != "" {
			ChatLog(
				msg.EditedMessage.Chat.ID,
				fmt.Sprintf(
					"User %s %s (username = %s, id = %d) in %s %s (%d) edited message: %s",
					msg.EditedMessage.From.FirstName,
					msg.EditedMessage.From.LastName,
					msg.EditedMessage.From.Username,
					msg.EditedMessage.From.ID,
					msg.EditedMessage.Chat.Type,
					msg.EditedMessage.Chat.Title,
					msg.EditedMessage.Chat.ID,
					msg.EditedMessage.Text,
				),
			)
		} else {
			ChatLog(
				msg.EditedMessage.Chat.ID,
				fmt.Sprintf(
					"User %s %s (username = %s, id = %d) in %s %s (%d) edited message",
					msg.EditedMessage.From.FirstName,
					msg.EditedMessage.From.LastName,
					msg.EditedMessage.From.Username,
					msg.EditedMessage.From.ID,
					msg.EditedMessage.Chat.Type,
					msg.EditedMessage.Chat.Title,
					msg.EditedMessage.Chat.ID,
				),
			)
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
