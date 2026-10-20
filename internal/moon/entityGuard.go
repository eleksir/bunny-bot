package moon

import (
	"bunny-bot/internal/mylog"

	"github.com/NicoNex/echotron/v3"
	"github.com/davecgh/go-spew/spew"
)

// The Bot API always returns a User/Chat carrying its id, and echotron validates
// nothing: lclient.dispatch (network.go:138) is json.Unmarshal plus a check() that
// reads only the ok/error_code envelope (helpers.go:37). An entity that arrives
// without an id therefore reaches us as a silent zero, with no layer below to
// complain about it. These predicates are the whole of our defence.
//
// They are SILENT. A refusal is a fact, not an incident: reporting happens once, at
// the ingress, where the offending update can be dumped whole. A predicate that logs
// would make one broken update produce one line per entity that happens to meet it.

// usableUser reports whether a user entity is one we can actually address.
//
// The test is against zero and not against zero-or-less on purpose: a chat id is
// legitimately negative for a group or supergroup, so a positivity check would
// refuse every group chat in the project.
func usableUser(u *echotron.User) bool {
	return u != nil && u.ID != 0
}

// usableChat reports whether a chat entity is one we can actually address, with the
// same reasoning and the same zero-not-positive rule.
func usableChat(c *echotron.Chat) bool {
	return c != nil && c.ID != 0
}

// usableUpdate reports whether an update is worth dispatching at all.
//
// Its branches mirror telegramMsgParser's, in the same order and under the same
// conditions, so a branch added to one and not the other is visible by reading. That
// is not tidiness: the alternative is a validator shaped differently from the
// dispatcher, whose coverage nobody can check by eye.
//
// nil is tolerated only where absence is normal. A message with no author is a
// channel post or a service message, and refusing those would disable chat logging
// for channels. A nil element inside NewChatMembers is different — the array asserts
// a newcomer exists — and so is a nil NewChatMember.User, which asserts who changed.
// A nil OldChatMember.User is tolerated: some transitions have no previous member.
func usableUpdate(u *echotron.Update) bool {
	if u == nil {
		return false
	}

	switch {
	case u.ChatMember != nil:
		return usableChatMember(u.ChatMember)
	case u.Message != nil:
		return usableMessage(u.Message)
	case u.EditedMessage != nil:
		return usableMessage(u.EditedMessage)
	default:
		// Nothing we act on. The dispatcher's spew fallback is where this belongs.
		return true
	}
}

// usableChatMember checks a membership change.
//
// Its entities are values, not pointers: ChatMemberUpdated.Chat, .From and both
// ChatMember structs can never be nil, so there is no nil test to fall back on and
// the id check is the only thing that can catch a malformed one.
func usableChatMember(cm *echotron.ChatMemberUpdated) bool {
	return usableChat(&cm.Chat) &&
		usableUser(cm.NewChatMember.User) &&
		(cm.OldChatMember.User == nil || usableUser(cm.OldChatMember.User)) &&
		usableUser(&cm.From)
}

// usableMessage checks a message, or an edit of one — the two carry the same fields
// and are treated the same way.
//
// Message.Chat is a value: a zero chat id would otherwise reach elog, which builds
// chat_logs/<id>/yyyy/mm/dd.txt from it and would silently create a chat_logs/0
// directory and write every malformed update into it.
func usableMessage(m *echotron.Message) bool {
	if !usableChat(&m.Chat) {
		return false
	}

	if m.From != nil && !usableUser(m.From) {
		return false
	}

	if m.LeftChatMember != nil && !usableUser(m.LeftChatMember) {
		return false
	}

	for _, member := range m.NewChatMembers {
		if !usableUser(member) {
			return false
		}
	}

	// One level deep, and the nested chat only. The nested author is never rendered,
	// so a zero there is cosmetic; the level bound exists to keep the work constant.
	if m.PinnedMessage != nil && !usableChat(&m.PinnedMessage.Chat) {
		return false
	}

	if m.ReplyToMessage != nil && !usableChat(&m.ReplyToMessage.Chat) {
		return false
	}

	return true
}

// reportBrokenUpdate records an update we are refusing, in the terms it deserves:
// this is not a user being difficult and not a transient API error, it is an
// upstream that broke its own contract.
//
// The whole update is dumped, once. The dump is the point — the condition is
// unreproducible and unqueryable, since there is no id to ask Telegram about — and a
// partial dump of one field out of four would not be evidence of anything.
func reportBrokenUpdate(u *echotron.Update) {
	mylog.Error(
		"Telegram sent an update with an entity carrying no id, refusing it: nothing in it can be greeted, recorded or banned, and an API that sends one is broken",
	)

	mylog.Error(spew.Sdump(u))
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
