package moon

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

// RunAndBan watches that noone arrives without newChatMember message.
func RunAndBan() {
	for {
		var (
			ListAppearedMembers, ListNewMembers []string
			empty                               string
		)

		AppearedMembers.Range(
			func(k, _ any) bool {
				ListAppearedMembers = append(ListAppearedMembers, k.(string))

				return true
			},
		)

		NewMembers.Range(
			func(k, _ any) bool {
				ListNewMembers = append(ListNewMembers, k.(string))

				return true
			},
		)

		for _, AppearedMember := range ListAppearedMembers {
			for _, NewMember := range ListNewMembers {
				if AppearedMember == NewMember {
					AppearedMembers.Delete(AppearedMember)
					NewMembers.Delete(NewMember)
				}
			}
		}

		currentTime := time.Now().Unix()

		ListAppearedMembers = ListAppearedMembers[:0]
		ListNewMembers = ListNewMembers[:0]
		ListAppearedMembers = slices.Clip(ListAppearedMembers)
		ListNewMembers = slices.Clip(ListNewMembers)

		AppearedMembers.Range(
			func(k, _ any) bool {
				ListAppearedMembers = append(ListAppearedMembers, k.(string))

				return true
			},
		)

		NewMembers.Range(
			func(k, _ any) bool {
				ListNewMembers = append(ListNewMembers, k.(string))

				return true
			},
		)

		for _, AppearedMember := range ListAppearedMembers {
			eventTimestamp, found := AppearedMembers.Get(AppearedMember)

			if !found {
				continue
			}

			eventTimeDelta := currentTime - eventTimestamp.(int64)

			// Spammer?
			if eventTimeDelta > 10 {
				item := strings.Split(AppearedMember, "+")
				chatid, err := strconv.ParseInt(item[0], 10, 64)

				if err != nil {
					Log.Errorf("Got broken key in AppearedMembers cache: %s", AppearedMember)

					continue
				}

				userid, err := strconv.ParseInt(item[1], 10, 64)

				if err != nil {
					Log.Errorf("Got broken key in AppearedMembers cache: %s", AppearedMember)

					continue
				}

				Log.Infof(
					"Ban user %s %s (username = %s, id = %d) in chat %s %s (%d) that have no newChatMember message since chatMemberUpdated for more than 10 seconds",
					empty,
					empty,
					empty,
					userid,
					empty,
					empty,
					chatid,
				)

				squash(
					chatid,
					empty,
					empty,
					userid,
					empty,
					empty,
					empty,
				)

				AppearedMembers.Delete(AppearedMember)
			}
		}

		for _, NewMember := range ListNewMembers {
			eventTimestamp, found := NewMembers.Get(NewMember)

			if !found {
				continue
			}

			eventTimeDelta := currentTime - eventTimestamp.(int64)

			// Spammer?
			if eventTimeDelta > 10 {
				item := strings.Split(NewMember, "+")
				chatid, err := strconv.ParseInt(item[0], 10, 64)

				if err != nil {
					Log.Errorf("Got broken key in AppearedMembers cache: %s", NewMember)

					continue
				}

				userid, err := strconv.ParseInt(item[1], 10, 64)

				if err != nil {
					Log.Errorf("Got broken key in AppearedMembers cache: %s", NewMember)

					continue
				}

				Log.Debugf("Clean newChatMember event record older than 10 seconds for user user %s %s (username = %s, id = %d) in chat %s %s (%d)",
					empty,
					empty,
					empty,
					userid,
					empty,
					empty,
					chatid,
				)

				NewMembers.Delete(NewMember)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

// Runs deferred check vs CAS db for recently joined users.
func DeferredCASCheck() {
	for {
		var (
			ListPendingMembers []string
			empty              string
		)

		PendingCASMembers.Range(
			func(k, _ any) bool {
				ListPendingMembers = append(ListPendingMembers, k.(string))

				return true
			},
		)

		currentTime := time.Now().Unix()

		for _, PendingMember := range ListPendingMembers {
			eventTimestamp, found := PendingCASMembers.Get(PendingMember)

			if !found {
				continue
			}

			eventTimeDelta := currentTime - eventTimestamp.(int64)

			// It's time to check for second time this id if it was banned in CAS.
			if eventTimeDelta > 300 {
				item := strings.Split(PendingMember, "+")
				chatid, err := strconv.ParseInt(item[0], 10, 64)

				if err != nil {
					Log.Errorf("Got broken key in PendingMembers cache: %s", PendingMember)

					continue
				}

				userid, err := strconv.ParseInt(item[1], 10, 64)

				if err != nil {
					Log.Errorf("Got broken key in PendingMembers cache: %s", PendingMember)

					continue
				}

				value := Cfg.GetValue(
					"BannedMembers",
					fmt.Sprintf("%d", chatid),
					fmt.Sprintf("%d", userid),
				)

				if value != "" {
					Log.Debugf("User %s %s (username = %s, id = %d) in %s %s (%d) is already banned, skip cas check.",
						empty,
						empty,
						empty,
						userid,
						empty,
						empty,
						chatid,
					)

					Log.Debugf(
						"Remove user %s %s (username = %s, id = %d) in %s %s (%d) from PendingCASMembers",
						empty,
						empty,
						empty,
						userid,
						empty,
						empty,
						chatid,
					)

					key := fmt.Sprintf("%d+%d", chatid, userid)
					PendingCASMembers.Delete(key)

					return
				}

				Log.Debugf("Checking user %s %s (username = %s, id = %d) in %s %s (%d) for second time in CAS",
					empty,
					empty,
					empty,
					userid,
					empty,
					empty,
					chatid,
				)

				banned, err := CasCheckID(userid)

				if err != nil {
					Log.Error(err)
				}

				if banned {
					Log.Infof(
						"Ban user %s %s (username = %s, id = %d) in %s %s (%d) by CAS blacklist",
						empty,
						empty,
						empty,
						userid,
						empty,
						empty,
						chatid,
					)

					squash(
						userid,
						empty,
						empty,
						userid,
						empty,
						empty,
						empty,
					)
				}

				Log.Debugf(
					"Remove user %s %s (username = %s, id = %d) in %s %s (%d) from PendingCASMembers",
					empty,
					empty,
					empty,
					userid,
					empty,
					empty,
					chatid,
				)

				PendingCASMembers.Delete(PendingMember)
			}

		}

		time.Sleep(10 * time.Minute)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
