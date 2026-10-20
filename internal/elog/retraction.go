package elog

import (
	"github.com/NicoNex/echotron/v3"
)

// LogGreetingRetracted records that we deleted a greeting of our own — one the
// addressed member never answered, because they were banned for the silence.
//
// The log is append-only, so the question's original line stays where it was and
// this line says what became of it. A reader must be able to tell that the room
// asked a question and later removed it, rather than finding a line that
// describes a message which no longer exists.
func LogGreetingRetracted(chat *echotron.Chat, user *echotron.User, messageID int) {
	if _, on := chatLogRoot(); !on || chat == nil {
		return
	}

	body := "greeting " + itoaInt(messageID) + " retracted, member banned without answering"

	appendLine(chat.ID, renderLine(*chat, user, body))
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
