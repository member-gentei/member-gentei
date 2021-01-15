// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi"
)

// ChannelMember defines model for ChannelMember.
type ChannelMember struct {
	Id *string `json:"id,omitempty"`
}

// AfterParam defines model for afterParam.
type AfterParam string

// ChannelSlugPathParam defines model for channelSlugPathParam.
type ChannelSlugPathParam string

// LimitParam defines model for limitParam.
type LimitParam int

// GetMembersParams defines parameters for GetMembers.
type GetMembersParams struct {

	// Specific Discord snowflake(s) to retrieve - omit to get all members.
	Snowflakes *[]string `json:"snowflakes,omitempty"`

	// Maximum page size.
	Limit *LimitParam `json:"limit,omitempty"`

	// Pagination cursor. Put the `next` string here if it's present in order to fetch the next page of the response.
	After *AfterParam `json:"after,omitempty"`
}

// CheckMembershipJSONBody defines parameters for CheckMembership.
type CheckMembershipJSONBody struct {

	// Discord snowflake ID
	Snowflake string `json:"snowflake"`
}

// CheckMembershipRequestBody defines body for CheckMembership for application/json ContentType.
type CheckMembershipJSONRequestBody CheckMembershipJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get members of a YouTube channel
	// (GET /v1/channel/{channelSlug}/members)
	GetMembers(w http.ResponseWriter, r *http.Request, channelSlug ChannelSlugPathParam, params GetMembersParams)
	// Initiate a check for Discord user membership
	// (POST /v1/channel/{channelSlug}/members/check)
	CheckMembership(w http.ResponseWriter, r *http.Request, channelSlug ChannelSlugPathParam)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetMembers operation middleware
func (siw *ServerInterfaceWrapper) GetMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "channelSlug" -------------
	var channelSlug ChannelSlugPathParam

	err = runtime.BindStyledParameter("simple", false, "channelSlug", chi.URLParam(r, "channelSlug"), &channelSlug)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter channelSlug: %s", err), http.StatusBadRequest)
		return
	}

	ctx = context.WithValue(ctx, "APIKey.Scopes", []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetMembersParams

	// ------------- Optional query parameter "snowflakes" -------------
	if paramValue := r.URL.Query().Get("snowflakes"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "snowflakes", r.URL.Query(), &params.Snowflakes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter snowflakes: %s", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "limit" -------------
	if paramValue := r.URL.Query().Get("limit"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter limit: %s", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "after" -------------
	if paramValue := r.URL.Query().Get("after"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "after", r.URL.Query(), &params.After)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter after: %s", err), http.StatusBadRequest)
		return
	}

	siw.Handler.GetMembers(w, r.WithContext(ctx), channelSlug, params)
}

// CheckMembership operation middleware
func (siw *ServerInterfaceWrapper) CheckMembership(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "channelSlug" -------------
	var channelSlug ChannelSlugPathParam

	err = runtime.BindStyledParameter("simple", false, "channelSlug", chi.URLParam(r, "channelSlug"), &channelSlug)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter channelSlug: %s", err), http.StatusBadRequest)
		return
	}

	ctx = context.WithValue(ctx, "APIKey.Scopes", []string{""})

	siw.Handler.CheckMembership(w, r.WithContext(ctx), channelSlug)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerFromMux(si, chi.NewRouter())
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	r.Group(func(r chi.Router) {
		r.Get("/v1/channel/{channelSlug}/members", wrapper.GetMembers)
	})
	r.Group(func(r chi.Router) {
		r.Post("/v1/channel/{channelSlug}/members/check", wrapper.CheckMembership)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
