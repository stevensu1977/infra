// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.2 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (DELETE /envs/{codeSnippetID})
	DeleteEnvsCodeSnippetID(c *gin.Context, codeSnippetID CodeSnippetID, params DeleteEnvsCodeSnippetIDParams)

	// (PATCH /envs/{codeSnippetID})
	PatchEnvsCodeSnippetID(c *gin.Context, codeSnippetID CodeSnippetID, params PatchEnvsCodeSnippetIDParams)

	// (POST /envs/{codeSnippetID})
	PostEnvsCodeSnippetID(c *gin.Context, codeSnippetID CodeSnippetID, params PostEnvsCodeSnippetIDParams)

	// (PUT /envs/{codeSnippetID}/state)
	PutEnvsCodeSnippetIDState(c *gin.Context, codeSnippetID CodeSnippetID, params PutEnvsCodeSnippetIDStateParams)

	// (GET /health)
	GetHealth(c *gin.Context)

	// (GET /sessions)
	GetSessions(c *gin.Context, params GetSessionsParams)

	// (POST /sessions)
	PostSessions(c *gin.Context, params PostSessionsParams)

	// (DELETE /sessions/{sessionID})
	DeleteSessionsSessionID(c *gin.Context, sessionID SessionID, params DeleteSessionsSessionIDParams)

	// (POST /sessions/{sessionID}/refresh)
	PostSessionsSessionIDRefresh(c *gin.Context, sessionID SessionID, params PostSessionsSessionIDRefreshParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// DeleteEnvsCodeSnippetID operation middleware
func (siw *ServerInterfaceWrapper) DeleteEnvsCodeSnippetID(c *gin.Context) {

	var err error

	// ------------- Path parameter "codeSnippetID" -------------
	var codeSnippetID CodeSnippetID

	err = runtime.BindStyledParameter("simple", false, "codeSnippetID", c.Param("codeSnippetID"), &codeSnippetID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter codeSnippetID: %s", err), http.StatusBadRequest)
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params DeleteEnvsCodeSnippetIDParams

	// ------------- Optional query parameter "api_key" -------------

	err = runtime.BindQueryParameter("form", true, false, "api_key", c.Request.URL.Query(), &params.ApiKey)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter api_key: %s", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DeleteEnvsCodeSnippetID(c, codeSnippetID, params)
}

// PatchEnvsCodeSnippetID operation middleware
func (siw *ServerInterfaceWrapper) PatchEnvsCodeSnippetID(c *gin.Context) {

	var err error

	// ------------- Path parameter "codeSnippetID" -------------
	var codeSnippetID CodeSnippetID

	err = runtime.BindStyledParameter("simple", false, "codeSnippetID", c.Param("codeSnippetID"), &codeSnippetID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter codeSnippetID: %s", err), http.StatusBadRequest)
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params PatchEnvsCodeSnippetIDParams

	// ------------- Optional query parameter "api_key" -------------

	err = runtime.BindQueryParameter("form", true, false, "api_key", c.Request.URL.Query(), &params.ApiKey)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter api_key: %s", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PatchEnvsCodeSnippetID(c, codeSnippetID, params)
}

// PostEnvsCodeSnippetID operation middleware
func (siw *ServerInterfaceWrapper) PostEnvsCodeSnippetID(c *gin.Context) {

	var err error

	// ------------- Path parameter "codeSnippetID" -------------
	var codeSnippetID CodeSnippetID

	err = runtime.BindStyledParameter("simple", false, "codeSnippetID", c.Param("codeSnippetID"), &codeSnippetID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter codeSnippetID: %s", err), http.StatusBadRequest)
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params PostEnvsCodeSnippetIDParams

	// ------------- Optional query parameter "api_key" -------------

	err = runtime.BindQueryParameter("form", true, false, "api_key", c.Request.URL.Query(), &params.ApiKey)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter api_key: %s", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostEnvsCodeSnippetID(c, codeSnippetID, params)
}

// PutEnvsCodeSnippetIDState operation middleware
func (siw *ServerInterfaceWrapper) PutEnvsCodeSnippetIDState(c *gin.Context) {

	var err error

	// ------------- Path parameter "codeSnippetID" -------------
	var codeSnippetID CodeSnippetID

	err = runtime.BindStyledParameter("simple", false, "codeSnippetID", c.Param("codeSnippetID"), &codeSnippetID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter codeSnippetID: %s", err), http.StatusBadRequest)
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params PutEnvsCodeSnippetIDStateParams

	// ------------- Optional query parameter "api_key" -------------

	err = runtime.BindQueryParameter("form", true, false, "api_key", c.Request.URL.Query(), &params.ApiKey)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter api_key: %s", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PutEnvsCodeSnippetIDState(c, codeSnippetID, params)
}

// GetHealth operation middleware
func (siw *ServerInterfaceWrapper) GetHealth(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetHealth(c)
}

// GetSessions operation middleware
func (siw *ServerInterfaceWrapper) GetSessions(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetSessionsParams

	// ------------- Optional query parameter "api_key" -------------

	err = runtime.BindQueryParameter("form", true, false, "api_key", c.Request.URL.Query(), &params.ApiKey)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter api_key: %s", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetSessions(c, params)
}

// PostSessions operation middleware
func (siw *ServerInterfaceWrapper) PostSessions(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostSessionsParams

	// ------------- Optional query parameter "api_key" -------------

	err = runtime.BindQueryParameter("form", true, false, "api_key", c.Request.URL.Query(), &params.ApiKey)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter api_key: %s", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostSessions(c, params)
}

// DeleteSessionsSessionID operation middleware
func (siw *ServerInterfaceWrapper) DeleteSessionsSessionID(c *gin.Context) {

	var err error

	// ------------- Path parameter "sessionID" -------------
	var sessionID SessionID

	err = runtime.BindStyledParameter("simple", false, "sessionID", c.Param("sessionID"), &sessionID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sessionID: %s", err), http.StatusBadRequest)
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params DeleteSessionsSessionIDParams

	// ------------- Optional query parameter "api_key" -------------

	err = runtime.BindQueryParameter("form", true, false, "api_key", c.Request.URL.Query(), &params.ApiKey)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter api_key: %s", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DeleteSessionsSessionID(c, sessionID, params)
}

// PostSessionsSessionIDRefresh operation middleware
func (siw *ServerInterfaceWrapper) PostSessionsSessionIDRefresh(c *gin.Context) {

	var err error

	// ------------- Path parameter "sessionID" -------------
	var sessionID SessionID

	err = runtime.BindStyledParameter("simple", false, "sessionID", c.Param("sessionID"), &sessionID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sessionID: %s", err), http.StatusBadRequest)
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params PostSessionsSessionIDRefreshParams

	// ------------- Optional query parameter "api_key" -------------

	err = runtime.BindQueryParameter("form", true, false, "api_key", c.Request.URL.Query(), &params.ApiKey)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter api_key: %s", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostSessionsSessionIDRefresh(c, sessionID, params)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router *gin.Engine, si ServerInterface) *gin.Engine {
	return RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router *gin.Engine, si ServerInterface, options GinServerOptions) *gin.Engine {

	errorHandler := options.ErrorHandler

	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.DELETE(options.BaseURL+"/envs/:codeSnippetID", wrapper.DeleteEnvsCodeSnippetID)

	router.PATCH(options.BaseURL+"/envs/:codeSnippetID", wrapper.PatchEnvsCodeSnippetID)

	router.POST(options.BaseURL+"/envs/:codeSnippetID", wrapper.PostEnvsCodeSnippetID)

	router.PUT(options.BaseURL+"/envs/:codeSnippetID/state", wrapper.PutEnvsCodeSnippetIDState)

	router.GET(options.BaseURL+"/health", wrapper.GetHealth)

	router.GET(options.BaseURL+"/sessions", wrapper.GetSessions)

	router.POST(options.BaseURL+"/sessions", wrapper.PostSessions)

	router.DELETE(options.BaseURL+"/sessions/:sessionID", wrapper.DeleteSessionsSessionID)

	router.POST(options.BaseURL+"/sessions/:sessionID/refresh", wrapper.PostSessionsSessionIDRefresh)

	return router
}
