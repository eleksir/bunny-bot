package moon

import (
	"fmt"
	"time"

	"github.com/NicoNex/echotron/v3"
)

// newChatMembers parses NewChatMember event.
func newChatMembers(msg *echotron.Update) {
	for _, user := range msg.Message.NewChatMembers {
		Log.Infof(
			"User %s %s (username = %s, id = %d) joined %s %s (%d)",
			user.FirstName,
			user.LastName,
			user.Username,
			user.ID,
			msg.Message.Chat.Type,
			msg.Message.Chat.Title,
			msg.Message.Chat.ID,
		)

		if err := Cfg.SaveKeyValue(
			"NewMembers",
			fmt.Sprintf("%d", msg.Message.Chat.ID),
			fmt.Sprintf("%d", user.ID),
			fmt.Sprintf("%d", time.Now().Unix()),
		); err != nil {
			Log.Error(err)
		} else {
			Log.Infof("Info about new userid %d saved to new users db", user.ID)
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
