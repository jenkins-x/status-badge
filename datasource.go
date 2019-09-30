package main

import (
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

type datasource struct {
}

func (f *datasource) GetBadge(repo string) (Badge, error) {
	return Badge{Label: "JX", Message: "Build Passing", Color: "green", SchemaVersion: 1}, nil
}

func (f *datasource) createKubeClient() (kubernetes.Interface, error) {
	log.Printf("CreateKubeClient()")
	cfg, err := f.createKubeConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("failed to create Kubernetes Client")
	}

	return client, nil
}

func (f *datasource) createKubeConfig() (*rest.Config, error) {
	log.Printf("CreateKubeConfig()")
	masterURL := ""

	kubeconfig := f.createKubeConfigText()
	var config *rest.Config
	var err error
	if kubeconfig != nil {
		exists, err := fileExists(*kubeconfig)
		if err == nil && exists {
			// use the current context in kubeconfig
			config, err = clientcmd.BuildConfigFromFlags(masterURL, *kubeconfig)
			if err != nil {
				return nil, err
			}
		}
	}
	if config == nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, errors.Wrapf(err, "failed to check if file exists %s", path)
}

func (f *datasource) createKubeConfigText() *string {
	var kubeconfig *string
	text := ""
	if home := homeDir(); home != "" {
		text = filepath.Join(home, ".kube", "config")
	}
	kubeconfig = &text
	return kubeconfig
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	h := os.Getenv("USERPROFILE") // windows
	if h == "" {
		h = "."
	}
	return h
}
