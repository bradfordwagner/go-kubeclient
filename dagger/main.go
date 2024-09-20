// A generated module for Dagger functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/dagger/internal/dagger"
)

type Dagger struct{}

func (m *Dagger) Test(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,

) (string, error) {
	k3s := dag.K3S("test")
	kServer := k3s.Server()

	kServer, err := kServer.Start(ctx)
	if err != nil {
		return "", err
	}

	configFile := k3s.Config()

	//time.Sleep(10 * time.Second)

	//return dag.Container().From("bitnami/kubectl").
	//	WithEnvVariable("KUBECONFIG", "/.kube/config").
	//	WithFile("/.kube/config", configFile, dagger.ContainerWithFileOpts{
	//		Permissions: 0755,
	//	}).
	//	WithExec([]string{"kubectl", "get", "pods", "-A"}).
	//	Stdout(ctx)
	return dag.Container().From("ghcr.io/bradfordwagner/go-builder:3.2.0-alpine_3.19").
		WithEnvVariable("KUBECONFIG", "/.kube/config").
		WithDirectory("/src", src).
		WithFile("/.kube/config", configFile, dagger.ContainerWithFileOpts{
			Permissions: 0755,
		}).
		//WithExec([]string{"apk", "add", "kubectl"}).
		//WithExec([]string{"kubectl", "get", "pods", "-A"}).
		// go test ./...
		//WithExec([]string{"go", "version"}).
		WithExec([]string{"go", "test", "./..."}).
		//WithExec([]string{"ls", "-la", "/src"}).
		Stdout(ctx)
}
