package moon

import (
	"fmt"
	"slices"
	"time"

	"github.com/NicoNex/echotron/v3"
)

// newChatMember parses one of types ChatMemberUpdated event.
// https://core.telegram.org/bots/api#chatmemberupdated
func chatMemberUpdated(msg *echotron.Update) {
	AddChat(&msg.ChatMember.Chat)
	AddUser(msg.ChatMember.NewChatMember.User)

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

	Log.Debugf(
		"Add %s %s (%d) to list of known chats",
		msg.ChatMember.Chat.Type,
		msg.ChatMember.Chat.Title,
		msg.ChatMember.Chat.ID,
	)

	ChatList = append(ChatList, msg.ChatMember.Chat.ID)
	slices.Sort(ChatList)
	ChatList = slices.Compact(ChatList)

	// Check some other attributes.
	switch {
	case msg.ChatMember.OldChatMember.Status == "left" && msg.ChatMember.NewChatMember.Status == "member":
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
				"User %s %s (username = %s, id = %d) in chat %s %s (%d) found in from LeftMembers DB",
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
			"User %s %s (username = %s, id = %d) in chat %s %s (%d) not found in LeftMembers DB",
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
				msg.ChatMember.NewChatMember.User.ID,
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

		key := fmt.Sprintf("%d+%d", msg.ChatMember.Chat.ID, msg.ChatMember.NewChatMember.User.ID)
		AppearedMembers.Set(key, time.Now().Unix(), 1*time.Minute)
		PendingCASMembers.Set(key, time.Now().Unix(), 30*time.Minute)

	case msg.ChatMember.OldChatMember.Status == "member" && msg.ChatMember.NewChatMember.Status == "kicked":
		Log.Infof(
			"User %s %s (username = %s, id = %d) in chat %s %s (%d) was banned",
			msg.ChatMember.NewChatMember.User.FirstName,
			msg.ChatMember.NewChatMember.User.LastName,
			msg.ChatMember.NewChatMember.User.Username,
			msg.ChatMember.NewChatMember.User.ID,
			msg.ChatMember.Chat.Type,
			msg.ChatMember.Chat.Title,
			msg.ChatMember.Chat.ID,
		)

		key := fmt.Sprintf("%d+%d", msg.ChatMember.Chat.ID, msg.ChatMember.NewChatMember.User.ID)

		Log.Debugf(
			"Put user %s %s (username = %s, id = %d) in %s %s (%d) into SquashedMembers to debounce ban.",
			msg.ChatMember.NewChatMember.User.FirstName,
			msg.ChatMember.NewChatMember.User.LastName,
			msg.ChatMember.NewChatMember.User.Username,
			msg.ChatMember.NewChatMember.User.ID,
			msg.ChatMember.Chat.Type,
			msg.ChatMember.Chat.Title,
			msg.ChatMember.Chat.ID,
		)

		SquashedMembers.Set(key, true, 30*time.Second)

		Log.Debugf(
			"Storing info about user %s %s (username = %s, id = %d) in chat %s %s (%d) to BannedMembers DB",
			msg.ChatMember.NewChatMember.User.FirstName,
			msg.ChatMember.NewChatMember.User.LastName,
			msg.ChatMember.NewChatMember.User.Username,
			msg.ChatMember.NewChatMember.User.ID,
			msg.ChatMember.Chat.Type,
			msg.ChatMember.Chat.Title,
			msg.ChatMember.Chat.ID,
		)

		if err := Cfg.SaveKeyValue(
			"BannedMembers",
			fmt.Sprintf("%d", msg.ChatMember.Chat.ID),
			fmt.Sprintf("%d", msg.ChatMember.NewChatMember.User.ID),
			fmt.Sprintf("%d", time.Now().Unix()),
		); err != nil {
			Log.Error(err)
		} else {
			Log.Infof("Info about banned userid %d saved to BannedMembers DB", msg.ChatMember.Chat.ID)
		}

	default:
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
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
