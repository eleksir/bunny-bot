package moon

import (
	"fmt"
	"time"

	"github.com/NicoNex/echotron/v3"
)

// leftChatMember handles left_chat_member event.
func leftChatMember(msg *echotron.Update) {
	Log.Infof(
		"User %s %s (username = %s, id = %d) left %s %s (%d)",
		msg.Message.LeftChatMember.FirstName,
		msg.Message.LeftChatMember.LastName,
		msg.Message.LeftChatMember.Username,
		msg.Message.LeftChatMember.ID,
		msg.Message.Chat.Type,
		msg.Message.Chat.Title,
		msg.Message.Chat.ID,
	)

	// Save people that left chat to pebbledb. To preserve their state.
	// Left folks can eventually return and in chat_mamber_updated event their state will be left->member.
	if err := Cfg.SaveKeyValue(
		"LeftMembers",
		fmt.Sprintf("%d", msg.Message.Chat.ID),
		fmt.Sprintf("%d", msg.Message.LeftChatMember.ID),
		fmt.Sprintf("%d", time.Now().Unix()),
	); err != nil {
		Log.Error(err)
	} else {
		Log.Infof("Info about left userid %d saved to left users db", msg.Message.LeftChatMember.ID)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
