package moon

import (
	"fmt"

	"github.com/NicoNex/echotron/v3"
)

func reaction(msg *echotron.Update) {
	AddChat(&msg.MessageReaction.Chat)
	AddUser(&msg.MessageReaction.User)

	if Cfg.LogChats {
		var reactionString string

		for _, item := range msg.MessageReaction.NewReaction {
			switch item.Type {
			case "Emoji":
				reactionString += item.Emoji
			case "CustomEmoji":
				reactionString += item.CustomEmoji
			}
		}

		if reactionString != "" {
			ChatLog(
				msg.MessageReaction.Chat.ID,
				fmt.Sprintf(
					"User %s %s (username = %s, id = %d) set reaction %s to chat %s %s (%d)",
					msg.MessageReaction.User.FirstName,
					msg.MessageReaction.User.LastName,
					msg.MessageReaction.User.Username,
					msg.MessageReaction.User.ID,
					reactionString,
					msg.MessageReaction.Chat.Type,
					msg.MessageReaction.Chat.Title,
					msg.MessageReaction.Chat.ID,
				),
			)
		} else {
			for _, item := range msg.MessageReaction.OldReaction {
				switch item.Type {
				case "Emoji":
					reactionString += item.Emoji
				case "CustomEmoji":
					reactionString += item.CustomEmoji
				}
			}

			ChatLog(
				msg.MessageReaction.Chat.ID,
				fmt.Sprintf(
					"User %s %s (username = %s, id = %d) removed reaction %s to chat %s %s (%d)",
					msg.MessageReaction.User.FirstName,
					msg.MessageReaction.User.LastName,
					msg.MessageReaction.User.Username,
					msg.MessageReaction.User.ID,
					reactionString,
					msg.MessageReaction.Chat.Type,
					msg.MessageReaction.Chat.Title,
					msg.MessageReaction.Chat.ID,
				),
			)
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
