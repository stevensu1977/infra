// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /health)
	GetHealth(c *gin.Context)

	// (GET /sandboxes)
	GetSandboxes(c *gin.Context)

	// (POST /sandboxes)
	PostSandboxes(c *gin.Context)

	// (DELETE /sandboxes/{sandboxID})
	DeleteSandboxesSandboxID(c *gin.Context, sandboxID SandboxID)

	// (GET /sandboxes/{sandboxID}/logs)
	GetSandboxesSandboxIDLogs(c *gin.Context, sandboxID SandboxID, params GetSandboxesSandboxIDLogsParams)

	// (POST /sandboxes/{sandboxID}/pause)
	PostSandboxesSandboxIDPause(c *gin.Context, sandboxID SandboxID)

	// (POST /sandboxes/{sandboxID}/refreshes)
	PostSandboxesSandboxIDRefreshes(c *gin.Context, sandboxID SandboxID)

	// (POST /sandboxes/{sandboxID}/resume)
	PostSandboxesSandboxIDResume(c *gin.Context, sandboxID SandboxID)

	// (POST /sandboxes/{sandboxID}/timeout)
	PostSandboxesSandboxIDTimeout(c *gin.Context, sandboxID SandboxID)

	// (GET /teams)
	GetTeams(c *gin.Context)

	// (GET /templates)
	GetTemplates(c *gin.Context, params GetTemplatesParams)

	// (POST /templates)
	PostTemplates(c *gin.Context)

	// (DELETE /templates/{templateID})
	DeleteTemplatesTemplateID(c *gin.Context, templateID TemplateID)

	// (POST /templates/{templateID})
	PostTemplatesTemplateID(c *gin.Context, templateID TemplateID)

	// (POST /templates/{templateID}/builds/{buildID})
	PostTemplatesTemplateIDBuildsBuildID(c *gin.Context, templateID TemplateID, buildID BuildID)

	// (GET /templates/{templateID}/builds/{buildID}/status)
	GetTemplatesTemplateIDBuildsBuildIDStatus(c *gin.Context, templateID TemplateID, buildID BuildID, params GetTemplatesTemplateIDBuildsBuildIDStatusParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetHealth operation middleware
func (siw *ServerInterfaceWrapper) GetHealth(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetHealth(c)
}

// GetSandboxes operation middleware
func (siw *ServerInterfaceWrapper) GetSandboxes(c *gin.Context) {

	c.Set(ApiKeyAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetSandboxes(c)
}

// PostSandboxes operation middleware
func (siw *ServerInterfaceWrapper) PostSandboxes(c *gin.Context) {

	c.Set(ApiKeyAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostSandboxes(c)
}

// DeleteSandboxesSandboxID operation middleware
func (siw *ServerInterfaceWrapper) DeleteSandboxesSandboxID(c *gin.Context) {

	var err error

	// ------------- Path parameter "sandboxID" -------------
	var sandboxID SandboxID

	err = runtime.BindStyledParameter("simple", false, "sandboxID", c.Param("sandboxID"), &sandboxID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sandboxID: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(ApiKeyAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeleteSandboxesSandboxID(c, sandboxID)
}

// GetSandboxesSandboxIDLogs operation middleware
func (siw *ServerInterfaceWrapper) GetSandboxesSandboxIDLogs(c *gin.Context) {

	var err error

	// ------------- Path parameter "sandboxID" -------------
	var sandboxID SandboxID

	err = runtime.BindStyledParameter("simple", false, "sandboxID", c.Param("sandboxID"), &sandboxID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sandboxID: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(ApiKeyAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetSandboxesSandboxIDLogsParams

	// ------------- Optional query parameter "start" -------------

	err = runtime.BindQueryParameter("form", true, false, "start", c.Request.URL.Query(), &params.Start)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter start: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", c.Request.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter limit: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetSandboxesSandboxIDLogs(c, sandboxID, params)
}

// PostSandboxesSandboxIDPause operation middleware
func (siw *ServerInterfaceWrapper) PostSandboxesSandboxIDPause(c *gin.Context) {

	var err error

	// ------------- Path parameter "sandboxID" -------------
	var sandboxID SandboxID

	err = runtime.BindStyledParameter("simple", false, "sandboxID", c.Param("sandboxID"), &sandboxID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sandboxID: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(ApiKeyAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostSandboxesSandboxIDPause(c, sandboxID)
}

// PostSandboxesSandboxIDRefreshes operation middleware
func (siw *ServerInterfaceWrapper) PostSandboxesSandboxIDRefreshes(c *gin.Context) {

	var err error

	// ------------- Path parameter "sandboxID" -------------
	var sandboxID SandboxID

	err = runtime.BindStyledParameter("simple", false, "sandboxID", c.Param("sandboxID"), &sandboxID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sandboxID: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(ApiKeyAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostSandboxesSandboxIDRefreshes(c, sandboxID)
}

// PostSandboxesSandboxIDResume operation middleware
func (siw *ServerInterfaceWrapper) PostSandboxesSandboxIDResume(c *gin.Context) {

	var err error

	// ------------- Path parameter "sandboxID" -------------
	var sandboxID SandboxID

	err = runtime.BindStyledParameter("simple", false, "sandboxID", c.Param("sandboxID"), &sandboxID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sandboxID: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(ApiKeyAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostSandboxesSandboxIDResume(c, sandboxID)
}

// PostSandboxesSandboxIDTimeout operation middleware
func (siw *ServerInterfaceWrapper) PostSandboxesSandboxIDTimeout(c *gin.Context) {

	var err error

	// ------------- Path parameter "sandboxID" -------------
	var sandboxID SandboxID

	err = runtime.BindStyledParameter("simple", false, "sandboxID", c.Param("sandboxID"), &sandboxID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sandboxID: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(ApiKeyAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostSandboxesSandboxIDTimeout(c, sandboxID)
}

// GetTeams operation middleware
func (siw *ServerInterfaceWrapper) GetTeams(c *gin.Context) {

	c.Set(AccessTokenAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetTeams(c)
}

// GetTemplates operation middleware
func (siw *ServerInterfaceWrapper) GetTemplates(c *gin.Context) {

	var err error

	c.Set(AccessTokenAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetTemplatesParams

	// ------------- Optional query parameter "teamID" -------------

	err = runtime.BindQueryParameter("form", true, false, "teamID", c.Request.URL.Query(), &params.TeamID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter teamID: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetTemplates(c, params)
}

// PostTemplates operation middleware
func (siw *ServerInterfaceWrapper) PostTemplates(c *gin.Context) {

	c.Set(AccessTokenAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostTemplates(c)
}

// DeleteTemplatesTemplateID operation middleware
func (siw *ServerInterfaceWrapper) DeleteTemplatesTemplateID(c *gin.Context) {

	var err error

	// ------------- Path parameter "templateID" -------------
	var templateID TemplateID

	err = runtime.BindStyledParameter("simple", false, "templateID", c.Param("templateID"), &templateID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter templateID: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(AccessTokenAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeleteTemplatesTemplateID(c, templateID)
}

// PostTemplatesTemplateID operation middleware
func (siw *ServerInterfaceWrapper) PostTemplatesTemplateID(c *gin.Context) {

	var err error

	// ------------- Path parameter "templateID" -------------
	var templateID TemplateID

	err = runtime.BindStyledParameter("simple", false, "templateID", c.Param("templateID"), &templateID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter templateID: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(AccessTokenAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostTemplatesTemplateID(c, templateID)
}

// PostTemplatesTemplateIDBuildsBuildID operation middleware
func (siw *ServerInterfaceWrapper) PostTemplatesTemplateIDBuildsBuildID(c *gin.Context) {

	var err error

	// ------------- Path parameter "templateID" -------------
	var templateID TemplateID

	err = runtime.BindStyledParameter("simple", false, "templateID", c.Param("templateID"), &templateID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter templateID: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "buildID" -------------
	var buildID BuildID

	err = runtime.BindStyledParameter("simple", false, "buildID", c.Param("buildID"), &buildID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter buildID: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(AccessTokenAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostTemplatesTemplateIDBuildsBuildID(c, templateID, buildID)
}

// GetTemplatesTemplateIDBuildsBuildIDStatus operation middleware
func (siw *ServerInterfaceWrapper) GetTemplatesTemplateIDBuildsBuildIDStatus(c *gin.Context) {

	var err error

	// ------------- Path parameter "templateID" -------------
	var templateID TemplateID

	err = runtime.BindStyledParameter("simple", false, "templateID", c.Param("templateID"), &templateID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter templateID: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "buildID" -------------
	var buildID BuildID

	err = runtime.BindStyledParameter("simple", false, "buildID", c.Param("buildID"), &buildID)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter buildID: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(AccessTokenAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetTemplatesTemplateIDBuildsBuildIDStatusParams

	// ------------- Optional query parameter "logsOffset" -------------

	err = runtime.BindQueryParameter("form", true, false, "logsOffset", c.Request.URL.Query(), &params.LogsOffset)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter logsOffset: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetTemplatesTemplateIDBuildsBuildIDStatus(c, templateID, buildID, params)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
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

	router.GET(options.BaseURL+"/health", wrapper.GetHealth)
	router.GET(options.BaseURL+"/sandboxes", wrapper.GetSandboxes)
	router.POST(options.BaseURL+"/sandboxes", wrapper.PostSandboxes)
	router.DELETE(options.BaseURL+"/sandboxes/:sandboxID", wrapper.DeleteSandboxesSandboxID)
	router.GET(options.BaseURL+"/sandboxes/:sandboxID/logs", wrapper.GetSandboxesSandboxIDLogs)
	router.POST(options.BaseURL+"/sandboxes/:sandboxID/pause", wrapper.PostSandboxesSandboxIDPause)
	router.POST(options.BaseURL+"/sandboxes/:sandboxID/refreshes", wrapper.PostSandboxesSandboxIDRefreshes)
	router.POST(options.BaseURL+"/sandboxes/:sandboxID/resume", wrapper.PostSandboxesSandboxIDResume)
	router.POST(options.BaseURL+"/sandboxes/:sandboxID/timeout", wrapper.PostSandboxesSandboxIDTimeout)
	router.GET(options.BaseURL+"/teams", wrapper.GetTeams)
	router.GET(options.BaseURL+"/templates", wrapper.GetTemplates)
	router.POST(options.BaseURL+"/templates", wrapper.PostTemplates)
	router.DELETE(options.BaseURL+"/templates/:templateID", wrapper.DeleteTemplatesTemplateID)
	router.POST(options.BaseURL+"/templates/:templateID", wrapper.PostTemplatesTemplateID)
	router.POST(options.BaseURL+"/templates/:templateID/builds/:buildID", wrapper.PostTemplatesTemplateIDBuildsBuildID)
	router.GET(options.BaseURL+"/templates/:templateID/builds/:buildID/status", wrapper.GetTemplatesTemplateIDBuildsBuildIDStatus)
}
