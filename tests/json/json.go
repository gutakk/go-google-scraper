package json

import (
	"encoding/json"

	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/helpers/log"
)

func JSONMarshaler(value interface{}) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		log.Fatal(errorconf.JSONMarshalFailure, err)
	}

	return data
}

func JSONUnmarshaler(data []byte, v interface{}) {
	err := json.Unmarshal(data, &v)
	if err != nil {
		log.Fatal(errorconf.JSONUnmarshalFailure, err)
	}
}
