// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetScheduleHandlerFunc turns a function with the right signature into a get schedule handler
type GetScheduleHandlerFunc func(GetScheduleParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetScheduleHandlerFunc) Handle(params GetScheduleParams) middleware.Responder {
	return fn(params)
}

// GetScheduleHandler interface for that can handle valid get schedule params
type GetScheduleHandler interface {
	Handle(GetScheduleParams) middleware.Responder
}

// NewGetSchedule creates a new http.Handler for the get schedule operation
func NewGetSchedule(ctx *middleware.Context, handler GetScheduleHandler) *GetSchedule {
	return &GetSchedule{Context: ctx, Handler: handler}
}

/*
	GetSchedule swagger:route GET /schedule getSchedule

# Get a specific schedule

Retrieve a schedule by user ID and schedule ID
*/
type GetSchedule struct {
	Context *middleware.Context
	Handler GetScheduleHandler
}

func (o *GetSchedule) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetScheduleParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
