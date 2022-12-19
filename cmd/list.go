package cmd

import (
	"context"
	"fmt"
	"github.com/pete911/kubectl-image/pkg/api"
	"github.com/pete911/kubectl-image/pkg/out"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

var (
	cmdList = &cobra.Command{
		Use:   "list",
		Short: "list images",
		Long:  "",
		RunE:  runListCmd,
	}
	listFlags ListFlags
)

func init() {

	RootCmd.AddCommand(cmdList)
	InitPodFlags(cmdList, &listFlags)
}

func runListCmd(_ *cobra.Command, args []string) error {

	// no namespace means all namespaces
	if listFlags.AllNamespaces {
		listFlags.Namespace = ""
	}

	// additional arguments are considered to be pod names, add to field selector flags
	for _, v := range args {
		fieldSelectors := strings.Split(listFlags.FieldSelector, ",")
		fieldSelectors = append(fieldSelectors, fmt.Sprintf("metadata.name=%s", v))
		listFlags.FieldSelector = strings.Join(fieldSelectors, ",")
	}

	client, err := api.NewClient(KubeconfigPath)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	registries, err := listRegistries(client)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	nodes, err := listNodes(client)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	if len(registries) == 0 {
		namespace := "all namespaces"
		if listFlags.Namespace != "" {
			namespace = fmt.Sprintf("%s namespace", listFlags.Namespace)
		}
		fmt.Printf("found 0 images in %s\n", namespace)
		return nil
	}

	out.PrintRegistries(registries, nodes)
	return nil
}

func listRegistries(client api.Client) (api.Registries, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.ListRegistries(ctx, listFlags.Namespace, listFlags.Label, listFlags.FieldSelector)
}

func listNodes(client api.Client) (api.Nodes, error) {

	if !listFlags.Size {
		return api.Nodes{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.ListNodes(ctx)
}
