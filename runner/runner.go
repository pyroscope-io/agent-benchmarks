package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
)

type runner struct {
	cli           *client.Client
	networkID     string
	ingestorID    string
	benchmarkedID string
}

func newRunner(cli *client.Client) *runner {
	return &runner{cli: cli}
}

func (r *runner) buildImage(ctx context.Context, path, tag string) error {
	// Let's remove the image to make sure it's properly built.
	// We rely on cache to rebuild it fast when existing.
	_, err := r.cli.ImageRemove(ctx, tag, types.ImageRemoveOptions{PruneChildren: false})
	if err != nil && !client.IsErrNotFound(err) {
		return err
	}
	tar, err := archive.Tar(path, archive.Gzip)
	if err != nil {
		return err
	}
	if _, err := r.cli.ImageBuild(ctx, tar, types.ImageBuildOptions{Tags: []string{tag}}); err != nil {
		return err
	}
	return nil
}

func (r *runner) createNetwork(ctx context.Context) error {
	res, err := r.cli.NetworkCreate(ctx, id, types.NetworkCreate{})
	if err != nil {
		return err
	}
	r.networkID = res.ID
	return nil
}

func (r *runner) connectIngestor(ctx context.Context) error {
	cfg := &network.EndpointSettings{Aliases: []string{"ingestor"}}
	return r.cli.NetworkConnect(ctx, r.networkID, r.ingestorID, cfg)
}

func (r *runner) connectBenchmarked(ctx context.Context) error {
	cfg := &network.EndpointSettings{Aliases: []string{"benchmarked"}}
	return r.cli.NetworkConnect(ctx, r.networkID, r.benchmarkedID, cfg)
}

func (r *runner) removeNetwork(ctx context.Context) {
	if r.networkID != "" {
		if err := r.cli.NetworkRemove(ctx, r.networkID); err != nil {
			fmt.Println("Unable to remove network:", err)
		}
		r.networkID = ""
	}
}

func (r *runner) createIngestor(ctx context.Context) error {
	cfg := &container.Config{Image: id + "/ingestor", ExposedPorts: nat.PortSet{"4040": struct{}{}}}
	cID, err := r.createContainer(ctx, cfg)
	if err != nil {
		return err
	}
	r.ingestorID = cID
	return nil
}

func (r *runner) removeIngestor(ctx context.Context) {
	if r.ingestorID != "" {
		if err := r.cli.ContainerRemove(ctx, r.ingestorID, types.ContainerRemoveOptions{Force: true}); err != nil {
			fmt.Println("Unable to remove ingestor container:", err)
		} else {
			r.ingestorID = ""
		}
	}
}

func (r *runner) createBenchmarked(ctx context.Context, profile bool) error {
	cfg := &container.Config{Image: id + "/benchmarked", ExposedPorts: nat.PortSet{"4040": struct{}{}}}
	if profile {
		cfg.Env = append(cfg.Env, "PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING=true")
	}
	cID, err := r.createContainer(ctx, cfg)
	if err != nil {
		return err
	}
	r.benchmarkedID = cID
	return nil
}

func (r *runner) removeBenchmarked(ctx context.Context) {
	if r.benchmarkedID != "" {
		if err := r.cli.ContainerRemove(ctx, r.benchmarkedID, types.ContainerRemoveOptions{Force: true}); err != nil {
			fmt.Println("Unable to remove benchmarked container:", err)
		}
		r.benchmarkedID = ""
	}
}

func (r *runner) createContainer(ctx context.Context, cfg *container.Config) (string, error) {
	// Give 1 CPU to each container
	hc := &container.HostConfig{
		CapAdd: []string{"sys_ptrace"},
		Resources: container.Resources{
			CPUPeriod:  100000,
			CPUQuota:   100000,
			CpusetCpus: "0",
		},
	}

	res, err := r.cli.ContainerCreate(ctx, cfg, hc, nil, nil, "")
	if err != nil {
		return "", err
	}
	return res.ID, nil
}

func (r *runner) startIngestor(ctx context.Context) error {
	if err := r.cli.ContainerStart(ctx, r.ingestorID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func (r *runner) startBenchmarked(ctx context.Context) error {
	if err := r.cli.ContainerStart(ctx, r.benchmarkedID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}
