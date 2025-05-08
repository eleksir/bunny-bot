package moon

import (
	"slices"
	"time"

	"github.com/NicoNex/echotron/v3"
	cache "github.com/akyoto/cache"
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

		// If there is no such chatid, means that no such cache too.
		if !slices.Contains(ChatList, msg.Message.Chat.ID) {
			Log.Debugf(
				"Creating caches for new %s %s (%d)",
				msg.Message.Chat.Type,
				msg.Message.Chat.Title,
				msg.Message.Chat.ID,
			)

			NewMembers[msg.Message.Chat.ID] = cache.New(1 * time.Minute)
			AppearedMembers[msg.Message.Chat.ID] = cache.New(1 * time.Minute)
			SquashedMembers[msg.Message.Chat.ID] = cache.New(60 * time.Minute)

			Log.Debugf(
				"Add %s %s (%d) to list of known chats",
				msg.Message.Chat.Type,
				msg.Message.Chat.Title,
				msg.Message.Chat.ID,
			)

			ChatList = append(ChatList, msg.Message.Chat.ID)
		}

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
				msg.Message.Chat.Type,
				msg.Message.Chat.Title,
				user.ID,
				user.FirstName,
				user.LastName,
				user.Username,
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

		NewMembers[msg.Message.Chat.ID].Set(user.ID, time.Now().Unix(), 1*time.Minute)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
