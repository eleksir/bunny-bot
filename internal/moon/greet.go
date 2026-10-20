package moon

import (
	"bunny-bot/internal/mylog"
	"math/rand/v2"
	"strings"

	"github.com/NicoNex/echotron/v3"
)

// greetings are the templates used to greet a newly joined member. The single
// verb is replaced with the newcomer's computed display name. The fourth one is
// intentionally kept verbatim despite looking like a typo.
var greetings = []string{
	"Дарована, %s. Представьтес, пожалуйста, и расскажите, что вас сюда привело.",
	"Добро пожаловать в наше скромное коммунити, %s. Представьтес, пожалуйста, и расскажите, что вас сюда привело.",
	"Доброе утро, день или вечер, %s. Представьтес, пожалуйста, и расскажите, что вас сюда привело.",
	"Наше вам с кисточкой тут, на канальчике, %s. Представьтес, пожалуйста, и расскажите, что вас сюда привело.",
}

// displayName builds the name substituted into a greeting for the given user.
// Preference order: "FirstName LastName", then LastName alone, then FirstName
// alone, then the username.
func displayName(u *echotron.User) string {
	if u == nil {
		return ""
	}

	name := strings.TrimSpace(strings.Join([]string{u.FirstName, u.LastName}, " "))

	if name != "" {
		return name
	}

	return u.Username
}

// sendGreeting sends a randomly chosen greeting for the given newcomer into the chat.
//
// It returns the id of the greeting it sent and whether it was actually
// delivered. Both matter: a newcomer who never received the question must not be
// judged on their silence, and a newcomer who is banned later must have the
// question taken back down.
//
// The mention makes the greeting a clickable link to the newcomer, and the
// message is sent in HTML mode because that is the only mode in which an
// arbitrary display name can be embedded safely — see greetMention.
func sendGreeting(chatid int64, u *echotron.User) (int, bool) {
	// Everything below dereferences u, including the error paths that name the
	// user, so an entity with no id is refused here rather than crashing later.
	if !usableUser(u) {
		return 0, false
	}

	// math/rand is sufficient here: we pick one of four static templates, not a
	// secret or a security-relevant value.
	greeting := greetings[rand.IntN(len(greetings))] //nolint: gosec // G404

	text := strings.Replace(greeting, "%s", greetingAddressee(u), 1)

	res, err := Tg.SendMessage(text, chatid, &echotron.MessageOptions{
		ParseMode: echotron.HTML,
	})
	if err != nil {
		mylog.Errorf("GREETING NOT DELIVERED to user %d in chat %d: %s", u.ID, chatid, err)

		return 0, false
	}

	mylog.Debugf("Greeted user %d in chat %d", u.ID, chatid)

	// A response carrying no message still means the greeting was delivered; we
	// simply have no id to retract later.
	return sentMessageID(res), true
}

// greetingAddressee is the name the template's %s is replaced with: a mention of
// the newcomer when they can be mentioned, and their plain display name when they
// cannot, so an unmentionable user is still greeted rather than skipped.
func greetingAddressee(u *echotron.User) string {
	if mention := greetMention(u); mention != "" {
		return mention
	}

	return displayName(u)
}

// sentMessageID pulls the id out of a send response, tolerating a response that
// carries no message at all.
func sentMessageID(res echotron.APIResponseMessage) int {
	if res.Result == nil {
		return 0
	}

	return res.Result.ID
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
