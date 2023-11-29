package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/posthog/posthog-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/e2b-dev/infra/packages/api/internal/api"
	"github.com/e2b-dev/infra/packages/api/internal/nomad"
	"github.com/e2b-dev/infra/packages/api/internal/utils"
	"github.com/e2b-dev/infra/packages/shared/pkg/telemetry"
)

func (a *APIStore) PostEnvs(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)

	// Prepare info for new env
	userID, team, tier, err := a.GetUserAndTeam(c)
	if err != nil {
		a.sendAPIStoreError(c, http.StatusInternalServerError, "Failed to get the default team")

		err = fmt.Errorf("error when getting default team: %w", err)
		telemetry.ReportCriticalError(ctx, err)

		return
	}

	envID := utils.GenerateID()

	telemetry.SetAttributes(ctx,
		attribute.String("env.user_id", userID.String()),
		attribute.String("env.team_id", team.ID.String()),
		attribute.String("team_name", team.Name),
		attribute.String("env.id", envID),
		attribute.String("tier", tier.ID),
	)

	buildID, err := uuid.NewRandom()
	if err != nil {
		err = fmt.Errorf("error when generating build id: %w", err)
		telemetry.ReportCriticalError(ctx, err)

		a.sendAPIStoreError(c, http.StatusInternalServerError, "Failed to generate build id")

		return
	}

	fileContent, fileHandler, err := c.Request.FormFile("buildContext")
	if err != nil {
		err = fmt.Errorf("error when parsing form data: %w", err)
		telemetry.ReportCriticalError(ctx, err)

		a.sendAPIStoreError(c, http.StatusBadRequest, "Failed to parse form data")

		return
	}

	defer func() {
		closeErr := fileContent.Close()
		if closeErr != nil {
			errMsg := fmt.Errorf("error when closing file: %w", closeErr)

			telemetry.ReportError(ctx, errMsg)
		}
	}()

	// Check if file is a tar.gz file
	if !strings.HasSuffix(fileHandler.Filename, ".tar.gz.e2b") {
		err = fmt.Errorf("build context doesn't have correct extension, the file is %s", fileHandler.Filename)
		telemetry.ReportCriticalError(ctx, err)

		a.sendAPIStoreError(c, http.StatusBadRequest, "Build context must be a tar.gz.e2b file")

		return
	}

	dockerfile := c.PostForm("dockerfile")
	alias := c.PostForm("alias")
	startCmd := c.PostForm("startCmd")

	if alias != "" {
		alias, err = utils.CleanEnvID(alias)
		if err != nil {
			a.sendAPIStoreError(c, http.StatusBadRequest, fmt.Sprintf("Invalid alias: %s", alias))

			err = fmt.Errorf("invalid alias: %w", err)
			telemetry.ReportCriticalError(ctx, err)

			return
		}
	}

	properties := a.GetPackageToPosthogProperties(&c.Request.Header)
	a.IdentifyAnalyticsTeam(team.ID.String(), team.Name)
	a.CreateAnalyticsUserEvent(userID.String(), team.ID.String(), "submitted environment build request", properties.
		Set("environment", envID).
		Set("build_id", buildID).
		Set("dockerfile", dockerfile).
		Set("alias", alias),
	)

	telemetry.SetAttributes(ctx,
		attribute.String("build.id", buildID.String()),
		attribute.String("build.alias", alias),
		attribute.String("build.dockerfile", dockerfile),
		attribute.String("build.startCmd", startCmd),
	)

	_, err = a.cloudStorage.streamFileUpload(strings.Join([]string{"v1", envID, buildID.String(), "context.tar.gz"}, "/"), fileContent)
	if err != nil {
		a.sendAPIStoreError(c, http.StatusInternalServerError, fmt.Sprintf("Error when uploading file to cloud storage: %s", err))

		err = fmt.Errorf("error when uploading file to cloud storage: %w", err)
		telemetry.ReportCriticalError(ctx, err)

		return
	}

	err = a.buildCache.Create(team.ID, envID, buildID)
	if err != nil {
		a.sendAPIStoreError(c, http.StatusConflict, fmt.Sprintf("there's already running build for %s", envID))

		err = fmt.Errorf("build is already running build for %s", envID)
		telemetry.ReportCriticalError(ctx, err)

		return
	}

	telemetry.ReportEvent(ctx, "started creating new environment")

	if alias != "" {
		err = a.supabase.EnsureEnvAlias(ctx, alias, envID)
		if err != nil {
			a.sendAPIStoreError(c, http.StatusInternalServerError, fmt.Sprintf("Error when inserting alias: %s", err))

			err = fmt.Errorf("error when inserting alias: %w", err)
			telemetry.ReportCriticalError(ctx, err)

			a.buildCache.Delete(envID, buildID)

			return
		} else {
			telemetry.ReportEvent(ctx, "inserted alias", attribute.String("alias", alias))
		}
	}

	go func() {
		buildContext, childSpan := a.tracer.Start(
			trace.ContextWithSpanContext(context.Background(), span.SpanContext()),
			"background-build-env",
		)

		var status api.EnvironmentBuildStatus

		buildErr := a.buildEnv(
			buildContext,
			userID.String(),
			team.ID,
			envID,
			buildID,
			dockerfile,
			startCmd,
			properties,
			nomad.BuildConfig{
				VCpuCount:  tier.Vcpu,
				MemoryMB:   tier.RAMMB,
				DiskSizeMB: tier.DiskMB,
			})

		if buildErr != nil {
			status = api.EnvironmentBuildStatusError

			errMsg := fmt.Errorf("error when building env: %w", buildErr)

			telemetry.ReportCriticalError(buildContext, errMsg)
		} else {
			status = api.EnvironmentBuildStatusReady

			telemetry.ReportEvent(buildContext, "created new environment", attribute.String("env_id", envID))
		}

		if status == api.EnvironmentBuildStatusError && alias != "" {
			errMsg := a.supabase.DeleteEnvAlias(buildContext, alias)
			if errMsg != nil {
				err = fmt.Errorf("error when deleting alias: %w", errMsg)
				telemetry.ReportError(buildContext, err)
			} else {
				telemetry.ReportEvent(buildContext, "deleted alias", attribute.String("alias", alias))
			}
		} else if status == api.EnvironmentBuildStatusReady && alias != "" {
			errMsg := a.supabase.UpdateEnvAlias(buildContext, alias, envID)
			if errMsg != nil {
				err = fmt.Errorf("error when updating alias: %w", errMsg)
				telemetry.ReportError(buildContext, err)
			} else {
				telemetry.ReportEvent(buildContext, "updated alias", attribute.String("alias", alias))
			}
		}

		cacheErr := a.buildCache.SetDone(envID, buildID, status)
		if cacheErr != nil {
			err = fmt.Errorf("error when setting build done in logs: %w", cacheErr)
			telemetry.ReportCriticalError(buildContext, cacheErr)
		}

		childSpan.End()
	}()

	aliases := []string{}

	if alias != "" {
		aliases = append(aliases, alias)
	}

	result := &api.Environment{
		EnvID:   envID,
		BuildID: buildID.String(),
		Public:  false,
		Aliases: &aliases,
	}

	c.JSON(http.StatusOK, result)
}

func (a *APIStore) PostEnvsEnvID(c *gin.Context, aliasOrEnvID api.EnvID) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)

	cleanedAliasOrEnvID, err := utils.CleanEnvID(aliasOrEnvID)
	if err != nil {
		a.sendAPIStoreError(c, http.StatusBadRequest, fmt.Sprintf("Invalid env ID: %s", aliasOrEnvID))

		err = fmt.Errorf("invalid env ID: %w", err)
		telemetry.ReportCriticalError(ctx, err)

		return
	}

	// Prepare info for rebuilding env
	userID, team, tier, err := a.GetUserAndTeam(c)
	if err != nil {
		a.sendAPIStoreError(c, http.StatusInternalServerError, fmt.Sprintf("Error when getting default team: %s", err))

		err = fmt.Errorf("error when getting default team: %w", err)
		telemetry.ReportCriticalError(ctx, err)

		return
	}

	telemetry.SetAttributes(ctx,
		attribute.String("env.user_id", userID.String()),
		attribute.String("env.team_id", team.ID.String()),
		attribute.String("tier", tier.ID),
	)

	buildID, err := uuid.NewRandom()
	if err != nil {
		err = fmt.Errorf("error when generating build id: %w", err)
		telemetry.ReportCriticalError(ctx, err)

		a.sendAPIStoreError(c, http.StatusInternalServerError, "Failed to generate build id")

		return
	}

	fileContent, fileHandler, err := c.Request.FormFile("buildContext")
	if err != nil {
		err = fmt.Errorf("error when parsing form data: %w", err)
		telemetry.ReportCriticalError(ctx, err)

		a.sendAPIStoreError(c, http.StatusBadRequest, "Failed to parse form data")

		return
	}

	// Check if file is a tar.gz file
	if !strings.HasSuffix(fileHandler.Filename, ".tar.gz.e2b") {
		err = fmt.Errorf("build context doesn't have correct extension, the file is %s", fileHandler.Filename)
		telemetry.ReportCriticalError(ctx, err)

		a.sendAPIStoreError(c, http.StatusBadRequest, "Build context must be a tar.gz.e2b file")

		return
	}

	dockerfile := c.PostForm("dockerfile")
	alias := c.PostForm("alias")
	if alias != "" {
		alias, err = utils.CleanEnvID(alias)
		if err != nil {
			a.sendAPIStoreError(c, http.StatusBadRequest, fmt.Sprintf("Invalid alias: %s", alias))

			err = fmt.Errorf("invalid alias: %w", err)
			telemetry.ReportCriticalError(ctx, err)

			return
		}
	}

	telemetry.SetAttributes(ctx,
		attribute.String("build.id", buildID.String()),
		attribute.String("build.alias", alias),
	)

	envID, hasAccess, accessErr := a.CheckTeamAccessEnv(ctx, cleanedAliasOrEnvID, team.ID, false)
	if accessErr != nil {
		a.sendAPIStoreError(c, http.StatusNotFound, fmt.Sprintf("the sandbox template '%s' does not exist", cleanedAliasOrEnvID))

		errMsg := fmt.Errorf("error env not found: %w", accessErr)
		telemetry.ReportError(ctx, errMsg)

		return
	}

	properties := a.GetPackageToPosthogProperties(&c.Request.Header)
	a.IdentifyAnalyticsTeam(team.ID.String(), team.Name)
	a.CreateAnalyticsUserEvent(userID.String(), team.ID.String(), "submitted environment build request", properties.
		Set("environment", envID).
		Set("build_id", buildID).
		Set("alias", alias).
		Set("dockerfile", dockerfile),
	)

	if !hasAccess {
		a.sendAPIStoreError(c, http.StatusForbidden, "You don't have access to this sandbox template")

		errMsg := fmt.Errorf("user doesn't have access to env '%s'", envID)
		telemetry.ReportError(ctx, errMsg)

		return
	}

	telemetry.SetAttributes(ctx, attribute.String("build.id", buildID.String()))

	_, err = a.cloudStorage.streamFileUpload(strings.Join([]string{"v1", envID, buildID.String(), "context.tar.gz"}, "/"), fileContent)
	if err != nil {
		a.sendAPIStoreError(c, http.StatusInternalServerError, fmt.Sprintf("Error when uploading file to cloud storage: %s", err))

		err = fmt.Errorf("error when uploading file to cloud storage: %w", err)
		telemetry.ReportCriticalError(ctx, err)

		return
	}

	err = a.buildCache.Create(team.ID, envID, buildID)
	if err != nil {
		a.sendAPIStoreError(c, http.StatusConflict, fmt.Sprintf("there's already running build for %s", envID))

		err = fmt.Errorf("build is already running build for %s", envID)
		telemetry.ReportCriticalError(ctx, err)

		return
	}

	telemetry.ReportEvent(ctx, "started updating environment")

	if alias != "" {
		err = a.supabase.EnsureEnvAlias(ctx, alias, envID)
		if err != nil {
			a.sendAPIStoreError(c, http.StatusInternalServerError, fmt.Sprintf("Error when inserting alias: %s", err))

			err = fmt.Errorf("error when inserting alias: %w", err)
			telemetry.ReportCriticalError(ctx, err)

			a.buildCache.Delete(envID, buildID)

			return
		} else {
			telemetry.ReportEvent(ctx, "inserted alias", attribute.String("alias", alias))
		}
	}

	startCmd := c.PostForm("startCmd")

	go func() {
		buildContext, childSpan := a.tracer.Start(
			trace.ContextWithSpanContext(context.Background(), span.SpanContext()),
			"background-build-env",
		)

		var status api.EnvironmentBuildStatus

		buildErr := a.buildEnv(
			buildContext,
			userID.String(),
			team.ID,
			envID,
			buildID,
			dockerfile,
			startCmd,
			properties,
			nomad.BuildConfig{
				VCpuCount:  tier.Vcpu,
				MemoryMB:   tier.RAMMB,
				DiskSizeMB: tier.DiskMB,
			})

		if buildErr != nil {
			status = api.EnvironmentBuildStatusError

			errMsg := fmt.Errorf("error when building env: %w", buildErr)

			telemetry.ReportCriticalError(buildContext, errMsg)
		} else {
			status = api.EnvironmentBuildStatusReady

			telemetry.ReportEvent(buildContext, "created new environment", attribute.String("env_id", envID))
		}

		if status == api.EnvironmentBuildStatusError && alias != "" {
			errMsg := a.supabase.DeleteNilEnvAlias(buildContext, alias)
			if errMsg != nil {
				err = fmt.Errorf("error when deleting alias: %w", errMsg)
				telemetry.ReportError(buildContext, err)
			} else {
				telemetry.ReportEvent(buildContext, "deleted alias", attribute.String("alias", alias))
			}
		} else if status == api.EnvironmentBuildStatusReady && alias != "" {
			errMsg := a.supabase.UpdateEnvAlias(buildContext, alias, envID)
			if errMsg != nil {
				err = fmt.Errorf("error when updating alias: %w", errMsg)
				telemetry.ReportError(buildContext, err)
			} else {
				telemetry.ReportEvent(buildContext, "updated alias", attribute.String("alias", alias))
			}
		}

		cacheErr := a.buildCache.SetDone(envID, buildID, status)
		if cacheErr != nil {
			errMsg := fmt.Errorf("error when setting build done in logs: %w", cacheErr)
			telemetry.ReportCriticalError(buildContext, errMsg)
		}

		childSpan.End()
	}()

	aliases := []string{}

	if alias != "" {
		aliases = append(aliases, alias)
	}

	result := &api.Environment{
		EnvID:   envID,
		BuildID: buildID.String(),
		Public:  false,
		Aliases: &aliases,
	}

	c.JSON(http.StatusOK, result)
}

func (a *APIStore) buildEnv(
	ctx context.Context,
	userID string,
	teamID uuid.UUID,
	envID string,
	buildID uuid.UUID,
	dockerfile,
	startCmd string,
	posthogProperties posthog.Properties,
	vmConfig nomad.BuildConfig,
) (err error) {
	childCtx, childSpan := a.tracer.Start(ctx, "build-env",
		trace.WithAttributes(
			attribute.String("env_id", envID),
			attribute.String("build_id", buildID.String()),
		),
	)
	defer childSpan.End()

	startTime := time.Now()

	defer func() {
		a.CreateAnalyticsUserEvent(userID, teamID.String(), "built environment", posthogProperties.
			Set("environment", envID).
			Set("build_id", buildID).
			Set("duration", time.Since(startTime).String()).
			Set("success", err != nil),
		)
	}()

	err = a.nomad.BuildEnvJob(
		a.tracer,
		childCtx,
		envID,
		buildID.String(),
		startCmd,
		a.apiSecret,
		a.googleServiceAccountBase64,
		vmConfig,
	)
	if err != nil {
		err = fmt.Errorf("error when building env: %w", err)
		telemetry.ReportCriticalError(childCtx, err)

		return err
	}

	err = a.supabase.UpsertEnv(ctx, teamID, envID, buildID, dockerfile)

	if err != nil {
		err = fmt.Errorf("error when updating env: %w", err)
		telemetry.ReportCriticalError(childCtx, err)

		return err
	}

	return nil
}
