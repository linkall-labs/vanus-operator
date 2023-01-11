// Code generated by go-swagger; DO NOT EDIT.

package cluster

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/linkall-labs/vanus-operator/api/models"
)

// PatchClusterHandlerFunc turns a function with the right signature into a patch cluster handler
type PatchClusterHandlerFunc func(PatchClusterParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PatchClusterHandlerFunc) Handle(params PatchClusterParams) middleware.Responder {
	return fn(params)
}

// PatchClusterHandler interface for that can handle valid patch cluster params
type PatchClusterHandler interface {
	Handle(PatchClusterParams) middleware.Responder
}

// NewPatchCluster creates a new http.Handler for the patch cluster operation
func NewPatchCluster(ctx *middleware.Context, handler PatchClusterHandler) *PatchCluster {
	return &PatchCluster{Context: ctx, Handler: handler}
}

/* PatchCluster swagger:route PATCH /cluster/ cluster patchCluster

集群手动伸缩

*/
type PatchCluster struct {
	Context *middleware.Context
	Handler PatchClusterHandler
}

func (o *PatchCluster) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPatchClusterParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PatchClusterOKBody patch cluster o k body
//
// swagger:model PatchClusterOKBody
type PatchClusterOKBody struct {

	// code
	// Required: true
	Code *int32 `json:"code"`

	// data
	// Required: true
	Data *models.ClusterInfo `json:"data"`

	// message
	// Required: true
	Message *string `json:"message"`
}

// Validate validates this patch cluster o k body
func (o *PatchClusterOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateCode(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateMessage(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PatchClusterOKBody) validateCode(formats strfmt.Registry) error {

	if err := validate.Required("patchClusterOK"+"."+"code", "body", o.Code); err != nil {
		return err
	}

	return nil
}

func (o *PatchClusterOKBody) validateData(formats strfmt.Registry) error {

	if err := validate.Required("patchClusterOK"+"."+"data", "body", o.Data); err != nil {
		return err
	}

	if o.Data != nil {
		if err := o.Data.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("patchClusterOK" + "." + "data")
			}
			return err
		}
	}

	return nil
}

func (o *PatchClusterOKBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("patchClusterOK"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this patch cluster o k body based on the context it is used
func (o *PatchClusterOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PatchClusterOKBody) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	if o.Data != nil {
		if err := o.Data.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("patchClusterOK" + "." + "data")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PatchClusterOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PatchClusterOKBody) UnmarshalBinary(b []byte) error {
	var res PatchClusterOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
