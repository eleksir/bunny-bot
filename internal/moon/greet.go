package moon

import (
	"fmt"
	"strings"
	"time"

	"github.com/NicoNex/echotron/v3"
)

// Greet makes greet for user and saves greet event to greet db.
func Greet(users *[]echotron.User, chat *echotron.Chat) {
	var (
		message   = GreetMessages[Random.Intn(len(GreetMessages))]
		usernames string
	)

	for _, u := range *users {
		switch {
		case strings.TrimSpace(u.FirstName) != "" && strings.TrimSpace(u.LastName) != "":
			usernames += fmt.Sprintf(
				"[%s %s](tg://user?id=%d), ",
				u.FirstName,
				u.LastName,
				u.ID,
			)

		case strings.TrimSpace(u.FirstName) != "":
			usernames += fmt.Sprintf(
				"[%s](tg://user?id=%d), ",
				u.FirstName,
				u.ID,
			)

		case strings.TrimSpace(u.LastName) != "":
			usernames += fmt.Sprintf(
				"[%s](tg://user?id=%d), ",
				u.LastName,
				u.ID,
			)

		case strings.TrimSpace(u.Username) != "":
			usernames += fmt.Sprintf(
				"[%s](tg://user?id=%d), ",
				u.Username,
				u.ID,
			)

		default:
			usernames += fmt.Sprintf(
				"[%d](tg://user?id=%d), ",
				u.ID,
				u.ID,
			)
		}
	}

	usernames = strings.TrimSpace(usernames)
	usernames = strings.TrimRight(usernames, ",")

	message += fmt.Sprintf(
		" %s. Представьтес, пожалуйста, и расскажите, что вас сюда привело.",
		usernames,
	)

	Tg.SendMessage(message, chat.ID, &echotron.MessageOptions{ParseMode: "MarkdownV2"})

	timestamp := time.Now().Unix()

	for _, u := range *users {
		Log.Infof(
			"Saving %s %s (username = %s, id = %d) from chat %s %s (%d) to GreetMembers DB",
			u.FirstName,
			u.LastName,
			u.Username,
			u.ID,
			chat.Type,
			chat.Title,
			chat.ID,
		)

		GreetChatMembers.Set(
			fmt.Sprintf("%d+%d", chat.ID, u.ID),
			timestamp,
			60*time.Minute,
		)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
