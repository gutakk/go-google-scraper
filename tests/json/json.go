package json

import (
	"encoding/json"

	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/helpers/log"
)

func JSONMarshaler(value interface{}) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		log.Error(errorconf.JSONMarshalFailure, err)
	}

	return data
}
