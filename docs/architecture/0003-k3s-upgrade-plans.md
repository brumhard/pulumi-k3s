# 3. K3s upgrade plans

Date: 2021-11-02

## Status

Accepted

## Context

One of the first features for the k3s provider should be automatic Kubernetes updates (by bumping the version defined in the resource).
There are multiple options to do so. For example, it's possible to [manually download](https://rancher.com/docs/k3s/latest/en/upgrades/basic/) a new version of k3s on the machine and executing that. This will stop the service, update it and restart it.
Another option is using the k3s [upgrade plan CRD](https://rancher.com/docs/k3s/latest/en/upgrades/automated/).

## Decision

For this provider the latter approach is used.
This has a couple of advantages.
First of all in a multi node setup it is not required to install and run the updated installation script on all nodes manually.
The only thing required is to set up the upgrade controller in the cluster and defining an upgrade plan a CRD with the version that should be installed.
If this version is bumped the controller will take care of updating all the nodes and restarting the k3s services.

To install this controller and the upgrade plans the provider will wait for the cluster to be up and running and then applies the required Kubernetes manifests.
While it would also be possible to put the manifests into the k3s [auto installation path](https://rancher.com/docs/k3s/latest/en/advanced/#auto-deploying-manifests) this has the advantage that the Pulumi cluster resource will only be ready after the Kubernetes cluster is reachable.

## Consequences

This is why k3s upgrade plans are used to manage version upgrades.
Another feature that is enabled by that are automatic updates using one of the various [k3s channels](https://rancher.com/docs/k3s/latest/en/upgrades/basic/#release-channels).
You just need to define a channel when creating the `cluster` resource in Pulumi, and it will automatically be updated as soon as a new version is released for the specified channel.
