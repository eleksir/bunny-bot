package moon

import (
	"time"
)

func RunAndBan() {
	for {
		// If user have both events, he is not spammer. So wipe him from unchecked events list, just for sure.
		for _, chatid := range ChatList {
			Log.Infof("RunAndBan by chatid: %d", chatid)

			var ListAppearedMembers, ListNewMembers []int64

			AppearedMembers[chatid].Range(
				func(k, _ any) bool {
					ListAppearedMembers = append(ListAppearedMembers, k.(int64))

					return true
				},
			)

			NewMembers[chatid].Range(
				func(k, _ any) bool {
					ListNewMembers = append(ListNewMembers, k.(int64))

					return true
				},
			)

			for _, AppearedMember := range ListAppearedMembers {
				for _, NewMember := range ListNewMembers {
					if AppearedMember == NewMember {
						// TODO: Store also additional attributes in temporary database.
						Log.Infof(
							"Skip banning user %s %s (username = %s, id = %d) in %s %s (%d), because it has both events newMember and updateMember",
							"",
							"",
							"",
							NewMember,
							"",
							"",
							chatid,
						)

						Log.Infof(
							"Deleting user %s %s (username = %s, id = %d) in %s %s (%d), from AppearedMember DB",
							"",
							"",
							"",
							NewMember,
							"",
							"",
							chatid,
						)

						AppearedMembers[chatid].Delete(AppearedMember)

						Log.Infof(
							"Deleting user %s %s (username = %s, id = %d) in %s %s (%d), from NewMember DB",
							"",
							"",
							"",
							NewMember,
							"",
							"",
							chatid,
						)

						NewMembers[chatid].Delete(NewMember)
					}
				}
			}
		}

		currentTime := time.Now().Unix()

		// We would like to process only events that definitely older than 10 seconds!
		for _, chatid := range ChatList {
			var (
				ListAppearedMembers, ListNewMembers []int64
			)

			AppearedMembers[chatid].Range(
				func(k, v interface{}) bool {
					ListAppearedMembers = append(ListAppearedMembers, k.(int64))

					return true
				},
			)

			NewMembers[chatid].Range(
				func(k, v interface{}) bool {
					ListNewMembers = append(ListNewMembers, k.(int64))

					return true
				},
			)

			for _, AppearedMember := range ListAppearedMembers {
				eventTimestamp, found := AppearedMembers[chatid].Get(AppearedMember)

				if !found {
					continue
				}

				eventTimeDelta := currentTime - eventTimestamp.(int64)

				// Looks like it is spammer!
				if eventTimeDelta > 10 {
					// TODO: Store also additional attributes in temporary database.
					Log.Infof(
						"Ban user %s %s (username = %s, id = %d) in chat %s %s (%d) that have no newChatMember message since chatMemberUpdated for more than 10 seconds",
						"",
						"",
						"",
						AppearedMember,
						"",
						"",
						chatid,
					)

					squash(
						chatid,
						"",
						"",
						AppearedMember,
						"",
						"",
						"",
					)

					AppearedMembers[chatid].Delete(AppearedMember)
				}
			}

			for _, NewMember := range ListNewMembers {
				eventTimestamp, found := NewMembers[chatid].Get(NewMember)

				if !found {
					continue
				}

				eventTimeDelta := currentTime - eventTimestamp.(int64)

				// Totally new members should not have chatMemberUpdated events.
				if eventTimeDelta > 15 {
					// TODO: Store also additional attributes in temporary database.
					Log.Debugf("Clean newChatMember event record older than 10 seconds for user user %s %s (username = %s, id = %d) in chat %s %s (%d)",
						"",
						"",
						"",
						NewMember,
						"",
						"",
						chatid,
					)

					NewMembers[chatid].Delete(NewMember)
				}
			}
		}

		time.Sleep(10 * time.Second)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
