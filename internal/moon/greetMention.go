package moon

import (
	"strconv"
	"strings"

	"github.com/NicoNex/echotron/v3"
)

// greetMention renders a clickable mention of a newcomer for use inside a
// greeting, or "" when the user cannot be mentioned at all.
//
// The mode is HTML, not markdown, and that is a correctness decision rather than
// a preference: the label is a name Telegram gave us, and MarkdownV2 requires
// escaping eighteen characters — "-" among them — so an ordinary "Анна-Мария"
// would be rejected with a 400. A rejected greeting makes newChatMembers delete
// the join event, which is how one hyphen would disable the only detector. HTML
// needs "&", "<" and ">" escaped, and those essentially never occur in a name.
//
// The mention works without a @username, which is the common case for someone who
// has only just joined a public group.
func greetMention(u *echotron.User) string {
	if u == nil || u.ID <= 0 {
		return ""
	}

	var b strings.Builder

	b.WriteString(`<a href="tg://user?id=`)
	b.WriteString(strconv.FormatInt(u.ID, 10))
	b.WriteString(`">`)

	for _, r := range mentionLabel(u) {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		default:
			b.WriteRune(r)
		}
	}

	b.WriteString("</a>")

	return b.String()
}

// escapeHTML escapes the three characters HTML requires. It exists for callers
// that assemble markup themselves; the greeting path uses greetMention's own
// single pass so the anchor and the label cannot drift apart.
func escapeHTML(s string) string {
	return htmlReplacer.Replace(s)
}

// mentionLabel is the text shown inside the mention link, in the preference order
// a reader expects: full name, then either name alone, then the @username, and
// finally the bare id so the link is never empty.
func mentionLabel(u *echotron.User) string {
	if name := strings.TrimSpace(u.FirstName + " " + u.LastName); name != "" {
		return name
	}

	if u.LastName != "" {
		return u.LastName
	}

	if u.FirstName != "" {
		return u.FirstName
	}

	if u.Username != "" {
		return "@" + u.Username
	}

	return strconv.FormatInt(u.ID, 10)
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
