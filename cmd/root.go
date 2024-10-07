package cmd

import (
	"fmt"
	"github.com/pete911/kubectl-image/pkg/api"
	"github.com/spf13/cobra"
	"os"
)

var (
	RootCmd = &cobra.Command{}

	Version     string
	GlobalFlags Flags
)

func init() {
	InitPersistentFlags(RootCmd, &GlobalFlags)
}

func GetRegistries() api.Registries {
	client, err := api.NewClient(GlobalFlags.KubeconfigPath, GlobalFlags.Context, GlobalFlags.Namespace)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	registries, err := client.ListRegistries(GlobalFlags.AllNamespaces)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	if len(registries) == 0 {
		namespace := fmt.Sprintf("%s namespace", client.Namespace)
		if GlobalFlags.AllNamespaces {
			namespace = "all namespaces"
		}
		fmt.Printf("found 0 images in %s\n", namespace)
		os.Exit(0)
	}
	return registries
}
