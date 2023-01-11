// Code generated by go-swagger; DO NOT EDIT.

package connector

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// CreateConnectorOKCode is the HTTP code returned for type CreateConnectorOK
const CreateConnectorOKCode int = 200

/*CreateConnectorOK OK

swagger:response createConnectorOK
*/
type CreateConnectorOK struct {

	/*
	  In: Body
	*/
	Payload *CreateConnectorOKBody `json:"body,omitempty"`
}

// NewCreateConnectorOK creates CreateConnectorOK with default headers values
func NewCreateConnectorOK() *CreateConnectorOK {

	return &CreateConnectorOK{}
}

// WithPayload adds the payload to the create connector o k response
func (o *CreateConnectorOK) WithPayload(payload *CreateConnectorOKBody) *CreateConnectorOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create connector o k response
func (o *CreateConnectorOK) SetPayload(payload *CreateConnectorOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateConnectorOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
