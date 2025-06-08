package moon

import (
	"github.com/NicoNex/echotron/v3"
	"github.com/davecgh/go-spew/spew"
)

func reactionCount(msg *echotron.Update) {
	Log.Infof(
		"Reaction count event in %s %s (%d): %",
		msg.MessageReaction.Chat.Type,
		msg.MessageReaction.Chat.Title,
		msg.MessageReaction.Chat.ID,
		spew.Sdump(msg),
	)
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
