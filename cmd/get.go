package cmd

import (
	"fmt"
	"github.com/pete911/kubectl-image/pkg/api"
	"github.com/pete911/kubectl-image/pkg/out"
	"github.com/spf13/cobra"
	"log/slog"
	"strings"
)

var (
	cmdGet = &cobra.Command{
		Use:   "get",
		Short: "get images",
		Long:  "",
		Run:   runGetCmd,
	}
)

func init() {
	RootCmd.AddCommand(cmdGet)
}

func runGetCmd(_ *cobra.Command, _ []string) {
	logger := GlobalFlags.Logger()
	registries, nodes := GetRegistriesAndNodes()
	PrintGet(logger, registries, nodes)
}

func PrintGet(logger *slog.Logger, registries api.Registries, nodes api.Nodes) {
	for _, registry := range registries {
		fmt.Printf("registry:   %s\n", registry.Name)
		for _, repository := range registry.ListRepositories() {
			fmt.Printf("  repository: %s\n", repository.Name)
			for _, id := range repository.ListIDs() {
				fmt.Printf("    id: %s tags: %s\n", id.Name, strings.Join(id.ListTags(), ", "))
				table := out.NewTable(logger, 80)
				table.AddRow("    CONTAINER", "RESTART", "STATE", "MESSAGE", "INIT", "POD", "PHASE")
				for _, container := range id.ListContainers() {
					containerName := fmt.Sprintf("    %s", container.Name)
					pod := fmt.Sprintf("%s/%s", container.Pod.Namespace, container.Pod.Name)
					restart := fmt.Sprintf("%d", container.RestartCount)
					init := fmt.Sprintf("%t", container.IsInit)
					table.AddRow(containerName, restart, container.State, container.Message, init, pod, container.Pod.Phase)
				}
				table.Print()
			}
		}
		fmt.Println()
	}
}
