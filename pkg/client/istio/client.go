package istio

import (
	"context"

	"istio.io/client-go/pkg/apis/networking/v1alpha3"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
)

type ClientInterface interface {
	GetEnvoyFilter(ctx context.Context, namespace string, name string) (*v1alpha3.EnvoyFilter, error)
	CreateEnvoyFilter(ctx context.Context, namespace string, envoyFilter *v1alpha3.EnvoyFilter) (*v1alpha3.EnvoyFilter, error)
	UpdateEnvoyFilter(ctx context.Context, namespace string, envoyFilter *v1alpha3.EnvoyFilter) (*v1alpha3.EnvoyFilter, error)
}
type Client struct {
	cfg    *rest.Config
	client versionedclient.Interface
}

func NewClient(cfg *rest.Config) (ClientInterface, error) {
	client, err := versionedclient.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{cfg: cfg, client: client}, nil
}

func (r *Client) GetEnvoyFilter(ctx context.Context, namespace string, name string) (*v1alpha3.EnvoyFilter, error) {
	return r.client.NetworkingV1alpha3().EnvoyFilters(namespace).Get(ctx, name, v1.GetOptions{})
}

func (r *Client) UpdateEnvoyFilter(ctx context.Context, namespace string, envoyFilter *v1alpha3.EnvoyFilter) (*v1alpha3.EnvoyFilter, error) {
	return r.client.NetworkingV1alpha3().EnvoyFilters(namespace).Update(ctx, envoyFilter, v1.UpdateOptions{})
}

func (r *Client) CreateEnvoyFilter(ctx context.Context, namespace string, envoyFilter *v1alpha3.EnvoyFilter) (*v1alpha3.EnvoyFilter, error) {
	return r.client.NetworkingV1alpha3().EnvoyFilters(namespace).Create(ctx, envoyFilter, v1.CreateOptions{})
}
