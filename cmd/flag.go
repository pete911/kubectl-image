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
	KubeconfigPath string
	logLevel       string
	Namespace      string
	AllNamespaces  bool
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

func InitPersistentFlags(cmd *cobra.Command, flags *Flags) {
	defaultKubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	cmd.PersistentFlags().StringVar(
		&flags.KubeconfigPath,
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
		&flags.Namespace,
		"namespace",
		"n",
		"",
		"kubernetes namespace",
	)
	cmd.PersistentFlags().BoolVarP(
		&flags.AllNamespaces,
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
