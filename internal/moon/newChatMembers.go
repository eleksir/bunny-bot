package moon

import (
	"fmt"
	"slices"
	"time"

	"github.com/NicoNex/echotron/v3"
)

// newChatMembers parses NewChatMember event.
func newChatMembers(msg *echotron.Update) {
	// If there is no such chatid, means that no such cache too.
	Log.Debugf(
		"Add %s %s (%d) to list of known chats",
		msg.Message.Chat.Type,
		msg.Message.Chat.Title,
		msg.Message.Chat.ID,
	)

	ChatList = append(ChatList, msg.Message.Chat.ID)
	slices.Sort(ChatList)
	ChatList = slices.Compact(ChatList)

	AddChat(&msg.Message.Chat)

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

		AddUser(user)

		casBanned, err := CasCheckID(user.ID)

		if err != nil {
			Log.Error(err)
		}

		// Ban user if cas report it as banned.
		if casBanned {
			Log.Infof(
				"Ban user %s %s (username = %s, id = %d) in %s %s (%d) by CAS blacklist",
				user.FirstName,
				user.LastName,
				user.Username,
				user.ID,
				msg.Message.Chat.Type,
				msg.Message.Chat.Title,
				msg.Message.Chat.ID,
			)

			squash(
				msg.Message.Chat.ID,
				user.ID,
			)
		}

		// Do not return on prev. step.

		Log.Infof(
			"Skip banning user %s %s (username = %s, id = %d) in %s %s (%d) because it was not found in CAS blacklist",
			user.FirstName,
			user.LastName,
			user.Username,
			user.ID,
			msg.Message.Chat.Type,
			msg.Message.Chat.Title,
			msg.Message.Chat.ID,
		)

		Log.Infof(
			"Add user %s %s (username = %s, id = %d) in %s %s (%d) to NewMember DB",
			user.FirstName,
			user.LastName,
			user.Username,
			user.ID,
			msg.Message.Chat.Type,
			msg.Message.Chat.Title,
			msg.Message.Chat.ID,
		)

		key := fmt.Sprintf("%d+%d", msg.Message.Chat.ID, user.ID)
		NewMembers.Set(key, time.Now().Unix(), 1*time.Minute)
		PendingCASMembers.Set(key, time.Now().Unix(), 30*time.Minute)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
