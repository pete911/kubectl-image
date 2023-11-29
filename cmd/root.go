package cmd

import (
	"context"
	"fmt"
	"github.com/pete911/kubectl-image/pkg/api"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	RootCmd = &cobra.Command{}

	Version     string
	GlobalFlags Flags
)

func init() {
	InitPersistentFlags(RootCmd, &GlobalFlags)
}

func GetRegistriesAndNodes() (api.Registries, api.Nodes) {
	client, err := api.NewClient(GlobalFlags.KubeconfigPath())
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	registries, err := listRegistries(client)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	if len(registries) == 0 {
		namespace := "all namespaces"
		if GlobalFlags.Namespace() != "" {
			namespace = fmt.Sprintf("%s namespace", GlobalFlags.Namespace())
		}
		fmt.Printf("found 0 images in %s\n", namespace)
		os.Exit(0)
	}

	nodes, err := listNodes(client)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	return registries, nodes
}

func listRegistries(client api.Client) (api.Registries, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.ListRegistries(ctx, GlobalFlags.Namespace())
}

func listNodes(client api.Client) (api.Nodes, error) {
	if !GlobalFlags.Size() {
		return api.Nodes{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.ListNodes(ctx)
}
