package cmd

import (
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"
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
	listHeader = []string{"REGISTRY", "REPOSITORY", "TAG", "ID", "SIZE", "PODS", "FAILED", "RESTART"}
	listSort   string
)

func init() {
	cmdList.Flags().StringVar(&listSort, "sort", "", fmt.Sprintf("sort by column name %v", listHeader))
	RootCmd.AddCommand(cmdList)
}

func runListCmd(_ *cobra.Command, _ []string) {
	sortColumn := -1
	if listSort != "" {
		sortColumn = slices.Index(listHeader, strings.ToUpper(listSort))
		if sortColumn < 0 {
			fmt.Printf("invalid sort column %q, valid columns: %v\n", listSort, listHeader)
			os.Exit(1)
		}
	}
	logger := GlobalFlags.Logger()
	registries := GetRegistries()
	PrintList(logger, registries, sortColumn)
}

// PrintList prints images table, rows are sorted by sortColumn index (if sortColumn is negative, rows are not sorted)
func PrintList(logger *slog.Logger, registries api.Registries, sortColumn int) {
	var rows [][]string
	for _, registry := range registries.List() {
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
				rows = append(rows, []string{registry.Name, repository.Name, tags, id.Name, sizeMb, pods, failedPods, restarts})
			}
		}
	}

	if sortColumn >= 0 {
		slices.SortStableFunc(rows, func(a, b []string) int {
			return compareColumn(a[sortColumn], b[sortColumn])
		})
	}

	table := out.NewTable(logger, 50)
	table.AddRow(listHeader...)
	for _, row := range rows {
		table.AddRow(row...)
	}
	table.Print()
}

// compareColumn compares values as numbers (size has "MB" suffix), or as strings if they are not numbers
func compareColumn(a, b string) int {
	af, aErr := strconv.ParseFloat(strings.TrimSuffix(a, "MB"), 64)
	bf, bErr := strconv.ParseFloat(strings.TrimSuffix(b, "MB"), 64)
	if aErr == nil && bErr == nil {
		return cmp.Compare(af, bf)
	}
	return strings.Compare(a, b)
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
