package K8s

import (
	"context"
	"os"
	v1 "task/controller/api/v1"

	"k8s.io/apimachinery/pkg/runtime" // Import runtime for scheme
	// Import for GroupVersionKind
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	// Import for error handling
)

type K8s struct {
	kubeconfigPath string
	client         *kubernetes.Clientset
	scheme         *runtime.Scheme // Add scheme to K8s struct
}

// Add the following function to create a Kubernetes client for local and in-cluster setup
func NewK8sClient(kubeconfigPath string) (*K8s, error) {
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

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	k := &K8s{
		kubeconfigPath: kubeconfigPath,
		client:         clientset,
		scheme:         runtime.NewScheme(), // Initialize the scheme
	}

	// Register your Task type with the scheme
	if err := v1.AddToScheme(k.scheme); err != nil {
		return nil, err
	}

	return k, nil
}

// CreateTask creates a new Task resource in the Kubernetes cluster
func (k *K8s) CreateTask(task *v1.Task) (*v1.Task, error) {
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
func (k *K8s) GetTask(namespace, name string) (*v1.Task, error) {
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
func (k *K8s) UpdateTask(task *v1.Task) (*v1.Task, error) {
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
func (k *K8s) DeleteTask(namespace, name string) error {
	return k.client.RESTClient().
		Delete().
		Resource("tasks").
		Namespace(namespace).
		Name(name).
		Do(context.TODO()).
		Error()
}
