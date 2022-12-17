package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

var (
	Version        string
	KubeconfigPath string

	RootCmd = &cobra.Command{
		Use:   "image",
		Short: "container images in kubernetes cluster",
		Long:  "",
	}
)

func init() {
	defaultKubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	RootCmd.PersistentFlags().StringVar(
		&KubeconfigPath,
		"kubeconfig",
		getStringEnv("KUBECONFIG", defaultKubeconfig),
		"path to kubeconfig file",
	)
}

func getStringEnv(envName string, defaultValue string) string {

	env, ok := os.LookupEnv(envName)
	if !ok {
		return defaultValue
	}
	return env
}
