// Package payload provides utilities for dealing with HTTP request and response payloads.
// It integrates with sibling packages log and errors.
package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ClientReporter provides information about an error such that client and
// server errors can be distinguished and handled appropriately.
type ClientReporter interface {
	error
	Message() map[string]string
	Status() int
}

// WriteError writes an appropriate error response to the given response
// writer. If the given error implements ClientReport, then the values from
// ErrorReport() and StatusCode() are written to the response, except in
// the case of a 5XX error, where the error is logged and a default message is
// written to the response.
func WriteError(c *gin.Context, e error) {
	if cr, ok := e.(ClientReporter); ok {
		status := cr.Status()
		if status >= http.StatusInternalServerError {
			handleInternalServerError(c, e)
			return
		}

		// log.FromRequest(r).Print(cr.Error())
		Write(c, cr.Message(), status)

		return
	}

	handleInternalServerError(c, e)
}

// Write writes the given payload to the response. If the payload
// cannot be marshaled, a 500 error is written instead. If the writer
// cannot be written to, then this function panics.
func Write(c *gin.Context, payload interface{}, status int) {
	// encoded, err := json.Marshal(payload)
	// if err != nil {
	// 	handleInternalServerError(w, r, errors.E(op, err))
	// 	return
	// }
	// Set the appropriate headers to indicate gzip compression and JSON content type
	c.JSON(status, payload)
}

// handleInternalServerError writes the given error to stderr and returns a
// 500 response with a default message.
func handleInternalServerError(c *gin.Context, e error) {
	log.Printf("Err: %v", e)
	c.JSON(http.StatusInternalServerError, gin.H{"message": "Something has gone wrong"})
}
