package main

import (
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"os"
	"strings"
)
import (
	"context"
	"log"

	"cloud.google.com/go/cloudtasks/apiv2beta3"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	service = "helloworld.Greeter"
	method = "SayHello"
)

func main() {
	baseurl := flag.String("baseurl", "", "target baseurl")
	queue := flag.String("queue", "default", "queue name")
	project := flag.String("project", "default", "project ID")
	location := flag.String("location", "", "Cloud Tasks location")
	flag.Parse()
	if len(flag.Args()) > 2 || *baseurl == "" || *project == "" || *location == "" {
		flag.Usage()
		os.Exit(1)
	}
	name := flag.Arg(0)
	if name == "" {
		name = "hello"
	}

	reqpb := &pb.HelloRequest{Name: name}
	b, err := proto.Marshal(reqpb)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	c, err := cloudtasks.NewClient(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	req := &taskspb.CreateTaskRequest{
		Parent: fmt.Sprintf(`projects/%s/locations/%s/queues/%s`, *project, *location, *queue),
		Task: &taskspb.Task{
			PayloadType:          &taskspb.Task_HttpRequest{
				// assemble gRPC fallback request
				// https://googleapis.github.io/HowToRPC.html#grpc-fallback-experimental
				HttpRequest: &taskspb.HttpRequest{
					Url:        fmt.Sprintf("%s/$rpc/%s/%s", *baseurl, service, method),
					HttpMethod: taskspb.HttpMethod_POST,
					Headers:    map[string]string{"content-type": "application/x-protobuf"},
					Body:       b,
				},
			},
		},
	}
	resp, err := c.CreateTask(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	_ = resp
	nameElems := strings.Split(resp.Name, "/")
	fmt.Println(nameElems[len(nameElems)-1])
}