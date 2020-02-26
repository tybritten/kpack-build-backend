package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
	"github.com/ghodss/yaml"
)

const (
	namespace = "default"
)

func (i *Image) GetMetadata() *metav1.ObjectMeta {
	return i.Metadata
}

func (i *ImageList) GetMetadata() *metav1.ListMeta {
	return i.Metadata
}
func init() {
	// Register resources with the k8s package.
	k8s.Register("build.pivotal.io", "v1alpha1", "images", true, &Image{})
	k8s.RegisterList("build.pivotal.io", "v1alpha1", "images", true, &ImageList{})
}
func test_k8s() {
	_ = k8s_client()
}

func get_logins() (username string, password string) {

	ctx := context.Background()
	client := k8s_client()
	var secret corev1.Secret
	if err := client.Get(ctx, namespace, "kpack-build-backend", &secret); err != nil {
		log.Println("Failed to Get Login from Secret: ", err)
		return "", ""
	}
	username = string(secret.Data["username"])
	password = string(secret.Data["password"])
	return username, password
}

func k8s_client() *k8s.Client {
	client, err := k8s.NewInClusterClient()
	var err1 error
	if err != nil {
		client, err1 = loadClient("kubeconfig")
	}
	if err1 != nil {
		log.Println("read cluster config:", err.Error())
		log.Fatal(err1)
	}
	return client

}
func image_list() (ImageList, error) {
	ctx := context.Background()
	client := k8s_client()
	var images ImageList
	if err := client.List(ctx, namespace, &images); err != nil {
		log.Println("Failed to Retrieve image list: ", err)
		return images, err
	}
	return images, nil
}
func create_git_image(name string, url string, revision string, tag string) error {
	ctx := context.Background()
	client := k8s_client()
	image := &Image{
		APIVersion: "build.pivotal.io/v1alpha1",
		Kind:       "Image",
		Metadata: &metav1.ObjectMeta{
			Name:      &name,
			Namespace: k8s.String(namespace),
		},
		Spec: Spec{
			Tag: tag,
			Source: Source{
				Git: Git{
					Revision: revision,
					URL:      url,
				},
			},
			Builder: Builder{
				Kind: "ClusterBuilder",
				Name: "default",
			},
			ServiceAccount: "kpack-service-account",
		},
	}
	log.Println(image)
	if err := client.Create(ctx, image); err != nil {
		return fmt.Errorf("create: %v", err)
	}

	return nil
}
func update_git_image(name string, url string, revision string) error {
	ctx := context.Background()
	client := k8s_client()
	var image Image
	image, _ = get_image_status(name)
	image.Spec.Source.Git.URL = url
	image.Spec.Source.Git.Revision = revision
	if err := client.Update(ctx, &image); err != nil {
		return fmt.Errorf("update: %v", err)
	}
	return nil
}
func get_image_status(name string) (Image, error) {
	ctx := context.Background()
	client := k8s_client()
	var imagestatus Image
	if err := client.Get(ctx, namespace, name, &imagestatus); err != nil {
		return imagestatus, fmt.Errorf("Cannot get Image: %v", err)
	}
	return imagestatus, nil
}
func delete_image(name string) error {
	ctx := context.Background()
	client := k8s_client()
	image := &Image{
		APIVersion: "build.pivotal.io/v1alpha1",
		Kind:       "Image",
		Metadata: &metav1.ObjectMeta{
			Name:      &name,
			Namespace: k8s.String(namespace),
		},
	}
	if err := client.Delete(ctx, image); err != nil {
		return fmt.Errorf("delete: %v", err)
	}
	return nil
}

// loadClient parses a kubeconfig from a file and returns a Kubernetes
// client. It does not support extensions or client auth providers.
func loadClient(kubeconfigPath string) (*k8s.Client, error) {
	data, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("read kubeconfig: %v", err)
	}

	// Unmarshal YAML into a Kubernetes config object.
	var config k8s.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal kubeconfig: %v", err)
	}
	return k8s.NewClient(&config)
}

type Image struct {
	APIVersion string             `json:"apiVersion"`
	Kind       string             `json:"kind"`
	Metadata   *metav1.ObjectMeta `json:"metadata"`
	Spec       Spec               `json:"spec"`
	Status     Status             `json:"status",omitempty`
}

type Spec struct {
	Builder        Builder `json:"builder"`
	ServiceAccount string  `json:"serviceAccount"`
	Source         Source  `json:"source"`
	Tag            string  `json:"tag"`
}

type Status struct {
	BuildCounter       int        `json:"buildCounter"`
	Conditions         Conditions `json:"conditions"`
	LatestBuildRef     string     `json:"latestBuildRef"`
	LatestImage        string     `json:"latestImage"`
	LatestStack        string     `json:"latestStack"`
	ObservedGeneration int        `json:"observedGeneration"`
}

type Builder struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type Source struct {
	Git Git `json:"git"`
}

type Conditions []struct {
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	Status             string    `json:"status"`
	Type               string    `json:"type"`
}

type Git struct {
	Revision string `json:"revision"`
	URL      string `json:"url"`
}

type ImageList struct {
	Metadata *metav1.ListMeta `json:"metadata"`
	Items    []Image          `json:"items"`
}
type ImageTaggingStrategy string

const (
	None        ImageTaggingStrategy = "None"
	BuildNumber ImageTaggingStrategy = "BuildNumber"
)
