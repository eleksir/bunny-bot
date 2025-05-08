package moon

import (
	"fmt"
	"time"

	"github.com/NicoNex/echotron/v3"
)

func textMessage(msg *echotron.Update) {
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

	// Update database of new members :)
	// Try to get timestamp when this member joined to chat
	if joinTimestampString := Cfg.GetValue(
		"NewMembers",
		fmt.Sprintf("%d", msg.ChatMember.Chat.ID),
		fmt.Sprintf("%d", msg.ChatMember.NewChatMember.User.ID),
	); joinTimestampString == "" {
		if err := Cfg.SaveKeyValue(
			"NewMembers",
			fmt.Sprintf("%d", msg.Message.Chat.ID),
			fmt.Sprintf("%d", msg.Message.From.ID),
			fmt.Sprintf("%d", time.Now().Unix()),
		); err != nil {
			Log.Error(err)
		} else {
			Log.Infof("Info about new userid %d saved to new users db", msg.Message.From.ID)
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
