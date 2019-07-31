package main

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/system"
	"github.com/pkg/errors"
)

type Directive struct {
	Flavor     string      `json:"flavor"`
	GoSpec     GoDirective `json:"goSpec"`
	Repository string      `json:"repository"`
}

type GoDirective struct {
	Binaries []string `json:"binaries"`
	Module   string   `json:"module,omitempty"`
}

func generateDefinition(drc Directive) (*llb.Definition, error) {
	if drc.Flavor != "go" {
		return nil, errors.Errorf("expected flavor to equal %q, but found", "go", drc.Flavor)
	}

	builder := llb.Image("golang:1.12-stretch").
		AddEnv("PATH", "/usr/local/go/bin:"+system.DefaultPathEnv).
		AddEnv("GOPATH", "/go")
	if drc.Repository != "" {
		builder = builder.Run(llb.Shlexf("git clone %s /build", drc.Repository)).Root()
	}
	builder = builder.Dir("/build")
	for _, b := range drc.GoSpec.Binaries {
		builder = builder.Run(llb.Shlexf("go install %s", b)).Root()
	}

	def, err := builder.Marshal(llb.LinuxAmd64)
	if err != nil {
		return nil, errors.Wrap(err, "builder.Marshal")
	}
	return def, nil
}
