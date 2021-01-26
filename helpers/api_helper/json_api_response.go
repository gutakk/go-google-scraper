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
	ID            string      `json:"id"`
	Type          string      `json:"type"`
	Attributes    interface{} `json:"attributes,omitempty"`
	Relationships interface{} `json:"relationships,omitempty"`
}

type DataResponse struct {
	Data DataResponseObject `json:"data"`
}

type DataResponseArray struct {
	Data []DataResponseObject `json:"data"`
}

func (d *DataResponseObject) GetRelationships() map[string]DataResponse {
	relationships := make(map[string]DataResponse)
	relationships[d.Type] = DataResponse{
		Data: *d,
	}

	return relationships
}
