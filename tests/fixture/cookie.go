package fixture

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

func GetCookie(key interface{}, value interface{}) *http.Cookie {
	codecs := securecookie.CodecsFromPairs([]byte("secret"))
	data := make(map[interface{}]interface{})
	data[key] = value
	encoded, _ := securecookie.EncodeMulti("go-google-scraper", data, codecs...)

	cookie := http.Cookie{
		Name:  "go-google-scraper",
		Value: encoded,
	}

	return &cookie
}
