package moon

import (
	"bunny-bot/internal/elog"
	"bunny-bot/internal/mylog"

	"github.com/NicoNex/echotron/v3"
)

// retractGreetings deletes the greetings we sent to a member who has just been
// banned for never answering them, so the room is not left asking a question of
// someone who no longer has access.
//
// Ids <= 0 mean "no greeting was ever sent for that join" and are skipped
// silently. A deletion Telegram refuses is logged and the pass continues: the
// greeting is not worth abandoning the rest over, and a failed retraction never
// changes the ban that has already landed.
func retractGreetings(chat *echotron.Chat, user *echotron.User, greetingIDs []int) retractOutcome {
	retracted, failed := 0, 0

	for _, id := range greetingIDs {
		if id <= 0 {
			continue
		}

		if _, err := Tg.DeleteMessage(chat.ID, id); err != nil {
			failed++

			reportRetractFailure(chat, user, id, err)

			continue
		}

		retracted++

		elog.LogGreetingRetracted(chat, user, id)
	}

	return retractOutcomeFor(retracted, failed)
}

// reportRetractFailure logs a refused deletion. It is a WARN rather than an ERROR
// because the ban it followed already succeeded and nothing depends on the
// greeting being gone.
func reportRetractFailure(chat *echotron.Chat, user *echotron.User, id int, err error) {
	mylog.Warnf(
		"Left an unanswered greeting in place: cannot delete message %d in chat %d for user %d: %s",
		id,
		chat.ID,
		userIDOf(user),
		err,
	)
}

// retractAskedGreetings takes down the questions this member was asked and never
// answered. One id per recorded join, so a member who joined twice had two
// greetings and both go.
func retractAskedGreetings(refs []joinRef) {
	ids := make([]int, 0, len(refs))
	for _, ref := range refs {
		ids = append(ids, ref.greetingID)
	}

	chat, user := refs[0].chat, refs[0].user

	switch retractGreetings(&chat, user, ids) {
	case greetingRetracted:
		mylog.Debugf(
			"Deleted %d unanswered greeting(s) asked of user %d in chat %d",
			countNamedGreetings(ids),
			user.ID,
			chat.ID,
		)
	case greetingRetractFailed:
		mylog.Warnf(
			"Banned user %d in chat %d, but could not delete every greeting we asked them",
			user.ID,
			chat.ID,
		)
	case nothingToRetract:
	}
}

// countNamedGreetings counts the ids that actually named a message; the rest
// mean no greeting was ever delivered.
func countNamedGreetings(ids []int) int {
	n := 0

	for _, id := range ids {
		if id > 0 {
			n++
		}
	}

	return n
}

// retractOutcomeFor collapses the two counts into the single value callers act on.
func retractOutcomeFor(retracted int, failed int) retractOutcome {
	switch {
	case retracted > 0 && failed == 0:
		return greetingRetracted
	case failed > 0:
		return greetingRetractFailed
	default:
		return nothingToRetract
	}
}

// userIDOf names a user in a log line without dereferencing a possibly nil one.
func userIDOf(user *echotron.User) int64 {
	if user == nil {
		return 0
	}

	return user.ID
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
