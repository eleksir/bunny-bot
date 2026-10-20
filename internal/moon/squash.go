package moon

import (
	"bunny-bot/internal/mylog"
	"fmt"
	"time"

	"github.com/NicoNex/echotron/v3"
)

// squash bans the given chat member and reports what it actually did.
//
// The three outcomes are deliberately distinct. A debounced call is not a
// failure: the user is already banned, which is the outcome we wanted, and the
// caller must not report it as an error.
//
// The user and chat come from the join event that triggered the judgement
// rather than from a cache. That event already carries both in full, so a cache
// would only ever be a second, staler copy of data we are holding — and a miss
// in it would produce a ban logged with blank names.
func squash(chat *echotron.Chat, user *echotron.User) banOutcome {
	// A chat or user with no id is an entity we cannot address, not a ban that
	// failed. In production the ingress has already refused such an update, so
	// reaching here means a caller handed us one directly; the guard is the
	// nil-safety that keeps that from becoming a call against chat 0.
	if !usableChat(chat) {
		return banFailed
	}

	if !usableUser(user) {
		return banFailed
	}

	// Telegram treats any ban shorter than 30 seconds (and longer than 1 year)
	// as permanent, so this deliberately short offset produces a permanent ban.
	// Do not "fix" it into a longer duration.
	banUntilDate := time.Now().Unix() + 10

	chatid, userid := chat.ID, user.ID

	key := fmt.Sprintf("%d+%d", chatid, userid)

	if _, exist := SquashedMembers.Get(key); exist {
		mylog.Infof(
			"Skip banning user %s %s (username = %s, id = %d) in %s %s (%d), already banned",
			user.FirstName,
			user.LastName,
			user.Username,
			user.ID,
			chat.Type,
			chat.Title,
			chat.ID,
		)

		return banSuppressed
	}

	res, err := Tg.BanChatMember(
		chatid,
		userid,
		&echotron.BanOptions{
			UntilDate: int(banUntilDate),
		},
	)
	if err != nil {
		mylog.Errorf(
			"Unable to ban user user %s %s (username = %s, id = %d) in %s %s (%d): %s",
			user.FirstName,
			user.LastName,
			user.Username,
			user.ID,
			chat.Type,
			chat.Title,
			chat.ID,
			err,
		)

		return banFailed
	}

	if !res.Ok {
		mylog.Errorf(
			"Unable to ban user %s %s (username = %s, id = %d) in %s %s (%d): %s",
			user.FirstName,
			user.LastName,
			user.Username,
			user.ID,
			chat.Type,
			chat.Title,
			chat.ID,
			res.Description,
		)

		return banFailed
	}

	mylog.Debugf(
		"Banned user %s %s (username = %s, id = %d) in %s %s (%d).",
		user.FirstName,
		user.LastName,
		user.Username,
		user.ID,
		chat.Type,
		chat.Title,
		chat.ID,
	)

	mylog.Debugf(
		"Put user %s %s (username = %s, id = %d) in %s %s (%d) into SquashedMembers to debounce ban.",
		user.FirstName,
		user.LastName,
		user.Username,
		user.ID,
		chat.Type,
		chat.Title,
		chat.ID,
	)

	SquashedMembers.Set(key, true, 30*time.Second)

	return banApplied
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
