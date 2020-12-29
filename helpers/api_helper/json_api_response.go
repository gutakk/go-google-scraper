package api_helper

import "github.com/gin-gonic/gin"

type ErrorResponseObject struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
}

func (e *ErrorResponseObject) ConstructErrorResponse() gin.H {
	errorResponse := []ErrorResponseObject{{
		Title:  e.Title,
		Detail: e.Detail,
		Status: e.Status,
	}}
	return gin.H{
		"errors": errorResponse,
	}
}
