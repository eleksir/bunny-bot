package moon

import (
	"bunny-bot/internal/elog"

	"github.com/NicoNex/echotron/v3"
)

// logOnlyMessage writes a message to the chat log and nothing else.
//
// It handles the messages the greeting detector must NOT count as answers: a photo with no
// caption, a sticker, a location, a leave, a title change. They are real events in the room
// and belong in the log — which is also what makes a later ban diagnosable, since a moderator
// can read what the member actually replied with — but a reply with no text is not an
// introduction, and recording it in GreetAnswers would let any attachment rescue a spammer
// from the ban.
//
// The dispatcher has already established that this is neither a join nor a text message; it
// is not re-checked here.
func logOnlyMessage(msg *echotron.Update) {
	elog.LogMessage(msg.Message)
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
