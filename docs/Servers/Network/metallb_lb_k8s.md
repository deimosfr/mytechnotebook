---
title: "MetalLB: Self-hosted Load Balancers for Kubernetes"
slug: metallb-lb-k8s/
description: "A guide to run a self-hosted load balancers for Kubernetes using MetalLB."
categories: ["Network", "Kubernetes", "Server", "Home Lab"]
tags: ["Kubernetes", "MetalLB", "Load Balancer", "Home Lab"]
---

## Introduction

[MetalLB](https://metallb.io/) hooks into your Kubernetes cluster, and provides a network load-balancer implementation. In short, it allows you to create Kubernetes services of type LoadBalancer in clusters that don’t run on a cloud provider, and thus cannot simply hook into paid products to provide load balancers.

Kubernetes does not offer an implementation of network load balancers (Services of type LoadBalancer) for bare-metal clusters. The implementations of network load balancers that Kubernetes does ship with are all glue code that calls out to various IaaS platforms (GCP, AWS, Azure…). If you’re not running on a supported IaaS platform (GCP, AWS, Azure…), LoadBalancers will remain in the “pending” state indefinitely when created.

Bare-metal cluster operators are left with two lesser tools to bring user traffic into their clusters, “NodePort” and “externalIPs” services. Both of these options have significant downsides for production use, which makes bare-metal clusters second-class citizens in the Kubernetes ecosystem.

MetalLB aims to redress this imbalance by offering a network load balancer implementation that integrates with standard network equipment, so that external services on bare-metal clusters also “just work” as much as possible.

## Installation

Installing MetalLB is pretty straightforward with Helm:

```bash
helm repo add metallb https://metallb.github.io/metallb
helm install metallb metallb/metallb -n kube-system --wait
```

## Configuration

Create a pool of IP addresses for the load balancers:

=== "ipaddresspool.yaml"

    ```yaml
    apiVersion: metallb.io/v1beta1
    kind: IPAddressPool
    metadata:
    name: metallb-pool
    namespace: kube-system
    spec:
    # 50 available IPs for the load balancers
    addresses:
        - 192.168.0.1-192.168.0.50
    ```

Then create a layer 2 configuration:

=== "l2advertisement.yaml"

    ```yaml
    apiVersion: metallb.io/v1beta1
    kind: L2Advertisement
    metadata:
    name: l2-advertisement
    namespace: kube-system
    spec:
    ipAddressPools:
        - metallb-pool
    ```

Then apply the configuration:

```bash
kubectl apply -f ipaddresspool.yaml -f l2advertisement.yaml
```

## Usage

Now when you want to ask for a load balancer, you need to use an annotation to tell MetalLB to use the load balancer. You need to set also the IP you want to use for the load balancer. Here is an example to add in a service:

```yaml
service:
  enabled: true
  externalTrafficPolicy: Local
  annotations:
    metallb.universe.tf/loadBalancerIPs: 192.168.0.1
```
