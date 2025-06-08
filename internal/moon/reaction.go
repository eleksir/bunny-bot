package moon

import "github.com/NicoNex/echotron/v3"

func reaction(msg *echotron.Update) {
	if Cfg.LogChats {
		var reactionString string

		for _, item := range msg.MessageReaction.NewReaction {
			if item.Type == "Emoji" {
				reactionString += item.Emoji
			} else if item.Type == "CustomEmoji" {
				reactionString += item.CustomEmoji
			}
		}

		if reactionString != "" {
			Log.Infof(
				"User %s %s (username = %s, id = %d) set reaction %s to chat %s %s (%d)",
				msg.MessageReaction.User.FirstName,
				msg.MessageReaction.User.LastName,
				msg.MessageReaction.User.Username,
				msg.MessageReaction.User.ID,
				reactionString,
				msg.MessageReaction.Chat.Type,
				msg.MessageReaction.Chat.Title,
				msg.MessageReaction.Chat.ID,
			)
		} else {
			for _, item := range msg.MessageReaction.OldReaction {
				if item.Type == "Emoji" {
					reactionString += item.Emoji
				} else if item.Type == "CustomEmoji" {
					reactionString += item.CustomEmoji
				}
			}

			Log.Infof(
				"User %s %s (username = %s, id = %d) removed reaction %s to chat %s %s (%d)",
				msg.MessageReaction.User.FirstName,
				msg.MessageReaction.User.LastName,
				msg.MessageReaction.User.Username,
				msg.MessageReaction.User.ID,
				reactionString,
				msg.MessageReaction.Chat.Type,
				msg.MessageReaction.Chat.Title,
				msg.MessageReaction.Chat.ID,
			)
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
