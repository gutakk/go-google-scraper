package api_helper

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponseObject struct {
	Detail string `json:"detail"`
	Status int    `json:"status"`
}

func (e *ErrorResponseObject) NewErrorResponse() gin.H {
	errorResponse := []ErrorResponseObject{{
		Detail: e.Detail,
		Status: e.Status,
	}}
	return gin.H{
		"errors": errorResponse,
	}
}

type DataResponseObject struct {
	ID            string                           `json:"id"`
	Type          string                           `json:"type"`
	Attributes    interface{}                      `json:"attributes"`
	Relationships map[string]RelationshipsResponse `json:"relationships"`
}

type DataResponse struct {
	Data DataResponseObject `json:"data"`
}

type DataResponseList struct {
	Data []DataResponseObject `json:"data"`
}
