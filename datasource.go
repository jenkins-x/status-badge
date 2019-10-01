package main

import (
	"fmt"
	v1 "github.com/jenkins-x/jx/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx/pkg/client/clientset/versioned"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type datasource struct {
}

func (f *datasource) GetBadge(repo string) (*Badge, error) {
	log.Printf("getting badge for '%s'", repo)
	config, err := f.createKubeConfig()
	if err != nil {
		log.Printf("error getting badge - %s", err)
	}

	client, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	options := metav1.ListOptions{}
	options.LabelSelector = fmt.Sprintf("repository=%s,branch=master", repo)

	list, err := client.JenkinsV1().PipelineActivities("").List(options)
	if err != nil {
		return nil, err
	}

	log.Printf("sorting '%s' items", len(list.Items))
	if len(list.Items) > 0 {

		activities := list.Items
		sort.Sort(pipelineActivitySorter(activities))

		last := activities[len(activities)-1]

		if last.Spec.Status == v1.ActivityStatusTypeSucceeded {
			return &Badge{Label: "JX", Message: fmt.Sprintf("Build %s", last.Spec.Status), Color: "success", SchemaVersion: 1}, nil
		} else if last.Spec.Status == v1.ActivityStatusTypeFailed {
			return &Badge{Label: "JX", Message: fmt.Sprintf("Build %s", last.Spec.Status), Color: "critical", SchemaVersion: 1}, nil
		} else {
			return &Badge{Label: "JX", Message: fmt.Sprintf("Build %s", last.Spec.Status), Color: "important", SchemaVersion: 1}, nil
		}

	}

	return nil, nil
}

type pipelineActivitySorter []v1.PipelineActivity

func (s pipelineActivitySorter) Len() int {
	return len(s)
}
func (s pipelineActivitySorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s pipelineActivitySorter) Less(i, j int) bool {
	return s[i].CreationTimestamp.Before(&s[j].CreationTimestamp)
}

func (f *datasource) createKubeConfig() (*rest.Config, error) {
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
