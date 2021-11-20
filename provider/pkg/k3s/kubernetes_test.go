package k3s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/meta/testrestmapper"
	fakedynamic "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/scheme"
)

func Test_K8sClient_apply(t *testing.T) {
	k8sClient := &K8sClient{
		dyn:    fakedynamic.NewSimpleDynamicClient(scheme.Scheme),
		mapper: testrestmapper.TestOnlyStaticRESTMapper(scheme.Scheme),
	}

	for i := 0; i < 2; i++ {
		err := k8sClient.CreateOrUpdateFromFile(context.Background(), []byte(`---
apiVersion: v1
kind: Namespace
metadata:
    name: system-upgrade
---
apiVersion: v1
kind: Namespace
metadata:
    name: system-upgrade2`))
		assert.NoError(t, err)
	}
}
