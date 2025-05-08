package moon

import (
	"slices"
	"time"

	"github.com/NicoNex/echotron/v3"
	"github.com/akyoto/cache"
)

// TODO: refactor to return bool, error

// squash bans given chat member.
func squash(chatid int64, chattype, chattitle string, userid int64, firstname, lastname, username string) {
	banUntilDate := time.Now().Unix() + 10

	// If there is no such chatid, means that no such cache too.
	if len(ChatList) == 0 || !slices.Contains(ChatList, chatid) {
		Log.Debugf("Creating caches for new %s %s (%d)", chattype, chattitle, chatid)

		NewMembers[chatid] = cache.New(1 * time.Minute)
		AppearedMembers[chatid] = cache.New(1 * time.Minute)
		SquashedMembers[chatid] = cache.New(5 * time.Minute)

		Log.Debugf("Add %s %s (%d) to list of known chats", chattype, chattitle, chatid)

		ChatList = append(ChatList, chatid)
	}

	if found, exist := SquashedMembers[chatid].Get(userid); exist && found.(bool) {
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

			SquashedMembers[chatid].Set(userid, true, 30*time.Second)
			NewMembers[chatid].Delete(userid)
			AppearedMembers[chatid].Delete(userid)
		}
	}
}
