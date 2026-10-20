package moon

import (
	"bunny-bot/internal/mylog"
	"time"

	"github.com/NicoNex/echotron/v3"
)

// SweepGreetings runs until a termination signal stops it. Every
// greetSweepInterval it looks at join events that are at least greetReplyWindow
// old and decides whether each newcomer ever answered the greeting:
//
//   - answered (a message from that user exists after the join, within the
//     reply window) -> newcomer is legitimate, drop the join event(s);
//   - not answered -> newcomer is a spammer, ban them and drop the join
//     event(s).
//
// A user may appear in PendingGreetings more than once (rejoins). For every user we
// re-query both caches so that a single answer rescues all of their join events
// and a single unanswered join bans them once.
//
// The stop is checked between passes and never during one, so a sweep in
// progress always finishes and can never be running while the caches are being
// closed underneath it.
func SweepGreetings() {
	if !startWorker() {
		return
	}

	defer stopWorker()

	timer := time.NewTimer(greetSweepInterval)
	defer timer.Stop()

	for {
		select {
		case <-stopRequested:
			mylog.Info("Greeting sweeper stopped")

			return

		case <-timer.C:
			sweepGreetings()

			timer.Reset(greetSweepInterval)
		}
	}
}

// sweepGreetings performs a single sweep pass against the current wall clock.
func sweepGreetings() {
	sweepGreetingsAt(time.Now())
}

// sweepGreetingsAt performs a single sweep pass as of now. Taking the clock as
// a parameter is what makes the reply-window boundary testable without sleeping
// for 15 minutes.
func sweepGreetingsAt(now time.Time) {
	joins := PendingGreetings.Snapshot()
	messages := GreetAnswers.Snapshot()

	windowMicros := greetReplyWindow.Microseconds()
	nowMicros := now.UnixMicro()

	// Group stale join events per (chat, user). A user is a distinct identity in
	// every chat: being silent in one chat says nothing about another, and
	// answering in one chat must not rescue them in another.
	staleByMember := make(map[chatUser][]joinRef)

	for joinMicros, recorded := range joins {
		join := recorded.join

		// Not stale yet.
		if nowMicros-joinMicros <= windowMicros {
			continue
		}

		// A join event carries the newcomer's user in NewChatMembers.
		if join == nil || len(join.NewChatMembers) == 0 {
			continue
		}

		user := join.NewChatMembers[0]

		// An entity with no id cannot be banned, and the dump has already been
		// taken. Drop the join rather than keeping it: re-reporting the same
		// broken update every minute for the rest of its TTL is noise that
		// teaches an operator to ignore this log.
		if !usableUser(user) {
			PendingGreetings.Delete(joinMicros)

			continue
		}

		member := chatUser{chatID: join.Chat.ID, userID: user.ID}
		staleByMember[member] = append(staleByMember[member], joinRef{
			micros:     joinMicros,
			chat:       join.Chat,
			user:       user,
			greetingID: recorded.greetingID,
		})
	}

	for member, refs := range staleByMember {
		// A single answer rescues every join this user has *in this chat*, no
		// matter which join the reply was actually answering.
		answered := false

		for _, ref := range refs {
			if userAnswered(member.chatID, member.userID, ref.micros, windowMicros, messages) {
				answered = true

				break
			}
		}

		if answered {
			mylog.Debugf(
				"User %d answered the greeting in chat %d, not banning",
				member.userID,
				member.chatID,
			)
		} else {
			mylog.Infof(
				"User %d did not answer the greeting in chat %d within %s, banning as spammer",
				member.userID,
				member.chatID,
				greetReplyWindow,
			)

			// squash already logs every outcome at the level it deserves: a
			// ban that landed is a debug line and a duplicate we deliberately
			// suppressed is an info line. Only a refused ban needs anything
			// added here, because it is the one case where the user is still
			// in the chat and will never be re-judged.
			outcome := squash(&refs[0].chat, refs[0].user)
			if outcome == banFailed {
				mylog.Errorf(
					"BAN FAILED for user %d in chat %d — they stayed in the chat and were not re-judged",
					member.userID,
					member.chatID,
				)
			} else {
				// The user is out of the chat, so every question we asked them
				// is a leftover. That includes the debounced repeat: they are
				// banned either way, and the room must not be left asking a
				// question nobody can answer.
				retractAskedGreetings(refs)
			}
		}

		// Drop every join event of this user *in this chat* so the next pass
		// does not re-decide them. Joins in their other chats are untouched.
		forgetPendingGreetings(member.chatID, member.userID)
	}
}

// userAnswered reports whether a message authored by userID in chatID exists
// in messages with a timestamp in (joinMicros, joinMicros+windowMicros].
//
// Both the chat and the author must match: the newcomer must answer in the same
// chat they were greeted in. Messages sent before the join do not count either,
// because the greeting is only sent in response to the join.
func userAnswered(
	chatID int64,
	userID int64,
	joinMicros int64,
	windowMicros int64,
	messages map[int64]*echotron.Message,
) bool {
	for msgMicros, msg := range messages {
		if msg == nil || msg.From == nil {
			continue
		}

		if msg.From.ID != userID || msg.Chat.ID != chatID {
			continue
		}

		if msgMicros > joinMicros && msgMicros <= joinMicros+windowMicros {
			return true
		}
	}

	return false
}

// forgetPendingGreetings removes every PendingGreetings entry whose newcomer is userID in the
// given chat. Joins by the same user in other chats are deliberately preserved:
// those are independent memberships with their own greetings and their own
// verdicts.
func forgetPendingGreetings(chatID int64, userID int64) {
	for joinMicros, recorded := range PendingGreetings.Snapshot() {
		join := recorded.join
		if join == nil || join.Chat.ID != chatID || len(join.NewChatMembers) == 0 {
			continue
		}

		if join.NewChatMembers[0] != nil && join.NewChatMembers[0].ID == userID {
			PendingGreetings.Delete(joinMicros)
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
