package fixture

import (
	"net/http"

	"github.com/gutakk/go-google-scraper/helpers/log"

	"github.com/gorilla/securecookie"
)

func GenerateCookie(key interface{}, value interface{}) *http.Cookie {
	codecs := securecookie.CodecsFromPairs([]byte("secret"))
	data := make(map[interface{}]interface{})
	data[key] = value
	encoded, err := securecookie.EncodeMulti("go-google-scraper", data, codecs...)
	if err != nil {
		log.Error("Failed to encode multi: ", err)
	}

	cookie := http.Cookie{
		Name:  "go-google-scraper",
		Value: encoded,
	}

	return &cookie
}
