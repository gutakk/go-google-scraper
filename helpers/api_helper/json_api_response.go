package api_helper

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponseObject struct {
	Detail string `json:"detail,omitempty"`
	Status int    `json:"status,omitempty"`
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
	ID            string      `json:"id,omitempty"`
	Type          string      `json:"type,omitempty"`
	Attributes    interface{} `json:"attributes,omitempty"`
	Relationships interface{} `json:"relationships,omitempty"`
}

type DataResponse struct {
	Data DataResponseObject `json:"data,omitempty"`
}

type DataResponseArray struct {
	Data []DataResponseObject `json:"data,omitempty"`
}

func (d *DataResponseObject) GetRelationships() map[string]DataResponse {
	relationships := make(map[string]DataResponse)
	relationships[d.Type] = DataResponse{
		Data: *d,
	}

	return relationships
}
