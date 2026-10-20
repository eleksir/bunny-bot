package elog

import (
	"github.com/NicoNex/echotron/v3"
)

// LogMessage records a chat message.
//
// It handles both ordinary messages — text, media, captions, replies — and the
// service messages Telegram delivers as a Message with a field set instead of
// text: joins, leaves, title changes, pins, migrations, forum topics. A message
// with nothing describable is skipped rather than logged as noise.
func LogMessage(msg *echotron.Message) {
	logMessage(msg, "")
}

// LogEditedMessage records an edited message. The line is marked "edited" so a
// reader can tell it apart from the original.
func LogEditedMessage(msg *echotron.Message) {
	logMessage(msg, "edited ")
}

// logMessage is the shared body of LogMessage and LogEditedMessage.
func logMessage(msg *echotron.Message, prefix string) {
	if _, on := chatLogRoot(); !on || msg == nil {
		return
	}

	for _, body := range describeMessage(msg) {
		appendLine(msg.Chat.ID, prefix+renderLine(msg.Chat, msg.From, body))
	}
}

// describeMessage returns one rendered body per line to write. It returns more
// than one only for a multi-member join, where one line per newcomer keeps each
// fact greppable instead of cramming every newcomer onto one line.
func describeMessage(msg *echotron.Message) []string {
	if bodies := describeService(msg); len(bodies) > 0 {
		return bodies
	}

	body := describePayload(msg)
	if body == "" {
		return nil
	}

	if msg.ReplyToMessage != nil {
		body += " reply_to=" + itoaInt(msg.ReplyToMessage.ID)
	}

	return []string{body}
}

// describeService renders the service-message fields, newest-first, and returns
// nil when the message is not a service message.
func describeService(msg *echotron.Message) []string {
	var bodies []string

	// A join with N newcomers writes N lines.
	for _, u := range msg.NewChatMembers {
		if u == nil {
			continue
		}

		bodies = append(bodies, describeReference(u)+" joined the chat")
	}

	if msg.LeftChatMember != nil {
		bodies = append(bodies, describeReference(msg.LeftChatMember)+" left the chat")
	}

	if msg.NewChatTitle != "" {
		bodies = append(bodies, "chat title changed to "+msg.NewChatTitle)
	}

	if len(msg.NewChatPhoto) > 0 {
		bodies = append(bodies, "chat photo changed")
	}

	if msg.PinnedMessage != nil {
		bodies = append(bodies, "pinned message "+itoaInt(msg.PinnedMessage.ID))
	}

	if msg.GroupChatCreated {
		bodies = append(bodies, "group chat created")
	}

	if msg.MigrateToChatID != 0 {
		bodies = append(bodies, "chat migrated to "+itoa(int64(msg.MigrateToChatID)))
	}

	if msg.MigrateFromChatID != 0 {
		bodies = append(bodies, "chat migrated from "+itoa(int64(msg.MigrateFromChatID)))
	}

	if msg.ConnectedWebsite != "" {
		bodies = append(bodies, "connected website "+msg.ConnectedWebsite)
	}

	if msg.ForumTopicCreated != nil {
		bodies = append(bodies, "forum topic "+msg.ForumTopicCreated.Name+" created")
	}

	if msg.ForumTopicClosed != nil {
		bodies = append(bodies, "forum topic closed")
	}

	if msg.ForumTopicReopened != nil {
		bodies = append(bodies, "forum topic reopened")
	}

	if msg.GeneralForumTopicHidden != nil {
		bodies = append(bodies, "general forum topic hidden")
	}

	if msg.GeneralForumTopicUnhidden != nil {
		bodies = append(bodies, "general forum topic unhidden")
	}

	if msg.VideoChatScheduled != nil {
		bodies = append(bodies, "video chat scheduled")
	}

	return bodies
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
