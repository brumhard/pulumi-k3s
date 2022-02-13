# 2. Build a Pulumi provider

Date: 2021-11-02

## Status

Accepted

## Context

To use the Pulumi provider for Kubernetes obviously a Kubernetes cluster is required.
One of the easiest ways to set one up for personal use is the creating some VMs in any cloud provider
and installing k3s on it.
To make that even easier tools like [k3sup](https://github.com/alexellis/k3sup) have been created by the community to install a cluster really easily in less than a minute.

To support that setup in Pulumi multiple steps would be required:

1. Create VMs with one of the [Pulumi providers](https://www.pulumi.com/registry/) like the Azure provider.
2. Apply the first VM stack and set the resulting public IP and credentials as outputs.
3. Use the outputs to connect to the VM and install k3s using a tool like k3sup or a configuration tool like Ansible.
4. Save the resulting kubeconfig.
5. Use the kubeconfig to configure the [Pulumi Kubernetes Provider](https://www.pulumi.com/registry/packages/kubernetes/) in a second stack and apply it.
6. Provision Kubernetes resources.

## Decision

This process can be simplified by creating a Pulumi provider. Like this the steps are:

1. Set up a single stack with VMs from any cloud provider.
2. Use the public IP and credentials as inputs for the k3s Pulumi provider.
3. Use the kubeconfig output property of the k3s `cluster` resource as input for the Kubernetes provider.
4. Apply the whole landscape in one apply.

To create a provider for Pulumi there are different options like native providers, dynamic providers or first writing a Terraform provider and translating that to pulumi afterwards.

For this specific project the full power of Pulumi should be leveraged. Also, it should be possible to define actions for all CRUD operations on the `cluster` resource. With dynamic provider it would be possible to install k3s on a single VM and configure it, but for more complex it would be hard to maintain all the logic needed.

A Terraform provider could be an option, but this project should also be used to deep dive into how Pulumi works.

## Consequences

Therefore, for this project a native Pulumi provider will be developed from scratch.
