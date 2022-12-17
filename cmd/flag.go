package cmd

import "github.com/spf13/cobra"

type PodFlags struct {
	Namespace     string
	AllNamespaces bool
	Label         string
	FieldSelector string
}

func InitPodFlags(cmd *cobra.Command, flags *PodFlags) {

	cmd.Flags().StringVarP(
		&flags.Namespace,
		"namespace",
		"n",
		"default",
		"kubernetes namespace",
	)
	cmd.Flags().BoolVarP(
		&flags.AllNamespaces,
		"all-namespaces",
		"A",
		false,
		"all kubernetes namespaces",
	)
	cmd.Flags().StringVarP(
		&flags.Label,
		"label",
		"l",
		"",
		"kubernetes label",
	)
	cmd.Flags().StringVarP(
		&flags.FieldSelector,
		"field-selector",
		"",
		"",
		"kubernetes field selector",
	)
}
