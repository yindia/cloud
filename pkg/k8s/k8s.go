package k8s

import (
	"context"
	"os"
	v1 "task/controller/api/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type k8s struct {
	kubeconfigPath string
	client         *kubernetes.Clientset
}

// Add the following function to create a Kubernetes client for local and in-cluster setup
func NewK8sClient(kubeconfigPath string) (*k8s, error) {
	var config *rest.Config
	var err error

	// In-cluster config
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else { // Local config
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, err
		}
	}

	return &k8s{
		kubeconfigPath: kubeconfigPath,
		client:         kubernetes.NewForConfigOrDie(config),
	}, nil
}

// CreateTask creates a new Task resource in the Kubernetes cluster
func (k *k8s) CreateTask(task *v1.Task) (*v1.Task, error) {
	tasksClient := k.client.RESTClient().
		Post().
		Resource("tasks").
		Namespace(task.Namespace).
		Body(task)

	result := &v1.Task{}
	err := tasksClient.Do(context.TODO()).Into(result)
	return result, err
}

// GetTask retrieves a Task resource from the Kubernetes cluster
func (k *k8s) GetTask(namespace, name string) (*v1.Task, error) {
	result := &v1.Task{}
	err := k.client.RESTClient().
		Get().
		Resource("tasks").
		Namespace(namespace).
		Name(name).
		Do(context.TODO()).
		Into(result)
	return result, err
}

// UpdateTask updates an existing Task resource in the Kubernetes cluster
func (k *k8s) UpdateTask(task *v1.Task) (*v1.Task, error) {
	result := &v1.Task{}
	err := k.client.RESTClient().
		Put().
		Resource("tasks").
		Namespace(task.Namespace).
		Name(task.Name).
		Body(task).
		Do(context.TODO()).
		Into(result)
	return result, err
}

// DeleteTask deletes a Task resource from the Kubernetes cluster
func (k *k8s) DeleteTask(namespace, name string) error {
	return k.client.RESTClient().
		Delete().
		Resource("tasks").
		Namespace(namespace).
		Name(name).
		Do(context.TODO()).
		Error
}
