package moon

import (
	"time"

	"bunny-bot/internal/elog"
	"bunny-bot/internal/mylog"

	"github.com/NicoNex/echotron/v3"
)

// newChatMembers parses NewChatMember event.
func newChatMembers(msg *echotron.Update) {
	// Telegram normally reports one newcomer per service message. If that ever
	// stops being true, the sweeper would only judge the first member, so make it
	// visible in the logs rather than silently ignoring the rest.
	if len(msg.Message.NewChatMembers) > 1 {
		mylog.Warnf(
			"MULTI-MEMBER JOIN in chat %d: %d newcomers in one event, only the first is moderated",
			msg.Message.Chat.ID,
			len(msg.Message.NewChatMembers),
		)
	}

	// Keep only the newcomers we can actually address. In production the ingress has
	// already refused such an update; this filter is the nil-safety that keeps a
	// direct caller from recording somebody we cannot name.
	newcomers := make([]*echotron.User, 0, len(msg.Message.NewChatMembers))

	for _, user := range msg.Message.NewChatMembers {
		if !usableUser(user) {
			continue
		}

		newcomers = append(newcomers, user)
	}

	if len(newcomers) == 0 {
		return
	}

	// Remember the join so the greeting sweeper can tell later whether this
	// newcomer answered, and so a ban can retract the greeting we sent.
	//
	// The join is recorded before the greeting exists, so the greeting id starts
	// at zero and is filled in once the send returns. A join whose greeting was
	// never delivered therefore carries zero, which is exactly the marker
	// retractGreetings skips.
	joinMicros := time.Now().UnixMicro()
	PendingGreetings.Set(joinMicros, pendingGreeting{join: msg.Message}, greetCacheTTL)

	// Record the join in the chat log before anything else, so the arrival is on
	// record even if the greeting later fails. elog writes one line per newcomer
	// itself, so this is called once for the whole event rather than per user.
	elog.LogMessage(msg.Message)

	for _, user := range newcomers {
		mylog.Infof(
			"User %s %s (username = %s, id = %d) joined %s %s (%d)",
			user.FirstName,
			user.LastName,
			user.Username,
			user.ID,
			msg.Message.Chat.Type,
			msg.Message.Chat.Title,
			msg.Message.Chat.ID,
		)

		// Greet the newcomer, so that a spammer who never answers becomes
		// detectable. If the greeting could not be delivered the newcomer never
		// saw the question, so they must not be judged on their silence.
		greetingID, delivered := sendGreeting(msg.Message.Chat.ID, user)
		if !delivered {
			PendingGreetings.Delete(joinMicros)

			return
		}

		// Re-record the join carrying the greeting's id. A member who joined
		// twice has one recorded join per greeting, so each id is retracted
		// separately.
		PendingGreetings.Set(joinMicros, pendingGreeting{
			join:       msg.Message,
			greetingID: greetingID,
		}, greetCacheTTL)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
