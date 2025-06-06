// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/infinity-ocean/ikakbolit/3-api-grpc-Homework/server/models"
)

// PostScheduleCreatedCode is the HTTP code returned for type PostScheduleCreated
const PostScheduleCreatedCode int = 201

/*
PostScheduleCreated Created

swagger:response postScheduleCreated
*/
type PostScheduleCreated struct {

	/*
	  In: Body
	*/
	Payload *models.ControllerResponseScheduleID `json:"body,omitempty"`
}

// NewPostScheduleCreated creates PostScheduleCreated with default headers values
func NewPostScheduleCreated() *PostScheduleCreated {

	return &PostScheduleCreated{}
}

// WithPayload adds the payload to the post schedule created response
func (o *PostScheduleCreated) WithPayload(payload *models.ControllerResponseScheduleID) *PostScheduleCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post schedule created response
func (o *PostScheduleCreated) SetPayload(payload *models.ControllerResponseScheduleID) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostScheduleCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostScheduleBadRequestCode is the HTTP code returned for type PostScheduleBadRequest
const PostScheduleBadRequestCode int = 400

/*
PostScheduleBadRequest Bad request

swagger:response postScheduleBadRequest
*/
type PostScheduleBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.ControllerAPIError `json:"body,omitempty"`
}

// NewPostScheduleBadRequest creates PostScheduleBadRequest with default headers values
func NewPostScheduleBadRequest() *PostScheduleBadRequest {

	return &PostScheduleBadRequest{}
}

// WithPayload adds the payload to the post schedule bad request response
func (o *PostScheduleBadRequest) WithPayload(payload *models.ControllerAPIError) *PostScheduleBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post schedule bad request response
func (o *PostScheduleBadRequest) SetPayload(payload *models.ControllerAPIError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostScheduleBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostScheduleNotFoundCode is the HTTP code returned for type PostScheduleNotFound
const PostScheduleNotFoundCode int = 404

/*
PostScheduleNotFound Resource not found

swagger:response postScheduleNotFound
*/
type PostScheduleNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.ControllerAPIError `json:"body,omitempty"`
}

// NewPostScheduleNotFound creates PostScheduleNotFound with default headers values
func NewPostScheduleNotFound() *PostScheduleNotFound {

	return &PostScheduleNotFound{}
}

// WithPayload adds the payload to the post schedule not found response
func (o *PostScheduleNotFound) WithPayload(payload *models.ControllerAPIError) *PostScheduleNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post schedule not found response
func (o *PostScheduleNotFound) SetPayload(payload *models.ControllerAPIError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostScheduleNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostScheduleInternalServerErrorCode is the HTTP code returned for type PostScheduleInternalServerError
const PostScheduleInternalServerErrorCode int = 500

/*
PostScheduleInternalServerError Internal server error

swagger:response postScheduleInternalServerError
*/
type PostScheduleInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ControllerAPIError `json:"body,omitempty"`
}

// NewPostScheduleInternalServerError creates PostScheduleInternalServerError with default headers values
func NewPostScheduleInternalServerError() *PostScheduleInternalServerError {

	return &PostScheduleInternalServerError{}
}

// WithPayload adds the payload to the post schedule internal server error response
func (o *PostScheduleInternalServerError) WithPayload(payload *models.ControllerAPIError) *PostScheduleInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post schedule internal server error response
func (o *PostScheduleInternalServerError) SetPayload(payload *models.ControllerAPIError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostScheduleInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
