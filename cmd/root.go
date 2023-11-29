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
	client, err := api.NewClient(GlobalFlags.KubeconfigPath())
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	registries, err := client.ListRegistries(GlobalFlags.Namespace())
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
	return registries
}
