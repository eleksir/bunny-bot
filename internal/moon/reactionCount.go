package moon

import (
	"fmt"

	"github.com/NicoNex/echotron/v3"
	"github.com/davecgh/go-spew/spew"
)

func reactionCount(msg *echotron.Update) {
	ChatLog(
		msg.MessageReaction.Chat.ID,
		fmt.Sprintf(
			"Reaction count event in %s %s (%d): %s",
			msg.MessageReaction.Chat.Type,
			msg.MessageReaction.Chat.Title,
			msg.MessageReaction.Chat.ID,
			spew.Sdump(msg),
		),
	)
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
