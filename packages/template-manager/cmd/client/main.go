package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/e2b-dev/infra/packages/shared/pkg/grpc/template-manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up connection to the server
	conn, err := grpc.Dial("localhost:5009", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create client
	client := pb.NewTemplateServiceClient(conn)

	// Example: Create a template
	template := &pb.TemplateConfig{
		TemplateID:         "example-template",
		BuildID:            "build-123",
		MemoryMB:           1024,
		VCpuCount:          2,
		DiskSizeMB:         10240,
		KernelVersion:      "5.10",
		FirecrackerVersion: "1.0",
		StartCommand:       "./start.sh",
		HugePages:          false,
	}

	req := &pb.TemplateCreateRequest{
		Template: template,
	}

	// Call TemplateCreate and stream the logs
	stream, err := client.TemplateCreate(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling TemplateCreate: %v", err)
	}

	for {
		buildLog, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving build log: %v", err)
		}
		fmt.Printf("Build log: %s\n", buildLog.Log)
	}

	// Example: Delete a template build
	deleteReq := &pb.TemplateBuildDeleteRequest{
		BuildID: "build-123",
	}

	_, err = client.TemplateBuildDelete(context.Background(), deleteReq)
	if err != nil {
		log.Fatalf("Error deleting template build: %v", err)
	}
	fmt.Println("Template build deleted successfully")
}
