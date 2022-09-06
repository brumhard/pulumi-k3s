// Copyright 2016-2020, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/brumhard/pulumi-k3s/provider/pkg/k3s"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	structpb "github.com/golang/protobuf/ptypes/struct"
)

type k3sProvider struct {
	host    *provider.HostClient
	name    string
	version string
}

func makeProvider(host *provider.HostClient, name, version string) (pulumirpc.ResourceProviderServer, error) {
	// Return the new provider
	return &k3sProvider{
		host:    host,
		name:    name,
		version: version,
	}, nil
}

func (k *k3sProvider) Attach(ctx context.Context, req *pulumirpc.PluginAttach) (*emptypb.Empty, error) {
	// copy pasted implementation from boilerplate repo
	host, err := provider.NewHostClient(req.GetAddress())
	if err != nil {
		return nil, err
	}
	k.host = host
	return &pbempty.Empty{}, nil
}

// Check validates that the given property bag is valid for a resource of the given type and returns
// the inputs that should be passed to successive calls to Diff, Create, or Update for this
// resource. As a rule, the provider inputs returned by a call to Check should preserve the original
// representation of the properties as present in the program inputs. Though this rule is not
// required for correctness, violations thereof can negatively impact the end-user experience, as
// the provider inputs are using for detecting and rendering diffs.
func (k *k3sProvider) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "k3s:index:Cluster" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}
	// TODO: add validation here?
	return &pulumirpc.CheckResponse{Inputs: req.News, Failures: nil}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (k *k3sProvider) Diff(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "k3s:index:Cluster" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}

	// Retrieve the old state.
	oldState, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, err
	}

	// Extract old inputs from the `__inputs` field of the old state.
	oldInputs := parseCheckpointObject(oldState)

	// Get new resource inputs. The user is submitting these as an update.
	newResInputs, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "diff failed because malformed resource inputs %v %v",
			oldState, newResInputs)
	}

	// Calculate the difference between old and new inputs.
	d := oldInputs.Diff(newResInputs, func(key resource.PropertyKey) bool {
		return strings.HasPrefix(string(key), "__")
	})

	if d == nil {
		return &pulumirpc.DiffResponse{
			Changes: pulumirpc.DiffResponse_DIFF_NONE,
		}, nil
	}

	changes := make([]string, 0, len(d.Updates)+len(d.Adds)+len(d.Deletes))
	for k := range d.Updates {
		changes = append(changes, string(k))
	}
	for k := range d.Adds {
		changes = append(changes, string(k))
	}
	for k := range d.Deletes {
		changes = append(changes, string(k))
	}

	changeType := pulumirpc.DiffResponse_DIFF_NONE
	if len(changes) > 0 {
		changeType = pulumirpc.DiffResponse_DIFF_SOME
	}

	return &pulumirpc.DiffResponse{
		Diffs:           changes,
		Changes:         changeType,
		HasDetailedDiff: false,
	}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.
func (k *k3sProvider) Create(ctx context.Context, req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "k3s:index:Cluster" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}

	name := urn.Name().String()

	cluster, err := propsToCluster(req.GetProperties())
	if err != nil {
		return nil, err
	}

	err = k3s.MakeOrUpdateCluster(name, cluster)
	if err != nil {
		return nil, err
	}

	// Read the inputs to persist them into state.
	newInputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		KeepUnknowns: true,
		KeepSecrets:  true,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "diff failed because malformed resource inputs")
	}

	outputProperties, err := plugin.MarshalProperties(
		checkpointObject(newInputs, cluster),
		plugin.MarshalOptions{KeepSecrets: true, SkipNulls: true},
	)
	if err != nil {
		return nil, err
	}

	return &pulumirpc.CreateResponse{
		Id:         name,
		Properties: outputProperties,
	}, nil
}

// Read the current live state associated with a resource.
func (k *k3sProvider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "k3s:index:Cluster" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}
	return nil, status.Error(codes.Unimplemented, "Read is not yet implemented for 'k3s:index:Cluster'")
}

// Update updates an existing resource with new values.
func (k *k3sProvider) Update(ctx context.Context, req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "k3s:index:Cluster" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}

	cluster, err := propsToCluster(req.GetNews())
	if err != nil {
		return nil, err
	}

	err = k3s.MakeOrUpdateCluster(urn.Name().String(), cluster)
	if err != nil {
		return nil, err
	}

	// Read the inputs to persist them into state.
	newInputs, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		KeepUnknowns: true,
		KeepSecrets:  true,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "diff failed because malformed resource inputs")
	}

	outputProperties, err := plugin.MarshalProperties(
		checkpointObject(newInputs, cluster),
		plugin.MarshalOptions{KeepSecrets: true, SkipNulls: true},
	)
	if err != nil {
		return nil, err
	}

	return &pulumirpc.UpdateResponse{
		Properties: outputProperties,
	}, nil
}

func propsToCluster(props *structpb.Struct) (*k3s.Cluster, error) {
	propMap, err := plugin.UnmarshalProperties(props, plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	propBytes, err := json.Marshal(propMap.Mappable())
	if err != nil {
		return nil, err
	}

	var cluster k3s.Cluster
	if err := json.Unmarshal(propBytes, &cluster); err != nil {
		return nil, err
	}

	return &cluster, err
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed
// to still exist.
func (k *k3sProvider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "k3s:index:Cluster" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}

	cluster, err := propsToCluster(req.GetProperties())
	if err != nil {
		return nil, err
	}

	err = k3s.DeleteCluster(cluster)
	if err != nil {
		return nil, err
	}

	// Note that for our Random resource, we don't have to do anything on Delete.
	return &pbempty.Empty{}, nil
}

// Call dynamically executes a method in the provider associated with a component resource.
func (k *k3sProvider) Call(ctx context.Context, req *pulumirpc.CallRequest) (*pulumirpc.CallResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Call is not yet implemented")
}

// Construct creates a new component resource.
func (k *k3sProvider) Construct(ctx context.Context, req *pulumirpc.ConstructRequest) (*pulumirpc.ConstructResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Construct is not yet implemented")
}

// CheckConfig validates the configuration for this provider.
func (k *k3sProvider) CheckConfig(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return &pulumirpc.CheckResponse{Inputs: req.GetNews()}, nil
}

// DiffConfig diffs the configuration for this provider.
func (k *k3sProvider) DiffConfig(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	return &pulumirpc.DiffResponse{}, nil
}

// Configure configures the resource provider with "globals" that control its behavior.
func (k *k3sProvider) Configure(_ context.Context, req *pulumirpc.ConfigureRequest) (*pulumirpc.ConfigureResponse, error) {
	return &pulumirpc.ConfigureResponse{}, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (k *k3sProvider) Invoke(_ context.Context, req *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	tok := req.GetTok()
	return nil, fmt.Errorf("Unknown Invoke token '%s'", tok)
}

// StreamInvoke dynamically executes a built-in function in the provider. The result is streamed
// back as a series of messages.
func (k *k3sProvider) StreamInvoke(req *pulumirpc.InvokeRequest, server pulumirpc.ResourceProvider_StreamInvokeServer) error {
	tok := req.GetTok()
	return fmt.Errorf("Unknown StreamInvoke token '%s'", tok)
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (k *k3sProvider) GetPluginInfo(context.Context, *pbempty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{
		Version: k.version,
	}, nil
}

// GetSchema returns the JSON-serialized schema for the provider.
func (k *k3sProvider) GetSchema(ctx context.Context, req *pulumirpc.GetSchemaRequest) (*pulumirpc.GetSchemaResponse, error) {
	return &pulumirpc.GetSchemaResponse{}, nil
}

// Cancel signals the provider to gracefully shut down and abort any ongoing resource operations.
// Operations aborted in this way will return an error (e.g., `Update` and `Create` will either a
// creation error or an initialization error). Since Cancel is advisory and non-blocking, it is up
// to the host to decide how long to wait after Cancel is called before (e.g.)
// hard-closing any gRPC connection.
func (k *k3sProvider) Cancel(context.Context, *pbempty.Empty) (*pbempty.Empty, error) {
	// TODO
	return &pbempty.Empty{}, nil
}

// checkpointObject puts inputs in the `__inputs` field of the state.
func checkpointObject(inputs resource.PropertyMap, outputs interface{}) resource.PropertyMap {
	object := resource.NewPropertyMap(outputs)
	object["__inputs"] = resource.MakeSecret(resource.NewObjectProperty(inputs))
	return object
}

// parseCheckpointObject returns inputs that are saved in the `__inputs` field of the state.
func parseCheckpointObject(obj resource.PropertyMap) resource.PropertyMap {
	if inputs, ok := obj["__inputs"]; ok {
		return inputs.ObjectValue()
	}

	return nil
}
