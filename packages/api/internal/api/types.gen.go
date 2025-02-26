// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	AccessTokenAuthScopes = "AccessTokenAuth.Scopes"
	AdminTokenAuthScopes  = "AdminTokenAuth.Scopes"
	ApiKeyAuthScopes      = "ApiKeyAuth.Scopes"
)

// Defines values for NodeStatus.
const (
	NodeStatusDraining NodeStatus = "draining"
	NodeStatusReady    NodeStatus = "ready"
)

// Defines values for TemplateBuildStatus.
const (
	TemplateBuildStatusBuilding TemplateBuildStatus = "building"
	TemplateBuildStatusError    TemplateBuildStatus = "error"
	TemplateBuildStatusReady    TemplateBuildStatus = "ready"
)

// CPUCount CPU cores for the sandbox
type CPUCount = int32

// EnvVars defines model for EnvVars.
type EnvVars map[string]string

// Error defines model for Error.
type Error struct {
	// Code Error code
	Code int32 `json:"code"`

	// Message Error
	Message string `json:"message"`
}

// MemoryMB Memory for the sandbox in MB
type MemoryMB = int32

// NewSandbox defines model for NewSandbox.
type NewSandbox struct {
	// AutoPause Automatically pauses the sandbox after the timeout
	AutoPause *bool            `json:"autoPause,omitempty"`
	EnvVars   *EnvVars         `json:"envVars,omitempty"`
	Metadata  *SandboxMetadata `json:"metadata,omitempty"`

	// TemplateID Identifier of the required template
	TemplateID string `json:"templateID"`

	// Timeout Time to live for the sandbox in seconds.
	Timeout *int32 `json:"timeout,omitempty"`
}

// Node defines model for Node.
type Node struct {
	// AllocatedCPU Number of allocated CPU cores
	AllocatedCPU int32 `json:"allocatedCPU"`

	// AllocatedMemoryMiB Amount of allocated memory in MiB
	AllocatedMemoryMiB int32 `json:"allocatedMemoryMiB"`

	// CreateFails Number of sandbox create fails
	CreateFails uint64 `json:"createFails"`

	// NodeID Identifier of the node
	NodeID string `json:"nodeID"`

	// SandboxCount Number of sandboxes running on the node
	SandboxCount int32 `json:"sandboxCount"`

	// Status Status of the node
	Status NodeStatus `json:"status"`
}

// NodeDetail defines model for NodeDetail.
type NodeDetail struct {
	// CachedBuilds List of cached builds id on the node
	CachedBuilds []string `json:"cachedBuilds"`

	// CreateFails Number of sandbox create fails
	CreateFails uint64 `json:"createFails"`

	// NodeID Identifier of the node
	NodeID string `json:"nodeID"`

	// Sandboxes List of sandboxes running on the node
	Sandboxes []RunningSandbox `json:"sandboxes"`

	// Status Status of the node
	Status NodeStatus `json:"status"`
}

// NodeStatus Status of the node
type NodeStatus string

// NodeStatusChange defines model for NodeStatusChange.
type NodeStatusChange struct {
	// Status Status of the node
	Status NodeStatus `json:"status"`
}

// ResumedSandbox defines model for ResumedSandbox.
type ResumedSandbox struct {
	// AutoPause Automatically pauses the sandbox after the timeout
	AutoPause *bool `json:"autoPause,omitempty"`

	// Timeout Time to live for the sandbox in seconds.
	Timeout *int32 `json:"timeout,omitempty"`
}

// RunningSandbox defines model for RunningSandbox.
type RunningSandbox struct {
	// Alias Alias of the template
	Alias *string `json:"alias,omitempty"`

	// ClientID Identifier of the client
	ClientID string `json:"clientID"`

	// CpuCount CPU cores for the sandbox
	CpuCount CPUCount `json:"cpuCount"`

	// EndAt Time when the sandbox will expire
	EndAt time.Time `json:"endAt"`

	// MemoryMB Memory for the sandbox in MB
	MemoryMB MemoryMB         `json:"memoryMB"`
	Metadata *SandboxMetadata `json:"metadata,omitempty"`

	// SandboxID Identifier of the sandbox
	SandboxID string `json:"sandboxID"`

	// StartedAt Time when the sandbox was started
	StartedAt time.Time `json:"startedAt"`

	// TemplateID Identifier of the template from which is the sandbox created
	TemplateID string `json:"templateID"`
}

// RunningSandboxWithMetrics defines model for RunningSandboxWithMetrics.
type RunningSandboxWithMetrics struct {
	// Alias Alias of the template
	Alias *string `json:"alias,omitempty"`

	// ClientID Identifier of the client
	ClientID string `json:"clientID"`

	// CpuCount CPU cores for the sandbox
	CpuCount CPUCount `json:"cpuCount"`

	// EndAt Time when the sandbox will expire
	EndAt time.Time `json:"endAt"`

	// MemoryMB Memory for the sandbox in MB
	MemoryMB MemoryMB         `json:"memoryMB"`
	Metadata *SandboxMetadata `json:"metadata,omitempty"`
	Metrics  *[]SandboxMetric `json:"metrics,omitempty"`

	// SandboxID Identifier of the sandbox
	SandboxID string `json:"sandboxID"`

	// StartedAt Time when the sandbox was started
	StartedAt time.Time `json:"startedAt"`

	// TemplateID Identifier of the template from which is the sandbox created
	TemplateID string `json:"templateID"`
}

// Sandbox defines model for Sandbox.
type Sandbox struct {
	// Alias Alias of the template
	Alias *string `json:"alias,omitempty"`

	// ClientID Identifier of the client
	ClientID string `json:"clientID"`

	// EnvdVersion Version of the envd running in the sandbox
	EnvdVersion string `json:"envdVersion"`

	// SandboxID Identifier of the sandbox
	SandboxID string `json:"sandboxID"`

	// TemplateID Identifier of the template from which is the sandbox created
	TemplateID string `json:"templateID"`
}

// SandboxLog Log entry with timestamp and line
type SandboxLog struct {
	// Line Log line content
	Line string `json:"line"`

	// Timestamp Timestamp of the log entry
	Timestamp time.Time `json:"timestamp"`
}

// SandboxLogs defines model for SandboxLogs.
type SandboxLogs struct {
	// Logs Logs of the sandbox
	Logs []SandboxLog `json:"logs"`
}

// SandboxMetadata defines model for SandboxMetadata.
type SandboxMetadata map[string]string

// SandboxMetric Metric entry with timestamp and line
type SandboxMetric struct {
	// CpuCount Number of CPU cores
	CpuCount int32 `json:"cpuCount"`

	// CpuUsedPct CPU usage percentage
	CpuUsedPct float32 `json:"cpuUsedPct"`

	// MemTotalMiB Total memory in MiB
	MemTotalMiB int64 `json:"memTotalMiB"`

	// MemUsedMiB Memory used in MiB
	MemUsedMiB int64 `json:"memUsedMiB"`

	// Timestamp Timestamp of the metric entry
	Timestamp time.Time `json:"timestamp"`
}

// Team defines model for Team.
type Team struct {
	// ApiKey API key for the team
	ApiKey string `json:"apiKey"`

	// IsDefault Whether the team is the default team
	IsDefault bool `json:"isDefault"`

	// Name Name of the team
	Name string `json:"name"`

	// TeamID Identifier of the team
	TeamID string `json:"teamID"`
}

// TeamUser defines model for TeamUser.
type TeamUser struct {
	// Email Email of the user
	Email string `json:"email"`

	// Id Identifier of the user
	Id openapi_types.UUID `json:"id"`
}

// Template defines model for Template.
type Template struct {
	// Aliases Aliases of the template
	Aliases *[]string `json:"aliases,omitempty"`

	// BuildCount Number of times the template was built
	BuildCount int32 `json:"buildCount"`

	// BuildID Identifier of the last successful build for given template
	BuildID string `json:"buildID"`

	// CpuCount CPU cores for the sandbox
	CpuCount CPUCount `json:"cpuCount"`

	// CreatedAt Time when the template was created
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy *TeamUser `json:"createdBy"`

	// LastSpawnedAt Time when the template was last used
	LastSpawnedAt time.Time `json:"lastSpawnedAt"`

	// MemoryMB Memory for the sandbox in MB
	MemoryMB MemoryMB `json:"memoryMB"`

	// Public Whether the template is public or only accessible by the team
	Public bool `json:"public"`

	// SpawnCount Number of times the template was used
	SpawnCount int64 `json:"spawnCount"`

	// TemplateID Identifier of the template
	TemplateID string `json:"templateID"`

	// UpdatedAt Time when the template was last updated
	UpdatedAt time.Time `json:"updatedAt"`
}

// TemplateBuild defines model for TemplateBuild.
type TemplateBuild struct {
	// BuildID Identifier of the build
	BuildID string `json:"buildID"`

	// Logs Build logs
	Logs []string `json:"logs"`

	// Status Status of the template
	Status TemplateBuildStatus `json:"status"`

	// TemplateID Identifier of the template
	TemplateID string `json:"templateID"`
}

// TemplateBuildStatus Status of the template
type TemplateBuildStatus string

// TemplateBuildRequest defines model for TemplateBuildRequest.
type TemplateBuildRequest struct {
	// Alias Alias of the template
	Alias *string `json:"alias,omitempty"`

	// CpuCount CPU cores for the sandbox
	CpuCount *CPUCount `json:"cpuCount,omitempty"`

	// Dockerfile Dockerfile for the template
	Dockerfile string `json:"dockerfile"`

	// MemoryMB Memory for the sandbox in MB
	MemoryMB *MemoryMB `json:"memoryMB,omitempty"`

	// StartCmd Start command to execute in the template after the build
	StartCmd *string `json:"startCmd,omitempty"`

	// TeamID Identifier of the team
	TeamID *string `json:"teamID,omitempty"`
}

// TemplateUpdateRequest defines model for TemplateUpdateRequest.
type TemplateUpdateRequest struct {
	// Public Whether the template is public or only accessible by the team
	Public *bool `json:"public,omitempty"`
}

// BuildID defines model for buildID.
type BuildID = string

// NodeID defines model for nodeID.
type NodeID = string

// SandboxID defines model for sandboxID.
type SandboxID = string

// TemplateID defines model for templateID.
type TemplateID = string

// N400 defines model for 400.
type N400 = Error

// N401 defines model for 401.
type N401 = Error

// N404 defines model for 404.
type N404 = Error

// N409 defines model for 409.
type N409 = Error

// N500 defines model for 500.
type N500 = Error

// GetSandboxesParams defines parameters for GetSandboxes.
type GetSandboxesParams struct {
	// Query A query used to filter the sandboxes (e.g. "user=abc&app=prod"). Query and each key and values must be URL encoded.
	Query *string `form:"query,omitempty" json:"query,omitempty"`
}

// GetSandboxesMetricsParams defines parameters for GetSandboxesMetrics.
type GetSandboxesMetricsParams struct {
	// Query A query used to filter the sandboxes (e.g. "user=abc&app=prod"). Query and each key and values must be URL encoded.
	Query *string `form:"query,omitempty" json:"query,omitempty"`
}

// GetSandboxesSandboxIDLogsParams defines parameters for GetSandboxesSandboxIDLogs.
type GetSandboxesSandboxIDLogsParams struct {
	// Start Starting timestamp of the logs that should be returned in milliseconds
	Start *int64 `form:"start,omitempty" json:"start,omitempty"`

	// Limit Maximum number of logs that should be returned
	Limit *int32 `form:"limit,omitempty" json:"limit,omitempty"`
}

// PostSandboxesSandboxIDRefreshesJSONBody defines parameters for PostSandboxesSandboxIDRefreshes.
type PostSandboxesSandboxIDRefreshesJSONBody struct {
	// Duration Duration for which the sandbox should be kept alive in seconds
	Duration *int `json:"duration,omitempty"`
}

// PostSandboxesSandboxIDTimeoutJSONBody defines parameters for PostSandboxesSandboxIDTimeout.
type PostSandboxesSandboxIDTimeoutJSONBody struct {
	// Timeout Timeout in seconds from the current time after which the sandbox should expire
	Timeout int32 `json:"timeout"`
}

// GetTemplatesParams defines parameters for GetTemplates.
type GetTemplatesParams struct {
	TeamID *string `form:"teamID,omitempty" json:"teamID,omitempty"`
}

// GetTemplatesTemplateIDBuildsBuildIDStatusParams defines parameters for GetTemplatesTemplateIDBuildsBuildIDStatus.
type GetTemplatesTemplateIDBuildsBuildIDStatusParams struct {
	// LogsOffset Index of the starting build log that should be returned with the template
	LogsOffset *int32 `form:"logsOffset,omitempty" json:"logsOffset,omitempty"`
}

// PostNodesNodeIDJSONRequestBody defines body for PostNodesNodeID for application/json ContentType.
type PostNodesNodeIDJSONRequestBody = NodeStatusChange

// PostSandboxesJSONRequestBody defines body for PostSandboxes for application/json ContentType.
type PostSandboxesJSONRequestBody = NewSandbox

// PostSandboxesSandboxIDRefreshesJSONRequestBody defines body for PostSandboxesSandboxIDRefreshes for application/json ContentType.
type PostSandboxesSandboxIDRefreshesJSONRequestBody PostSandboxesSandboxIDRefreshesJSONBody

// PostSandboxesSandboxIDResumeJSONRequestBody defines body for PostSandboxesSandboxIDResume for application/json ContentType.
type PostSandboxesSandboxIDResumeJSONRequestBody = ResumedSandbox

// PostSandboxesSandboxIDTimeoutJSONRequestBody defines body for PostSandboxesSandboxIDTimeout for application/json ContentType.
type PostSandboxesSandboxIDTimeoutJSONRequestBody PostSandboxesSandboxIDTimeoutJSONBody

// PostTemplatesJSONRequestBody defines body for PostTemplates for application/json ContentType.
type PostTemplatesJSONRequestBody = TemplateBuildRequest

// PatchTemplatesTemplateIDJSONRequestBody defines body for PatchTemplatesTemplateID for application/json ContentType.
type PatchTemplatesTemplateIDJSONRequestBody = TemplateUpdateRequest

// PostTemplatesTemplateIDJSONRequestBody defines body for PostTemplatesTemplateID for application/json ContentType.
type PostTemplatesTemplateIDJSONRequestBody = TemplateBuildRequest
