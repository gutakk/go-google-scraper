package oauth_test

import "time"

type TokenStoreItem struct {
	ID        int64
	CreatedAt time.Time
	ExpiresAt time.Time
	Code      string
	Access    string
	Refresh   string
	Data      []byte
}

type TokenData struct {
	ClientID         string
	UserID           string
	RedirectURI      string
	Scope            string
	Code             string
	CodeCreateAt     time.Time
	CodeExpiresIn    time.Duration
	Access           string
	AccessCreateAt   time.Time
	AccessExpiresIn  time.Duration
	Refresh          string
	RefreshCreateAt  time.Time
	RefreshExpiresIn time.Duration
}
