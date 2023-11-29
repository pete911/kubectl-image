package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var logLevels = map[string]slog.Level{"debug": slog.LevelDebug, "info": slog.LevelInfo, "warn": slog.LevelWarn, "error": slog.LevelError}

type Flags struct {
	kubeconfigPath string
	logLevel       string
	namespace      string
	allNamespaces  bool
}

func (f Flags) KubeconfigPath() string {
	return f.kubeconfigPath
}

func (f Flags) Logger() *slog.Logger {
	if level, ok := logLevels[strings.ToLower(f.logLevel)]; ok {
		opts := &slog.HandlerOptions{Level: level}
		return slog.New(slog.NewJSONHandler(os.Stderr, opts))
	}

	fmt.Printf("invalid log level %s", f.logLevel)
	os.Exit(1)
	return nil
}

func (f Flags) Namespace() string {
	if f.allNamespaces {
		return ""
	}
	return f.namespace
}

func InitPersistentFlags(cmd *cobra.Command, flags *Flags) {
	defaultKubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	cmd.PersistentFlags().StringVar(
		&flags.kubeconfigPath,
		"kubeconfig",
		getStringEnv("KUBECONFIG", defaultKubeconfig),
		"path to kubeconfig file",
	)
	cmd.PersistentFlags().StringVar(
		&flags.logLevel,
		"log-level",
		"warn",
		"log level - debug, info, warn, error",
	)
	cmd.PersistentFlags().StringVarP(
		&flags.namespace,
		"namespace",
		"n",
		"default",
		"kubernetes namespace",
	)
	cmd.PersistentFlags().BoolVarP(
		&flags.allNamespaces,
		"all-namespaces",
		"A",
		false,
		"all kubernetes namespaces",
	)
}

func getStringEnv(envName string, defaultValue string) string {
	env, ok := os.LookupEnv(envName)
	if !ok {
		return defaultValue
	}
	return env
}
