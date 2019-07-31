package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/containerd/console"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/buildid"
	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	pb "github.com/moby/buildkit/frontend/gateway/pb"
	"github.com/moby/buildkit/util/appcontext"
	ui "github.com/moby/buildkit/util/progress/progressui"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func exportAsImage(name string) []client.ExportEntry {
	return []client.ExportEntry{
		{
			Type: client.ExporterImage,
			Attrs: map[string]string{
				"name": name,
			},
		},
	}
}

func main() {
	buildkitAddr := flag.String("addr", "tcp://127.0.0.1:1234", "Buildkitd address")
	definitionFile := flag.String("definition", "", "Definition filename")
	outputImageName := flag.String("output", "", "Fully-qualified output image name, e.g., docker.io/ripta/vib:latest")
	flag.Parse()

	if *buildkitAddr == "" {
		log.Fatal("error: --addr must not be empty")
	}

	if *definitionFile == "" {
		log.Fatal("error: --definition must not be empty")
	}
	if *outputImageName == "" {
		log.Fatal("error: --output must not be empty")
	}

	if err := run(*buildkitAddr, *definitionFile, *outputImageName); err != nil {
		log.Fatalf("error: %s", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(addr, filename, img string) error {
	r, err := os.Open(filename)
	if err != nil {
		return errors.Wrapf(err, "could not load file %s", filename)
	}

	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrapf(err, "could not read from file %s", filename)
	}

	var drc Directive
	if err := json.Unmarshal(raw, &drc); err != nil {
		return errors.Wrapf(err, "could not parse JSON from file %s", filename)
	}

	ctx := appcontext.Context()
	opts := []client.ClientOpt{
		client.WithFailFast(),
	}
	c, err := client.New(ctx, addr, opts...)
	if err != nil {
		return errors.Wrap(err, "client.New")
	}

	solver := client.SolveOpt{
		Exports: exportAsImage(img),
	}
	progressCh := make(chan *client.SolveStatus)

	gr, ctx := errgroup.WithContext(ctx)
	gr.Go(func() error {
		req, err := generateDefinition(drc)
		if err != nil {
			close(progressCh)
			return errors.Wrap(err, "generateDefinition")
		}

		rsp, err := c.Solve(ctx, req, solver, progressCh)
		if err != nil {
			return errors.Wrap(err, "client.Build")
		}
		for k, v := range rsp.ExporterResponse {
			fmt.Printf("Build response %s = %+v\n", k, v)
		}
		return nil
	})
	gr.Go(progressDisplayer(context.Background(), "build", progressCh))

	return gr.Wait()
}

func progressDisplayer(ctx context.Context, phase string, ch chan *client.SolveStatus) func() error {
	return func() error {
		cons, err := console.ConsoleFromFile(os.Stderr)
		if err != nil {
			return errors.Wrap(err, "console.ConsoleFromFile")
		}
		return ui.DisplaySolveStatus(ctx, phase, cons, os.Stderr, ch)
	}
}

func run2() error {
	ctx, cancel := context.WithTimeout(appcontext.Context(), 2*time.Second)
	defer cancel()

	opts := []client.ClientOpt{
		client.WithFailFast(),
	}
	c, err := client.New(ctx, "tcp://127.0.0.1:1234", opts...)
	if err != nil {
		return errors.Wrap(err, "client.New")
	}

	du, err := c.DiskUsage(ctx)
	if err != nil {
		return errors.Wrap(err, "client.DiskUsage")
	}
	fmt.Printf("disk usage returned %d usage info objects\n", len(du))

	fs := client.WithFilter([]string{})
	ws, err := c.ListWorkers(ctx, fs)
	if err != nil {
		return errors.Wrap(err, "client.ListWorkers")
	}
	for _, w := range ws {
		var ps bytes.Buffer
		for _, p := range w.Platforms {
			ps.WriteString(p.Architecture + ",")
		}
		fmt.Printf("%20s\t%20s\n", w.ID, ps.String())
	}

	return nil
}

func run3() error {
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
