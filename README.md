vib - Vastly Indifferent Builder

Experiments in using Moby's Low-Level Builder (LLB) directly.

To use, you'll first need to run buildkitd. If your docker already has it
enabled, you need to find out its listening address; the default listening
socket is:

    unix:///run/buildkit/buildkitd.sock

If your docker does not have buildkitd but you can run a privileged container,
then you could run buildkitd in docker:

    docker run --detach --rm --privileged --publish 1234:1234 ripta/buildkit:commit-fb5324c609465f9b0713cbce5f8a36eb119be144 --addr tcp://0.0.0.0:1234

for any arbitrary free port 1234, which would make the buildkit gRPC endpoint
available on the host machine:

    tcp://127.0.0.1:1234

Compile vib from the root of this repository:

    go build ./cmd/vib

and look at what buildkit workers are available to you:

    ./vib --definition examples/demo.json --output docker.io/ripta/vib-demo:test --addr tcp://127.0.0.1:1234
