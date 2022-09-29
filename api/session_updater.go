package api

import (
	"net/http"
	"maid/api/rest"
	"maid/util"
	"time"
)

type X19SessionUpdater struct {
	hasFinalized bool
}

func (u *X19SessionUpdater) StartUpdating(client *http.Client, userAgent string, release rest.X19ReleaseInfo, user *util.X19User) {
	u.hasFinalized = false
	for !u.hasFinalized {
		time.Sleep(1 * time.Minute)
		go func() {
			err := rest.X19DoAuthenticationUpdate(client, userAgent, release, user)
			if err != nil {
				u.hasFinalized = true
			}
		}()
	}
}

func (u *X19SessionUpdater) Finalize() {
	u.hasFinalized = true
}