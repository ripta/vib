package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/util/system"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// GraphTest dumps the LLB instructions to STDOUT.
func GraphTest() error {
	st := llb.Image("ubuntu:18.04")

	def, err := st.Marshal()
	if err != nil {
		return errors.Wrapf(err, "marshaling state")
	}

	return llb.WriteTo(def, os.Stdout)
}

// Builder generates an image.
func Builder(ctx context.Context, c client.Client) (*client.Result, error) {
	st := llb.Image("ubuntu:18.04")

	def, err := st.Marshal(llb.LinuxAmd64)
	if err != nil {
		return nil, errors.Wrapf(err, "marshaling state")
	}

	req := client.SolveRequest{
		Definition: def.ToPB(),
	}
	res, err := c.Solve(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "resolving dockerfile")
	}

	img := specs.Image{
		Architecture: "amd64",
		OS:           "linux",
	}
	img.RootFS.Type = "layers"
	img.Config.WorkingDir = "/"
	img.Config.Env = []string{"PATH=" + system.DefaultPathEnv}
	img.Config.Cmd = []string{"ls"}

	cfg, err := json.Marshal(img)
	if err != nil {
		return nil, errors.Wrapf(err, "marshaling image definition")
	}
	log.Println(cfg)

	plat := []string{
		exptypes.ExporterImageConfigKey,
		platforms.Format(platforms.DefaultSpec()),
	}
	res.AddMeta(strings.Join(plat, "/"), cfg)

	// ref, err := res.SingleRef()
	// if err != nil {
	// 	return nil, err
	// }
	// res.SetRef(ref)

	return res, nil
}
