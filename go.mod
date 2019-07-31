module github.com/ripta/vib

go 1.12

require (
	github.com/containerd/console v0.0.0-20181022165439-0650fd9eeb50
	github.com/containerd/containerd v1.3.0-0.20190426060238-3a3f0aac8819
	github.com/docker/docker v1.14.0-0.20190319215453-e7b5f7dbe98c
	github.com/moby/buildkit v0.5.1
	github.com/opencontainers/image-spec v1.0.1
	github.com/pkg/errors v0.8.1
	golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f
	google.golang.org/grpc v1.12.0
)

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
