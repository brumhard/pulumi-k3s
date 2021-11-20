package k3s

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	test "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	mapper meta.RESTMapper
	dyn    dynamic.Interface
}

func newK8sClient(kubeconfig string) (*K8sClient, error) {
	clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(kubeconfig))
	if err != nil {
		return nil, err
	}

	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	return &K8sClient{
		mapper: mapper,
		dyn:    dyn,
	}, nil
}

func (k *K8sClient) CreateOrUpdateFromFile(ctx context.Context, fileBytes []byte) error {
	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	reader := test.NewYAMLReader(bufio.NewReader(bytes.NewReader(fileBytes)))
	for {
		documentBytes, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}

		var obj unstructured.Unstructured
		_, gvk, err := decUnstructured.Decode(documentBytes, nil, &obj)
		if err != nil {
			return err
		}

		mapping, err := k.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}

		var dr dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			dr = k.dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
		} else {
			dr = k.dyn.Resource(mapping.Resource)
		}

		_, err = dr.Update(ctx, &obj, metav1.UpdateOptions{
			FieldManager: "pulumi-k3s",
		})
		if apierrors.IsNotFound(err) {
			_, err = dr.Create(ctx, &obj, metav1.CreateOptions{
				FieldManager: "pulumi-k3s",
			})
		}
		if err != nil {
			return err
		}
	}

	return nil
}
