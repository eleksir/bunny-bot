package moon

import (
	"encoding/json"
	"time"

	"github.com/NicoNex/echotron/v3"
)

// AddUser remebers user info in People cache.
func AddUser(u *echotron.User) {
	bytes, err := json.Marshal(u)

	if err != nil {
		Log.Errorf("Unable to encode user info to json string: %s", err)
	}

	People.Set(u.ID, bytes, 24*time.Hour)
}

// GetUser returns user info structure.
func GetUser(id int64) *echotron.User {
	var u *echotron.User

	bytes, ok := People.Get(id)

	if !ok {
		return u
	}

	err := json.Unmarshal(bytes.([]byte), &u)

	if err != nil {
		Log.Errorf("Unable to decode user info json: %s", err)

		return u
	}

	return u
}

// AddChat remebers chat info in Groups cache.
func AddChat(c *echotron.Chat) {
	bytes, err := json.Marshal(c)

	if err != nil {
		Log.Errorf("Unable to encode chat info to json string: %s", err)
	}

	Groups.Set(c.ID, bytes, 24*time.Hour)
}

// GetChat returns chat info structure.
func GetChat(id int64) *echotron.Chat {
	var c *echotron.Chat

	bytes, ok := Groups.Get(id)

	if !ok {
		Log.Errorf("No group %d in Groups cache", id)

		return c
	}

	err := json.Unmarshal(bytes.([]byte), &c)

	if err != nil {
		Log.Errorf("Unable to decode chat info json: %s", err)

		return c
	}

	return c
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
