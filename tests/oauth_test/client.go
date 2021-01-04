package oauth_test

import "gorm.io/datatypes"

type OAuthClient struct {
	ID     string
	Secret string
	Domain string
	Data   datatypes.JSON
}
