package cmd

import "github.com/spf13/cobra"

type ListFlags struct {
	Namespace     string
	AllNamespaces bool
	Size          bool
	Label         string
	FieldSelector string
}

func InitPodFlags(cmd *cobra.Command, flags *ListFlags) {

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
	cmd.Flags().BoolVarP(
		&flags.Size,
		"size",
		"",
		true,
		"print image size",
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
