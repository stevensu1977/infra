package server

import (
	"context"
	"fmt"
	"os"
	"time"

	artifactregistry "cloud.google.com/go/artifactregistry/apiv1"
	"github.com/docker/docker/client"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"

	"github.com/e2b-dev/infra/packages/shared/pkg/consts"
	e2bgrpc "github.com/e2b-dev/infra/packages/shared/pkg/grpc"
	templatemanager "github.com/e2b-dev/infra/packages/shared/pkg/grpc/template-manager"
	l "github.com/e2b-dev/infra/packages/shared/pkg/logger"
	"github.com/e2b-dev/infra/packages/template-manager/internal/constants"
	"github.com/e2b-dev/infra/packages/template-manager/internal/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

type serverStore struct {
	templatemanager.UnimplementedTemplateServiceServer
	server             *grpc.Server
	tracer             trace.Tracer
	logger             *zap.Logger
	buildLogger        *zap.Logger
	dockerClient       *client.Client
	legacyDockerClient *docker.Client
	artifactRegistry   *artifactregistry.Client
	templateStorage    *template.Storage
	ecrClient          *ecr.Client
}

func New(logger *zap.Logger, buildLogger *zap.Logger) *grpc.Server {
	ctx := context.Background()
	logger.Info("Initializing template manager")

	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.PayloadReceived, logging.PayloadSent, logging.FinishCall),
		logging.WithLevels(logging.DefaultServerCodeToLevel),
		logging.WithFieldsFromContext(logging.ExtractFields),
	}

	s := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second, // Minimum time between pings from client
			PermitWithoutStream: true,            // Allow pings even when no active streams
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    15 * time.Second, // Server sends keepalive pings every 15s
			Timeout: 5 * time.Second,  // Wait 5s for response before considering dead
		}),
		grpc.StatsHandler(e2bgrpc.NewStatsWrapper(otelgrpc.NewServerHandler())),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
			selector.UnaryServerInterceptor(
				logging.UnaryServerInterceptor(l.GRPCLogger(logger), opts...),
				l.WithoutHealthCheck(),
			),
		),
	)
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	legacyClient, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	var ecrClient *ecr.Client
	var artifactRegistry *artifactregistry.Client

	if consts.CloudProviderEnv == consts.AWS {
		region := os.Getenv("AWS_REGION")
		if region == "" {
			region = "us-east-1"
		}
		cfg, err := loadAWSConfig(ctx, region)
		if err != nil {
			panic(err)
		}

		ecrClient = ecr.NewFromConfig(cfg)

	}
	if consts.CloudProviderEnv == consts.GCP {
		artifactRegistry, err = artifactregistry.NewClient(ctx)
		if err != nil {
			panic(err)
		}
	}

	if consts.CloudProviderEnv != consts.AWS && consts.CloudProviderEnv != consts.GCP {
		panic(fmt.Errorf("unsupported cloud provider: %s", consts.CloudProviderEnv))
	}

	templateStorage := template.NewStorage(ctx)

	templatemanager.RegisterTemplateServiceServer(s, &serverStore{
		tracer:             otel.Tracer(constants.ServiceName),
		logger:             logger,
		buildLogger:        buildLogger,
		dockerClient:       dockerClient,
		legacyDockerClient: legacyClient,
		artifactRegistry:   artifactRegistry,
		templateStorage:    templateStorage,
		ecrClient:          ecrClient,
	})

	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	return s
}

func loadAWSConfig(ctx context.Context, region string) (aws.Config, error) {
	configOpts := []func(*config.LoadOptions) error{}

	if region != "" {
		configOpts = append(configOpts, config.WithRegion(region))
	}

	return config.LoadDefaultConfig(ctx, configOpts...)
}
