package main

import (
	"log"
	"os"

	"github.com/moby/buildkit/client/buildid"
	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	pb "github.com/moby/buildkit/frontend/gateway/pb"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("error: %s", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	// err := GraphTest()
	// return errors.Wrap(err, "GraphTest")

	// ctx := appcontext.Context()
	// err := grpcclient.RunFromEnvironment(ctx, Builder)
	// return errors.Wrap(err, "RunFromEnvironment")

	ctx := buildid.AppendToOutgoingContext(appcontext.Context(), "cfdb77b4-d8ab-46d8-95a1-ad75d18514c8")
	opts := []grpc.DialOption{
		// grpc.WithDialer(...),
		grpc.WithInsecure(),
	}
	conn, err := grpc.DialContext(ctx, "127.0.0.1:1234", opts...)
	if err != nil {
		return errors.Wrap(err, "grpc.DialContext")
	}

	client, err := grpcclient.New(ctx, nil, "", "", pb.NewLLBBridgeClient(conn), nil)
	if err != nil {
		return errors.Wrap(err, "grpcclient.New")
	}

	return errors.Wrap(client.Run(ctx, Builder), "Run")
}
