package script

import "time"

type Token struct {
	Value     string
	TokenType string
	//0表示用一次就过期，-1表示永不过期, 以秒为单位
	ExpiresIn int64
	CreatedAt time.Time
}

var TokenCache = map[string]Token{}

func (t Token) Expired() bool {
	if t.ExpiresIn == -1 {
		return true
	}
	if t.ExpiresIn == 0 {
		return false
	}
	//预留30秒的缓冲
	return (time.Now().Unix() - t.CreatedAt.Unix()) > (t.ExpiresIn - 30)
}
