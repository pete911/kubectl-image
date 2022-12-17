package cmd

import (
	"context"
	"fmt"
	"github.com/pete911/kubectl-image/pkg/api"
	"github.com/pete911/kubectl-image/pkg/out"
	"github.com/spf13/cobra"
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
	podFlags PodFlags
)

func init() {

	RootCmd.AddCommand(cmdList)
	InitPodFlags(cmdList, &podFlags)
}

func runListCmd(_ *cobra.Command, args []string) error {

	// no namespace means all namespaces
	if podFlags.AllNamespaces {
		podFlags.Namespace = ""
	}

	// additional arguments are considered to be pod names, add to field selector flags
	for _, v := range args {
		fieldSelectors := strings.Split(podFlags.FieldSelector, ",")
		fieldSelectors = append(fieldSelectors, fmt.Sprintf("metadata.name=%s", v))
		podFlags.FieldSelector = strings.Join(fieldSelectors, ",")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	registries, err := listRegistries(ctx, podFlags.Namespace, podFlags.Label, podFlags.FieldSelector)
	if err != nil {
		return fmt.Errorf("list registries: %w", err)
	}

	if len(registries) == 0 {
		namespace := "all namespaces"
		if podFlags.Namespace != "" {
			namespace = fmt.Sprintf("%s namespace", podFlags.Namespace)
		}
		fmt.Printf("found 0 images in %s\n", namespace)
		return nil
	}

	out.Print(registries)
	return nil
}

func listRegistries(ctx context.Context, namespace, labelSelector, fieldSelector string) (api.Registries, error) {

	client, err := api.NewClient(KubeconfigPath)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.ListRegistries(ctx, namespace, labelSelector, fieldSelector)
}
