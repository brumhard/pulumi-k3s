# k3s Pulumi Provider

> This is a side project and not intended for production usage.
> That's why currently only single node setups are supported and there are no released binaries.

`pulumi-k3s` closes the gap between existing Pulumi providers for IaaS (VMs) and the Kubernetes provider.
With `pulumi-k3s` it is possible to define all your required infrastructure as well as Kubernetes resources in one pulumi stack
with proper dependency tracking while provisioning your whole landscape.

An example of using this provider's `Cluster` resource can be found in the [examples](examples/simple/index.ts).

To get a quick preview, here's a minimal example:

```typescript
const cluster = new k3s.Cluster("mycluster", {
    // currently only one master node is supported in the following
    // list. To not break compatibility with future releases it's already
    // designed to have multiple master nodes.
    masterNodes: [
        {
            host: cfg.require("master_ip"),
            user: "ubuntu",
            privateKey: cfg.requireSecret("private_key"),
        }
    ],
    versionConfig: {
        version: "v1.22.3+k3s1"
    }
});

// cluster.kubeconfig is the output that contains the kubeconfig used
// for the Pulumi Kubernetes provider.
...
```

## Installation

Currently, this provider is not optimized for distribution. Therefore, if you want to use it, you need to follow the instructions defined in the [Development section](#development) and build the provider as well as the needed SDK for you preferred language from the main branch and use that.

## Development

Most of the code for the provider implementation is in `pkg/provider/provider.go`.  

A code generator is available which generates SDKs in TypeScript, Python, Go and .NET which are also checked in to the `sdk` folder.  The SDKs are generated from a schema in `provider/cmd/pulumi-resource-k3s/schema.json`.  This file should be kept aligned with the resources, functions and types supported by the provider implementation.

Note that the generated provider plugin (`pulumi-resource-k3s`) must be on your `PATH` to be used by Pulumi deployments. If creating a provider for distribution to other users, you should ensure they install this plugin to their `PATH`.

### Pre-requisites

Install the `pulumictl` CLI from the [releases](https://github.com/pulumi/pulumictl/releases) page or follow the [installation instructions](https://github.com/pulumi/pulumictl#installation).

> NB: Usage of `pulumictl` is optional. If not using it, hard code the version in the [Makefile](Makefile) of when building explicitly pass version as `VERSION=0.0.1 make build`

### Build and Test

```bash
# build and install the resource provider plugin
$ make build install

# test
$ cd examples/simple
$ yarn link @pulumi/k3s
$ yarn install
$ pulumi stack init test
$ pulumi up
```