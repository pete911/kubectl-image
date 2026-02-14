package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/pete911/kubectl-image/pkg/api"
	"github.com/pete911/kubectl-image/pkg/out"
	"github.com/spf13/cobra"
)

var (
	cmdList = &cobra.Command{
		Use:   "list",
		Short: "list images",
		Long:  "",
		Run:   runListCmd,
	}
)

func init() {
	RootCmd.AddCommand(cmdList)
}

func runListCmd(_ *cobra.Command, _ []string) {
	logger := GlobalFlags.Logger()
	registries := GetRegistries()
	PrintList(logger, registries)
}

func PrintList(logger *slog.Logger, registries api.Registries) {
	table := out.NewTable(logger, 50)
	table.AddRow("REGISTRY", "REPOSITORY", "TAG", "ID", "SIZE", "PODS", "FAILED", "RESTART")
	for _, registry := range registries {
		for _, repository := range registry.ListRepositories() {
			for _, id := range repository.ListIDs() {
				containers := id.ListContainers()
				if len(containers) == 0 {
					// very unlikely
					continue
				}
				// nodes can have image either by tag or id, we just need first container to find the size
				size := id.ListContainers()[0].ImageSizeBytes
				sizeMb := fmt.Sprintf("%.2fMB", float64(size)/1000000)
				tags := strings.Join(id.ListTags(), ", ")
				pods := getNumPods(containers, false)
				failedPods := getNumPods(containers, true)
				restarts := getNumRestarts(containers)
				table.AddRow(registry.Name, repository.Name, tags, id.Name, sizeMb, pods, failedPods, restarts)
			}
		}
	}
	table.Print()
}

func getNumRestarts(containers api.Containers) string {
	var restarts int
	for _, container := range containers {
		restarts = restarts + container.RestartCount
	}
	return fmt.Sprintf("%d", restarts)
}

func getNumPods(containers []api.Container, failed bool) string {
	searchedSet := make(map[string]struct{})
	for _, container := range containers {
		key := fmt.Sprintf("%s/%s", container.Pod.Namespace, container.Pod.Name)
		if _, ok := searchedSet[key]; ok {
			continue
		}
		if failed {
			if strings.ToLower(container.Pod.Phase) == "failed" {
				searchedSet[key] = struct{}{}
			}
			continue
		}
		searchedSet[key] = struct{}{}
	}
	return fmt.Sprintf("%d", len(searchedSet))
}
