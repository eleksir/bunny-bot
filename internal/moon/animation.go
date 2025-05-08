package moon

import "github.com/NicoNex/echotron/v3"

func animation(msg *echotron.Update) {
	Log.Infof(
		"User %s %s (username = %s, id = %d) sent animation to chat %s %s (%d)",
		msg.Message.From.FirstName,
		msg.Message.From.LastName,
		msg.Message.From.Username,
		msg.Message.From.ID,
		msg.Message.Chat.Type,
		msg.Message.Chat.Title,
		msg.Message.Chat.ID,
	)
}
