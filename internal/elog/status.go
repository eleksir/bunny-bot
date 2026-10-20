package elog

import (
	"strings"

	"github.com/NicoNex/echotron/v3"
)

// LogUserStatus records a change to a chat member's status — join, leave, ban,
// promote, demote, restrict — as carried by Update.ChatMember.
func LogUserStatus(ev *echotron.ChatMemberUpdated) {
	if _, on := chatLogRoot(); !on || ev == nil {
		return
	}

	who := ev.NewChatMember.User
	if who == nil {
		who = ev.OldChatMember.User
	}

	var b strings.Builder

	// The subject is already rendered by renderLine's "user <id> (<desc>)"
	// prefix, so the body starts at the transition itself.
	b.WriteString("status ")
	b.WriteString(orUnknown(ev.OldChatMember.Status))
	b.WriteString(" -> ")
	b.WriteString(orUnknown(ev.NewChatMember.Status))

	if ev.NewChatMember.CustomTitle != "" {
		b.WriteString(" custom_title=")
		b.WriteString(ev.NewChatMember.CustomTitle)
	}

	b.WriteString(" by ")
	b.WriteString(describeUserRef(&ev.From))

	appendLine(ev.Chat.ID, renderLine(ev.Chat, who, b.String()))
}

// describeUserRef renders a value-type User reference.
func describeUserRef(u *echotron.User) string {
	if u == nil {
		return "user unknown"
	}

	if desc := describeUser(u); desc != "" {
		return "user " + itoa(u.ID) + " (" + desc + ")"
	}

	return "user " + itoa(u.ID)
}

// orUnknown renders a status, or a placeholder when Telegram sent none.
func orUnknown(status string) string {
	if status == "" {
		return "unknown"
	}

	return status
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
