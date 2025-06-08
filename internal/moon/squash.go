package moon

import (
	"fmt"
	"slices"
	"time"

	"github.com/NicoNex/echotron/v3"
)

// TODO: refactor to return bool, error

// squash bans given chat member.
func squash(chatid int64, chattype, chattitle string, userid int64, firstname, lastname, username string) {
	banUntilDate := time.Now().Unix() + 10

	// If there is no such chatid, means that no such cache too.
	if len(ChatList) == 0 || !slices.Contains(ChatList, chatid) {
		Log.Debugf("Add %s %s (%d) to list of known chats", chattype, chattitle, chatid)

		ChatList = append(ChatList, chatid)
	}

	key := fmt.Sprintf("%d+%d", chatid, userid)
	if _, exist := SquashedMembers.Get(key); exist {
		Log.Infof(
			"Skip banning user %s %s (username = %s, id = %d) in %s %s (%d), already banned",
			firstname,
			lastname,
			username,
			userid,
			chattype,
			chattitle,
			chatid,
		)

		return
	}

	if res, err := Tg.BanChatMember(
		chatid,
		userid,
		&echotron.BanOptions{
			UntilDate: int(banUntilDate),
		},
	); err != nil {
		Log.Errorf(
			"Unable to ban user user %s %s (username = %s, id = %d) in %s %s (%d): %s",
			firstname,
			lastname,
			username,
			userid,
			chattype,
			chattitle,
			chatid,
			res.Description,
		)
	} else {
		if !res.Ok {
			Log.Errorf(
				"Unable to ban user %s %s (username = %s, id = %d) in %s %s (%d): %s",
				firstname,
				lastname,
				username,
				userid,
				chattype,
				chattitle,
				chatid,
				res.Description,
			)
		} else {
			Log.Debugf(
				"Banned user %s %s (username = %s, id = %d) in %s %s (%d).",
				firstname,
				lastname,
				username,
				userid,
				chattype,
				chattitle,
				chatid,
			)

			key := fmt.Sprintf("%d+%d", chatid, userid)
			Log.Debugf(
				"Put user %s %s (username = %s, id = %d) in %s %s (%d) into SquashedMembers to debounce ban.",
				firstname,
				lastname,
				username,
				userid,
				chattype,
				chattitle,
				chatid,
			)
			SquashedMembers.Set(key, true, 30*time.Second)

			Log.Debugf(
				"Remove user %s %s (username = %s, id = %d) in %s %s (%d) from NewMembers.",
				firstname,
				lastname,
				username,
				userid,
				chattype,
				chattitle,
				chatid,
			)
			NewMembers.Delete(key)

			Log.Debugf(
				"Remove user %s %s (username = %s, id = %d) in %s %s (%d) from AppearedMembers.",
				firstname,
				lastname,
				username,
				userid,
				chattype,
				chattitle,
				chatid,
			)
			AppearedMembers.Delete(key)

			if err := Cfg.SaveKeyValue(
				"BannedMembers",
				fmt.Sprintf("%d", chatid),
				fmt.Sprintf("%d", userid),
				fmt.Sprintf("%d", time.Now().Unix()),
			); err != nil {
				Log.Error(err)
			} else {
				Log.Infof("Info about banned userid %d saved to BannedMembers db", userid)
			}
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
