package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	pb "github.com/e2b-dev/infra/packages/shared/pkg/grpc/template-manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Define command line flags
	host := flag.String("host", "localhost", "Host")
	port := flag.Int("port", 5009, "Server port")
	templateID := flag.String("template-id", "example-template", "Template ID")
	buildID := flag.String("build-id", "build-123", "Build ID")
	memoryMB := flag.Int("memory", 1024, "Memory in MB")
	vcpuCount := flag.Int("vcpu", 2, "Number of virtual CPUs")
	diskSizeMB := flag.Int("disk", 10240, "Disk size in MB")
	kernelVersion := flag.String("kernel", "5.10", "Kernel version")
	firecrackerVersion := flag.String("firecracker", "1.0", "Firecracker version")
	startCommand := flag.String("start-cmd", "", "Start command")
	hugePages := flag.Bool("huge-pages", false, "Enable huge pages")
	
	flag.Parse()

	// Set up connection to the server
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *host, *port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create client
	client := pb.NewTemplateServiceClient(conn)

	// Example: Create a template
	template := &pb.TemplateConfig{
		TemplateID:         *templateID,
		BuildID:            *buildID,
		MemoryMB:          int32(*memoryMB),
		VCpuCount:         int32(*vcpuCount),
		DiskSizeMB:        int32(*diskSizeMB),
		KernelVersion:     *kernelVersion,
		FirecrackerVersion: *firecrackerVersion,
		StartCommand:      *startCommand,
		HugePages:         *hugePages,
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
		BuildID: *buildID,
	}

	_, err = client.TemplateBuildDelete(context.Background(), deleteReq)
	if err != nil {
		log.Fatalf("Error deleting template build: %v", err)
	}
	fmt.Println("Template build deleted successfully")
}
