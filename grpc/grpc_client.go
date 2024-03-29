package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"strings"
	"time"
	"trace-app/grpc/pb"
)

func DoSingleGrpcRequest(url string, interval int, cpc bool) {
	if url == "" {
		url = "localhost:8080"
	}
	grpcUrl := url
	log.Printf("grpc request => %s is begin.", url)
	go func() {
		for {
			time.Sleep(time.Duration(interval) * time.Second)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if cpc {
				grpcUrl = "127.0.0.1:12345"
			}
			opts := grpc.WithInsecure()
			cc, err := grpc.Dial(grpcUrl, opts, grpc.WithAuthority(url))
			if err != nil {
				log.Println(err)
				continue
			}
			defer cc.Close()

			client := pb.NewGreetingServiceClient(cc)
			request := &pb.GreetingServiceRequest{Name: "Tubiers"}

			resp, err := client.Greeting(ctx, request)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Printf("Receive grpc response => %s from => %s \n", resp.Message, url)
		}
	}()
}

func DoGrpcRequest(url string, interval int, cpc bool) {
	urls := strings.Split(url, ",")
	for _, u := range urls {
		DoSingleGrpcRequest(u, interval, cpc)
	}
}
