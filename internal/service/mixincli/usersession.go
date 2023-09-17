package mixincli

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	// Set the default expiration to 10 minutes
	SessionCache = cache.New(time.Minute*10, cache.NoExpiration)
)

const (
	UserSessionStateInit = iota
	UserSessionStateWaiting
	UserSessionStateSuccess
)

type UserSession struct {
	State int
}

func setSession(userID string, sess *UserSession) {
	SessionCache.Set(userID, sess, cache.DefaultExpiration)
}

func getSession(userID string) *UserSession {
	if x, found := SessionCache.Get(userID); found {
		return x.(*UserSession)
	}
	return nil
}
