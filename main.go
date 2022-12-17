package main

import (
	"context"
	"github.com/pete911/kubectl-image/pkg/api"
	"github.com/pete911/kubectl-image/pkg/out"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"time"
)

func main() {

	registries, err := listRegistries()
	if err != nil {
		log.Fatal(err)
	}
	out.Print(registries)
}

func listRegistries() (api.Registries, error) {

	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	client, err := api.NewClient(kubeconfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.ListRegistries(ctx)
}
