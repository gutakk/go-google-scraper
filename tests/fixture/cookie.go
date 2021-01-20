package fixture

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/securecookie"
)

func GenerateCookie(key interface{}, value interface{}) *http.Cookie {
	codecs := securecookie.CodecsFromPairs([]byte("secret"))
	data := make(map[interface{}]interface{})
	data[key] = value
	encoded, encodeMultiErr := securecookie.EncodeMulti("go-google-scraper", data, codecs...)
	if encodeMultiErr != nil {
		glog.Errorf("Cannot encode multi: %s", encodeMultiErr)
	}

	cookie := http.Cookie{
		Name:  "go-google-scraper",
		Value: encoded,
	}

	return &cookie
}
