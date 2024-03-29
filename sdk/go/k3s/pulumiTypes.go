// Code generated by pulumigen DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package k3s

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CRIConfiguration struct {
	EnableGVisor *bool `pulumi:"enableGVisor"`
}

// CRIConfigurationInput is an input type that accepts CRIConfigurationArgs and CRIConfigurationOutput values.
// You can construct a concrete instance of `CRIConfigurationInput` via:
//
//	CRIConfigurationArgs{...}
type CRIConfigurationInput interface {
	pulumi.Input

	ToCRIConfigurationOutput() CRIConfigurationOutput
	ToCRIConfigurationOutputWithContext(context.Context) CRIConfigurationOutput
}

type CRIConfigurationArgs struct {
	EnableGVisor pulumi.BoolPtrInput `pulumi:"enableGVisor"`
}

func (CRIConfigurationArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*CRIConfiguration)(nil)).Elem()
}

func (i CRIConfigurationArgs) ToCRIConfigurationOutput() CRIConfigurationOutput {
	return i.ToCRIConfigurationOutputWithContext(context.Background())
}

func (i CRIConfigurationArgs) ToCRIConfigurationOutputWithContext(ctx context.Context) CRIConfigurationOutput {
	return pulumi.ToOutputWithContext(ctx, i).(CRIConfigurationOutput)
}

func (i CRIConfigurationArgs) ToCRIConfigurationPtrOutput() CRIConfigurationPtrOutput {
	return i.ToCRIConfigurationPtrOutputWithContext(context.Background())
}

func (i CRIConfigurationArgs) ToCRIConfigurationPtrOutputWithContext(ctx context.Context) CRIConfigurationPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(CRIConfigurationOutput).ToCRIConfigurationPtrOutputWithContext(ctx)
}

// CRIConfigurationPtrInput is an input type that accepts CRIConfigurationArgs, CRIConfigurationPtr and CRIConfigurationPtrOutput values.
// You can construct a concrete instance of `CRIConfigurationPtrInput` via:
//
//	        CRIConfigurationArgs{...}
//
//	or:
//
//	        nil
type CRIConfigurationPtrInput interface {
	pulumi.Input

	ToCRIConfigurationPtrOutput() CRIConfigurationPtrOutput
	ToCRIConfigurationPtrOutputWithContext(context.Context) CRIConfigurationPtrOutput
}

type criconfigurationPtrType CRIConfigurationArgs

func CRIConfigurationPtr(v *CRIConfigurationArgs) CRIConfigurationPtrInput {
	return (*criconfigurationPtrType)(v)
}

func (*criconfigurationPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**CRIConfiguration)(nil)).Elem()
}

func (i *criconfigurationPtrType) ToCRIConfigurationPtrOutput() CRIConfigurationPtrOutput {
	return i.ToCRIConfigurationPtrOutputWithContext(context.Background())
}

func (i *criconfigurationPtrType) ToCRIConfigurationPtrOutputWithContext(ctx context.Context) CRIConfigurationPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(CRIConfigurationPtrOutput)
}

type CRIConfigurationOutput struct{ *pulumi.OutputState }

func (CRIConfigurationOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*CRIConfiguration)(nil)).Elem()
}

func (o CRIConfigurationOutput) ToCRIConfigurationOutput() CRIConfigurationOutput {
	return o
}

func (o CRIConfigurationOutput) ToCRIConfigurationOutputWithContext(ctx context.Context) CRIConfigurationOutput {
	return o
}

func (o CRIConfigurationOutput) ToCRIConfigurationPtrOutput() CRIConfigurationPtrOutput {
	return o.ToCRIConfigurationPtrOutputWithContext(context.Background())
}

func (o CRIConfigurationOutput) ToCRIConfigurationPtrOutputWithContext(ctx context.Context) CRIConfigurationPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v CRIConfiguration) *CRIConfiguration {
		return &v
	}).(CRIConfigurationPtrOutput)
}

func (o CRIConfigurationOutput) EnableGVisor() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v CRIConfiguration) *bool { return v.EnableGVisor }).(pulumi.BoolPtrOutput)
}

type CRIConfigurationPtrOutput struct{ *pulumi.OutputState }

func (CRIConfigurationPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**CRIConfiguration)(nil)).Elem()
}

func (o CRIConfigurationPtrOutput) ToCRIConfigurationPtrOutput() CRIConfigurationPtrOutput {
	return o
}

func (o CRIConfigurationPtrOutput) ToCRIConfigurationPtrOutputWithContext(ctx context.Context) CRIConfigurationPtrOutput {
	return o
}

func (o CRIConfigurationPtrOutput) Elem() CRIConfigurationOutput {
	return o.ApplyT(func(v *CRIConfiguration) CRIConfiguration {
		if v != nil {
			return *v
		}
		var ret CRIConfiguration
		return ret
	}).(CRIConfigurationOutput)
}

func (o CRIConfigurationPtrOutput) EnableGVisor() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *CRIConfiguration) *bool {
		if v == nil {
			return nil
		}
		return v.EnableGVisor
	}).(pulumi.BoolPtrOutput)
}

type Node struct {
	Args       []string          `pulumi:"args"`
	CriConfig  *CRIConfiguration `pulumi:"criConfig"`
	Host       string            `pulumi:"host"`
	PrivateKey string            `pulumi:"privateKey"`
	User       *string           `pulumi:"user"`
}

// Defaults sets the appropriate defaults for Node
func (val *Node) Defaults() *Node {
	if val == nil {
		return nil
	}
	tmp := *val
	if isZero(tmp.User) {
		user_ := "root"
		tmp.User = &user_
	}
	return &tmp
}

// NodeInput is an input type that accepts NodeArgs and NodeOutput values.
// You can construct a concrete instance of `NodeInput` via:
//
//	NodeArgs{...}
type NodeInput interface {
	pulumi.Input

	ToNodeOutput() NodeOutput
	ToNodeOutputWithContext(context.Context) NodeOutput
}

type NodeArgs struct {
	Args       pulumi.StringArrayInput  `pulumi:"args"`
	CriConfig  CRIConfigurationPtrInput `pulumi:"criConfig"`
	Host       pulumi.StringInput       `pulumi:"host"`
	PrivateKey pulumi.StringInput       `pulumi:"privateKey"`
	User       pulumi.StringPtrInput    `pulumi:"user"`
}

// Defaults sets the appropriate defaults for NodeArgs
func (val *NodeArgs) Defaults() *NodeArgs {
	if val == nil {
		return nil
	}
	tmp := *val
	if isZero(tmp.User) {
		tmp.User = pulumi.StringPtr("root")
	}
	return &tmp
}
func (NodeArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*Node)(nil)).Elem()
}

func (i NodeArgs) ToNodeOutput() NodeOutput {
	return i.ToNodeOutputWithContext(context.Background())
}

func (i NodeArgs) ToNodeOutputWithContext(ctx context.Context) NodeOutput {
	return pulumi.ToOutputWithContext(ctx, i).(NodeOutput)
}

// NodeArrayInput is an input type that accepts NodeArray and NodeArrayOutput values.
// You can construct a concrete instance of `NodeArrayInput` via:
//
//	NodeArray{ NodeArgs{...} }
type NodeArrayInput interface {
	pulumi.Input

	ToNodeArrayOutput() NodeArrayOutput
	ToNodeArrayOutputWithContext(context.Context) NodeArrayOutput
}

type NodeArray []NodeInput

func (NodeArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]Node)(nil)).Elem()
}

func (i NodeArray) ToNodeArrayOutput() NodeArrayOutput {
	return i.ToNodeArrayOutputWithContext(context.Background())
}

func (i NodeArray) ToNodeArrayOutputWithContext(ctx context.Context) NodeArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(NodeArrayOutput)
}

type NodeOutput struct{ *pulumi.OutputState }

func (NodeOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Node)(nil)).Elem()
}

func (o NodeOutput) ToNodeOutput() NodeOutput {
	return o
}

func (o NodeOutput) ToNodeOutputWithContext(ctx context.Context) NodeOutput {
	return o
}

func (o NodeOutput) Args() pulumi.StringArrayOutput {
	return o.ApplyT(func(v Node) []string { return v.Args }).(pulumi.StringArrayOutput)
}

func (o NodeOutput) CriConfig() CRIConfigurationPtrOutput {
	return o.ApplyT(func(v Node) *CRIConfiguration { return v.CriConfig }).(CRIConfigurationPtrOutput)
}

func (o NodeOutput) Host() pulumi.StringOutput {
	return o.ApplyT(func(v Node) string { return v.Host }).(pulumi.StringOutput)
}

func (o NodeOutput) PrivateKey() pulumi.StringOutput {
	return o.ApplyT(func(v Node) string { return v.PrivateKey }).(pulumi.StringOutput)
}

func (o NodeOutput) User() pulumi.StringPtrOutput {
	return o.ApplyT(func(v Node) *string { return v.User }).(pulumi.StringPtrOutput)
}

type NodeArrayOutput struct{ *pulumi.OutputState }

func (NodeArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]Node)(nil)).Elem()
}

func (o NodeArrayOutput) ToNodeArrayOutput() NodeArrayOutput {
	return o
}

func (o NodeArrayOutput) ToNodeArrayOutputWithContext(ctx context.Context) NodeArrayOutput {
	return o
}

func (o NodeArrayOutput) Index(i pulumi.IntInput) NodeOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) Node {
		return vs[0].([]Node)[vs[1].(int)]
	}).(NodeOutput)
}

type VersionConfiguration struct {
	Channel *string `pulumi:"channel"`
	Version *string `pulumi:"version"`
}

// VersionConfigurationInput is an input type that accepts VersionConfigurationArgs and VersionConfigurationOutput values.
// You can construct a concrete instance of `VersionConfigurationInput` via:
//
//	VersionConfigurationArgs{...}
type VersionConfigurationInput interface {
	pulumi.Input

	ToVersionConfigurationOutput() VersionConfigurationOutput
	ToVersionConfigurationOutputWithContext(context.Context) VersionConfigurationOutput
}

type VersionConfigurationArgs struct {
	Channel pulumi.StringPtrInput `pulumi:"channel"`
	Version pulumi.StringPtrInput `pulumi:"version"`
}

func (VersionConfigurationArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*VersionConfiguration)(nil)).Elem()
}

func (i VersionConfigurationArgs) ToVersionConfigurationOutput() VersionConfigurationOutput {
	return i.ToVersionConfigurationOutputWithContext(context.Background())
}

func (i VersionConfigurationArgs) ToVersionConfigurationOutputWithContext(ctx context.Context) VersionConfigurationOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VersionConfigurationOutput)
}

func (i VersionConfigurationArgs) ToVersionConfigurationPtrOutput() VersionConfigurationPtrOutput {
	return i.ToVersionConfigurationPtrOutputWithContext(context.Background())
}

func (i VersionConfigurationArgs) ToVersionConfigurationPtrOutputWithContext(ctx context.Context) VersionConfigurationPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VersionConfigurationOutput).ToVersionConfigurationPtrOutputWithContext(ctx)
}

// VersionConfigurationPtrInput is an input type that accepts VersionConfigurationArgs, VersionConfigurationPtr and VersionConfigurationPtrOutput values.
// You can construct a concrete instance of `VersionConfigurationPtrInput` via:
//
//	        VersionConfigurationArgs{...}
//
//	or:
//
//	        nil
type VersionConfigurationPtrInput interface {
	pulumi.Input

	ToVersionConfigurationPtrOutput() VersionConfigurationPtrOutput
	ToVersionConfigurationPtrOutputWithContext(context.Context) VersionConfigurationPtrOutput
}

type versionConfigurationPtrType VersionConfigurationArgs

func VersionConfigurationPtr(v *VersionConfigurationArgs) VersionConfigurationPtrInput {
	return (*versionConfigurationPtrType)(v)
}

func (*versionConfigurationPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**VersionConfiguration)(nil)).Elem()
}

func (i *versionConfigurationPtrType) ToVersionConfigurationPtrOutput() VersionConfigurationPtrOutput {
	return i.ToVersionConfigurationPtrOutputWithContext(context.Background())
}

func (i *versionConfigurationPtrType) ToVersionConfigurationPtrOutputWithContext(ctx context.Context) VersionConfigurationPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VersionConfigurationPtrOutput)
}

type VersionConfigurationOutput struct{ *pulumi.OutputState }

func (VersionConfigurationOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*VersionConfiguration)(nil)).Elem()
}

func (o VersionConfigurationOutput) ToVersionConfigurationOutput() VersionConfigurationOutput {
	return o
}

func (o VersionConfigurationOutput) ToVersionConfigurationOutputWithContext(ctx context.Context) VersionConfigurationOutput {
	return o
}

func (o VersionConfigurationOutput) ToVersionConfigurationPtrOutput() VersionConfigurationPtrOutput {
	return o.ToVersionConfigurationPtrOutputWithContext(context.Background())
}

func (o VersionConfigurationOutput) ToVersionConfigurationPtrOutputWithContext(ctx context.Context) VersionConfigurationPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v VersionConfiguration) *VersionConfiguration {
		return &v
	}).(VersionConfigurationPtrOutput)
}

func (o VersionConfigurationOutput) Channel() pulumi.StringPtrOutput {
	return o.ApplyT(func(v VersionConfiguration) *string { return v.Channel }).(pulumi.StringPtrOutput)
}

func (o VersionConfigurationOutput) Version() pulumi.StringPtrOutput {
	return o.ApplyT(func(v VersionConfiguration) *string { return v.Version }).(pulumi.StringPtrOutput)
}

type VersionConfigurationPtrOutput struct{ *pulumi.OutputState }

func (VersionConfigurationPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**VersionConfiguration)(nil)).Elem()
}

func (o VersionConfigurationPtrOutput) ToVersionConfigurationPtrOutput() VersionConfigurationPtrOutput {
	return o
}

func (o VersionConfigurationPtrOutput) ToVersionConfigurationPtrOutputWithContext(ctx context.Context) VersionConfigurationPtrOutput {
	return o
}

func (o VersionConfigurationPtrOutput) Elem() VersionConfigurationOutput {
	return o.ApplyT(func(v *VersionConfiguration) VersionConfiguration {
		if v != nil {
			return *v
		}
		var ret VersionConfiguration
		return ret
	}).(VersionConfigurationOutput)
}

func (o VersionConfigurationPtrOutput) Channel() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *VersionConfiguration) *string {
		if v == nil {
			return nil
		}
		return v.Channel
	}).(pulumi.StringPtrOutput)
}

func (o VersionConfigurationPtrOutput) Version() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *VersionConfiguration) *string {
		if v == nil {
			return nil
		}
		return v.Version
	}).(pulumi.StringPtrOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*CRIConfigurationInput)(nil)).Elem(), CRIConfigurationArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*CRIConfigurationPtrInput)(nil)).Elem(), CRIConfigurationArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*NodeInput)(nil)).Elem(), NodeArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*NodeArrayInput)(nil)).Elem(), NodeArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*VersionConfigurationInput)(nil)).Elem(), VersionConfigurationArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*VersionConfigurationPtrInput)(nil)).Elem(), VersionConfigurationArgs{})
	pulumi.RegisterOutputType(CRIConfigurationOutput{})
	pulumi.RegisterOutputType(CRIConfigurationPtrOutput{})
	pulumi.RegisterOutputType(NodeOutput{})
	pulumi.RegisterOutputType(NodeArrayOutput{})
	pulumi.RegisterOutputType(VersionConfigurationOutput{})
	pulumi.RegisterOutputType(VersionConfigurationPtrOutput{})
}
