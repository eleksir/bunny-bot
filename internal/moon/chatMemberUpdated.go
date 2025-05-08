package moon

import (
	"fmt"
	"time"

	"github.com/NicoNex/echotron/v3"
)

// newChatMember parses one of types ChatMemberUpdated event.
// https://core.telegram.org/bots/api#chatmemberupdated
func newChatMember(msg *echotron.Update) {
	var (
		joinTimestampString string
	)

	// Spam bots come with these attributes.
	// msg.ChatMember.NewChatMember.Status = "member"
	// msg.ChatMember.OldChatMember.User != nil
	// msg.ChatMember.OldChatMember.Status = "left"

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

	// To avoid race conditions.
	// To really avoid them, this particular event and newChatMember handler must be serialized via channel.
	time.Sleep(10 * time.Second)

	Log.Info("Processing UpdateUser event")

	// Try to get timestamp when this member joined to chat
	joinTimestampString = Cfg.GetValue(
		"NewMembers",
		fmt.Sprintf("%d", msg.ChatMember.Chat.ID),
		fmt.Sprintf("%d", msg.ChatMember.NewChatMember.User.ID),
	)

	// Нормальный человек, который зашёл с правильным ивентом, через message с атрибутом new_chat_members
	// https://core.telegram.org/bots/api#message
	if joinTimestampString != "" {
		return
	}

	// Check some other attributes.
	if msg.ChatMember.OldChatMember.Status == "left" && msg.ChatMember.NewChatMember.Status == "member" {
		leftTimestampString := Cfg.GetValue(
			"LeftMembers",
			fmt.Sprintf("%d", msg.ChatMember.Chat.ID),
			fmt.Sprintf("%d", msg.ChatMember.NewChatMember.User.ID),
		)

		// Yes, this member definitely left chat at some point, we have no question to him/her.
		if leftTimestampString != "" {
			return
		}

		// Very suspicious coincedence, looks like we've got bot here.
		banDate := time.Now().Unix() + 45

		if res, err := Tg.BanChatMember(
			msg.ChatID(),
			msg.ChatMember.NewChatMember.User.ID,
			&echotron.BanOptions{
				UntilDate: int(banDate),
			},
		); err != nil {
			Log.Errorf(
				"Unable to ban user with id %d in chat %s: %s",
				msg.ChatMember.NewChatMember.User.ID,
				msg.ChatID(),
				res.Description,
			)
		} else {
			if !res.Ok {
				Log.Errorf(
					"Unable to ban user with id %d in chat %s: %s",
					msg.ChatMember.NewChatMember.User.ID,
					msg.ChatID(),
					res.Description,
				)
			} else {
				Log.Infof(
					"Banned user with id %d in chat %s.",
					msg.ChatMember.NewChatMember.User.ID,
					msg.ChatID(),
				)
			}
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
