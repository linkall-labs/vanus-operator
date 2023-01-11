// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ClusterInfo 集群信息
//
// swagger:model Cluster_info
type ClusterInfo struct {

	// cloudevents端点
	CloudeventsEndpoints string `json:"cloudevents_endpoints,omitempty"`

	// gateway端点
	GatewayEndpoints string `json:"gateway_endpoints,omitempty"`
}

// Validate validates this cluster info
func (m *ClusterInfo) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this cluster info based on context it is used
func (m *ClusterInfo) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ClusterInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ClusterInfo) UnmarshalBinary(b []byte) error {
	var res ClusterInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}