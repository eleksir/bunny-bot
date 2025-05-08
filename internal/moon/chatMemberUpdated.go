package moon

import (
	"fmt"
	"slices"
	"time"

	"github.com/NicoNex/echotron/v3"
	"github.com/akyoto/cache"
)

// newChatMember parses one of types ChatMemberUpdated event.
// https://core.telegram.org/bots/api#chatmemberupdated
func newChatMember(msg *echotron.Update) {
	Log.Infof(
		"User %s %s (username = %s, id = %d) in chat %s %s (%d) now %s was %s",
		msg.ChatMember.NewChatMember.User.FirstName,
		msg.ChatMember.NewChatMember.User.LastName,
		msg.ChatMember.NewChatMember.User.Username,
		msg.ChatMember.NewChatMember.User.ID,
		msg.ChatMember.Chat.Type,
		msg.ChatMember.Chat.Title,
		msg.ChatMember.Chat.ID,
		msg.ChatMember.NewChatMember.Status,
		msg.ChatMember.OldChatMember.Status,
	)

	// If there is no such chatid, means that no such cache too.
	if len(ChatList) == 0 || !slices.Contains(ChatList, msg.Message.Chat.ID) {
		Log.Debugf(
			"Creating caches for new %s %s (%d)",
			msg.ChatMember.Chat.Type,
			msg.ChatMember.Chat.Title,
			msg.ChatMember.Chat.ID,
		)

		NewMembers[msg.Message.Chat.ID] = cache.New(1 * time.Minute)
		AppearedMembers[msg.Message.Chat.ID] = cache.New(1 * time.Minute)
		SquashedMembers[msg.Message.Chat.ID] = cache.New(60 * time.Minute)

		Log.Debugf(
			"Add %s %s (%d) to list of known chats",
			msg.ChatMember.Chat.Type,
			msg.ChatMember.Chat.Title,
			msg.ChatMember.Chat.ID,
		)

		ChatList = append(ChatList, msg.Message.Chat.ID)
	}

	// Check some other attributes.
	if msg.ChatMember.OldChatMember.Status == "left" && msg.ChatMember.NewChatMember.Status == "member" {
		Log.Infof(
			"Getting info about user %s %s (username = %s, id = %d) in chat %s %s (%d) from LeftMembers DB",
			msg.ChatMember.NewChatMember.User.FirstName,
			msg.ChatMember.NewChatMember.User.LastName,
			msg.ChatMember.NewChatMember.User.Username,
			msg.ChatMember.NewChatMember.User.ID,
			msg.ChatMember.Chat.Type,
			msg.ChatMember.Chat.Title,
			msg.ChatMember.Chat.ID,
		)

		leftTimestampString := Cfg.GetValue(
			"LeftMembers",
			fmt.Sprintf("%d", msg.ChatMember.Chat.ID),
			fmt.Sprintf("%d", msg.ChatMember.NewChatMember.User.ID),
		)

		// Yes, this member definitely left chat at some point, we have no question to him/her.
		if leftTimestampString != "" {
			Log.Infof(
				"User %s %s (username = %s, id = %d) in chat %s %s (%d) found in from LeftMembers DB ",
				msg.ChatMember.NewChatMember.User.FirstName,
				msg.ChatMember.NewChatMember.User.LastName,
				msg.ChatMember.NewChatMember.User.Username,
				msg.ChatMember.NewChatMember.User.ID,
				msg.ChatMember.Chat.Type,
				msg.ChatMember.Chat.Title,
				msg.ChatMember.Chat.ID,
			)

			return
		}

		Log.Infof(
			"User %s %s (username = %s, id = %d) in chat %s %s (%d) not found in LeftMembers DB user %s %s (username = %s, id = %d) in chat %s %s (%d)",
			msg.ChatMember.NewChatMember.User.FirstName,
			msg.ChatMember.NewChatMember.User.LastName,
			msg.ChatMember.NewChatMember.User.Username,
			msg.ChatMember.NewChatMember.User.ID,
			msg.ChatMember.Chat.Type,
			msg.ChatMember.Chat.Title,
			msg.ChatMember.Chat.ID,
		)

		// Spam bots come with these attributes.
		// msg.ChatMember.NewChatMember.Status = "member"
		// msg.ChatMember.OldChatMember.User != nil
		// msg.ChatMember.OldChatMember.Status = "left"

		casBanned, err := CasCheckID(msg.ChatMember.NewChatMember.User.ID)

		if err != nil {
			Log.Errorf(
				"Unable to check user %s %s (username = %s, id = %d) in chat %s %s (%d): %s",
				msg.ChatMember.NewChatMember.User.FirstName,
				msg.ChatMember.NewChatMember.User.LastName,
				msg.ChatMember.NewChatMember.User.Username,
				msg.ChatMember.NewChatMember.User.ID,
				msg.ChatMember.Chat.Type,
				msg.ChatMember.Chat.Title,
				msg.ChatMember.Chat.ID,
				err,
			)
		} else if casBanned {
			// Ban user if cas report it as banned.
			Log.Infof(
				"Ban user %s %s (username = %s, id = %d) in chat %s %s (%d) by CAS blacklist",
				msg.ChatMember.NewChatMember.User.FirstName,
				msg.ChatMember.NewChatMember.User.LastName,
				msg.ChatMember.NewChatMember.User.Username,
				msg.ChatMember.NewChatMember.User.ID,
				msg.ChatMember.Chat.Type,
				msg.ChatMember.Chat.Title,
				msg.ChatMember.Chat.ID,
			)

			squash(
				msg.ChatMember.Chat.ID,
				msg.ChatMember.Chat.Type,
				msg.ChatMember.Chat.Title,
				msg.ChatMember.NewChatMember.User.ID,
				msg.ChatMember.NewChatMember.User.FirstName,
				msg.ChatMember.NewChatMember.User.LastName,
				msg.ChatMember.NewChatMember.User.Username,
			)

			return
		}

		Log.Infof(
			"Storing info about user %s %s (username = %s, id = %d) in chat %s %s (%d) to AppearedMember DB",
			msg.ChatMember.NewChatMember.User.FirstName,
			msg.ChatMember.NewChatMember.User.LastName,
			msg.ChatMember.NewChatMember.User.Username,
			msg.ChatMember.NewChatMember.User.ID,
			msg.ChatMember.Chat.Type,
			msg.ChatMember.Chat.Title,
			msg.ChatMember.Chat.ID,
		)

		AppearedMembers[msg.ChatMember.Chat.ID].Set(msg.ChatMember.NewChatMember.User.ID, time.Now().Unix(), 1*time.Minute)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
