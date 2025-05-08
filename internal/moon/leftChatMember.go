package moon

import (
	"fmt"
	"time"

	"github.com/NicoNex/echotron/v3"
)

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
