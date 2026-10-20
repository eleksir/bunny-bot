package elog

import (
	"fmt"
	"strings"

	"github.com/NicoNex/echotron/v3"
)

// renderLine assembles one log line:
//
//	21:33:07.123 in chat -1001 (title=…) user 42 (first_name=…) says …
//
// Empty field groups are dropped entirely rather than printed as empty parens.
func renderLine(chat echotron.Chat, user *echotron.User, body string) string {
	var b strings.Builder

	b.WriteString(nowUTC().Format(stampLayout))
	b.WriteString(" in chat ")
	b.WriteString(itoa(chat.ID))

	if desc := describeChat(&chat); desc != "" {
		b.WriteString(" (")
		b.WriteString(desc)
		b.WriteString(")")
	}

	if user != nil {
		b.WriteString(" user ")
		b.WriteString(itoa(user.ID))

		if desc := describeUser(user); desc != "" {
			b.WriteString(" (")
			b.WriteString(desc)
			b.WriteString(")")
		}
	}

	if body != "" {
		b.WriteString(" ")
		b.WriteString(body)
	}

	return b.String()
}

// describeChat renders a chat's identity fields. Type is deliberately omitted:
// the user-specified format has no slot for it, and which of title / first_name
// is present already distinguishes group from private chat.
func describeChat(c *echotron.Chat) string {
	if c == nil {
		return ""
	}

	return joinFields([]string{
		field("title", c.Title),
		field("username", atUsername(c.Username)),
		field("first_name", c.FirstName),
		field("last_name", c.LastName),
	})
}

// describeUser renders a user's identity fields.
func describeUser(u *echotron.User) string {
	if u == nil {
		return ""
	}

	return joinFields([]string{
		field("first_name", u.FirstName),
		field("last_name", u.LastName),
		field("username", atUsername(u.Username)),
	})
}

// describePayload renders what a message carries: media first, then a caption,
// because "post photo (…) caption" reads better than the reverse.
func describePayload(m *echotron.Message) string {
	kind, details := describeMedia(m)

	var b strings.Builder

	if kind != "" {
		b.WriteString("post ")
		b.WriteString(kind)

		if details != "" {
			b.WriteString(" (")
			b.WriteString(details)
			b.WriteString(")")
		}
	}

	// A caption accompanies the media rather than replacing it: "post photo (…)"
	// alone loses what the user actually said.
	if m.Text != "" {
		if b.Len() > 0 {
			b.WriteString(" ")
		}

		b.WriteString("says ")
		b.WriteString(m.Text)
	}

	return b.String()
}

// describeMedia returns the media kind and its key info. The first media field
// set on the message wins, in the order listed in the SPEC.
func describeMedia(m *echotron.Message) (string, string) {
	switch {
	case len(m.Photo) > 0:
		// The array is ordered smallest first, so the last entry is the largest.
		best := m.Photo[len(m.Photo)-1]

		return "photo", joinFields([]string{
			dims(best.Width, best.Height),
			sizeField(int64(best.FileSize)),
		})

	case m.Video != nil:
		return "video", joinFields([]string{
			field("file_name", m.Video.FileName),
			field("mime", m.Video.MimeType),
			dims(m.Video.Width, m.Video.Height),
			field("duration", seconds(m.Video.Duration)),
			sizeField(m.Video.FileSize),
		})

	case m.Animation != nil:
		return "animation", joinFields([]string{
			field("file_name", m.Animation.FileName),
			field("mime", m.Animation.MimeType),
			dims(m.Animation.Width, m.Animation.Height),
			field("duration", seconds(m.Animation.Duration)),
			sizeField(m.Animation.FileSize),
		})

	case m.Audio != nil:
		return "audio", joinFields([]string{
			field("title", m.Audio.Title),
			field("file_name", m.Audio.FileName),
			field("mime", m.Audio.MimeType),
			field("duration", seconds(m.Audio.Duration)),
			sizeField(m.Audio.FileSize),
		})

	case m.Voice != nil:
		return "voice", joinFields([]string{
			field("mime", m.Voice.MimeType),
			field("duration", seconds(m.Voice.Duration)),
			sizeField(m.Voice.FileSize),
		})

	case m.VideoNote != nil:
		return "video_note", joinFields([]string{
			dims(m.VideoNote.Length, m.VideoNote.Length),
			field("duration", seconds(m.VideoNote.Duration)),
			sizeField(int64(m.VideoNote.FileSize)),
		})

	case m.Sticker != nil:
		return "sticker", joinFields([]string{
			field("emoji", m.Sticker.Emoji),
			field("set", m.Sticker.SetName),
			dims(m.Sticker.Width, m.Sticker.Height),
			field("video", itoaBool(m.Sticker.IsVideo)),
			field("animated", itoaBool(m.Sticker.IsAnimated)),
		})

	case m.Document != nil:
		return "document", joinFields([]string{
			field("file_name", m.Document.FileName),
			field("mime", m.Document.MimeType),
			sizeField(m.Document.FileSize),
		})

	case m.Contact != nil:
		return "contact", joinFields([]string{
			field("phone", m.Contact.PhoneNumber),
			field("first_name", m.Contact.FirstName),
			field("last_name", m.Contact.LastName),
			field("user_id", itoaInt(m.Contact.UserID)),
		})

	case m.Location != nil:
		return "location", joinFields([]string{
			field("lat", itoaFloat(m.Location.Latitude)),
			field("lon", itoaFloat(m.Location.Longitude)),
		})

	case m.Venue != nil:
		parts := []string{
			field("title", m.Venue.Title),
			field("address", m.Venue.Address),
		}

		if m.Venue.Location != nil {
			parts = append(parts,
				field("lat", itoaFloat(m.Venue.Location.Latitude)),
				field("lon", itoaFloat(m.Venue.Location.Longitude)))
		}

		return "venue", joinFields(parts)

	case m.Dice != nil:
		return "dice", joinFields([]string{
			field("emoji", m.Dice.Emoji),
			field("value", itoaInt(m.Dice.Value)),
		})
	}

	return "", ""
}

// dims renders "WxH", or empty when the dimensions are unknown.
func dims(w int, h int) string {
	if w <= 0 || h <= 0 {
		return ""
	}

	return itoaInt(w) + "x" + itoaInt(h)
}

// itoaBool renders a boolean for an inline field.
func itoaBool(b bool) string {
	if b {
		return "true"
	}

	return "false"
}

// describeReference renders "user 42 (first_name=Ann)", used for the newcomer in
// a service message and for the member in a status change.
func describeReference(u *echotron.User) string {
	if u == nil {
		return "user unknown"
	}

	if desc := describeUser(u); desc != "" {
		return fmt.Sprintf("user %d (%s)", u.ID, desc)
	}

	return fmt.Sprintf("user %d", u.ID)
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
