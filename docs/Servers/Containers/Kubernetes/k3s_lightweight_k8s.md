---
title: "K3s: A lightweight Kubernetes"
slug: k3s-lightweight-k8s/
description: "A guide to run a lightweight Kubernetes ideal for home lab."
categories: ["Kubernetes", "Server", "Home Lab"]
tags: ["Kubernetes", "K3s", "Home Lab"]
---

## Introduction

[K3s](https://k3s.io/) is a lightweight Kubernetes distribution that is designed to be easy to install and use. It is a great way to run a Kubernetes cluster in a home lab.

In this documentation, I assume you're running a Linux server on [Debian](https://www.debian.org/).

## Requirements

You can find a list of requirements on the [K3s website](https://docs.k3s.io/installation/requirements). Basically the most important thing is to disable the firewall system to not overlap with the k3s network.

## Installation

To install k3s on multiple nodes, its is preferable to use a configuration management tool like [Ansible](../../Configuration%20Managers/index.md).

Here we'll see how to install k3s on a single node manually for simplicity. K3s has 2 roles:

- Server (K3s Control plane)
- Agent (K3s client/workers)

On large cluster, you generally have 3 (or more) dedicated nodes for the control plane and the rest for the workers. But here we'll keep it simple and run the control plane and the worker on the same node.

Start by creating a file to configure the kubelet:

=== "kubelet.config"

    ```yaml
    apiVersion: kubelet.config.k8s.io/v1beta1
    kind: KubeletConfiguration
    shutdownGracePeriod: 180s
    shutdownGracePeriodCriticalPods: 60s
    failSwapOn: false
    featureGates:
        NodeSwap: true
    memorySwap:
        swapBehavior: LimitedSwap
    ```

Then run this command with root privileges:

```
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--kubelet-arg 'config=/etc/rancher/k3s/kubelet.config' --etcd-expose-metrics --flannel-backend=none --disable-network-policy --disable=traefik --disable=metrics-server --disable servicelb --bind-address=0.0.0.0" sh -
```

- `etcd-expose-metrics` is used to expose the etcd metrics to the Prometheus server.
- `flannel-backend=none` is used to disable the flannel network plugin because I prefer using [Cilium](https://docs.cilium.io/en/stable/) for the network policy.
- `disable-network-policy` required by Cilium.
- `disable=traefik` you can keep it, but I prefer letting [Qovery](https://qovery.com/) handle the ingress with Nginx.
- `disable=metrics-server` same here, I prefer using [Qovery](https://qovery.com/) for the metrics.
- `disable servicelb` we'll use [metallb](../../Network/metallb_lb_k8s.md) for the load balancer.
- `bind-address=0.0.0.0` is used to bind the kubelet to all interfaces.

## Graceful shutdown

I personaly don't like how k3s shutdown the server when a reboot is triggered. It's not graceful and dangerous if you're hosting stateful apps likes databases.

To gracefully shutdown the k3s server, you can use this script:

=== "/usr/local/bin/k3s-node-drain.sh"

    ```bash
    #!/bin/bash
    set -e

    # Get the node name
    NODE_NAME=$(hostname)

    # Log the start of drain process
    echo "Starting drain of node ${NODE_NAME} before reboot" | systemd-cat -t k3s-drain

    # Attempt to drain the node
    if kubectl drain ${NODE_NAME} --ignore-daemonsets --delete-emptydir-data --timeout=300s --grace-period=120; then
        echo "Successfully drained node ${NODE_NAME}" | systemd-cat -t k3s-drain
        exit 0
    else
        echo "Failed to drain node ${NODE_NAME}, but continuing with reboot" | systemd-cat -t k3s-drain
        # We still exit with 0 to allow the reboot to proceed
        exit 0
    fi
    ```

Then create a systemd service to run it on reboot/shutdown:

=== "/etc/systemd/system/k3s-node-drain.service"

    ```ini
    [Unit]
    Description=Drain K3s node before shutdown
    DefaultDependencies=no
    # Run before k3s stops and before shutdown sequence
    Before=k3s.service shutdown.target reboot.target halt.target
    # Ensure k3s is still running when we try to drain
    After=network-online.target
    Requires=k3s.service
    # Only run on shutdown
    Conflicts=shutdown.target reboot.target halt.target

    [Service]
    Type=oneshot
    ExecStart=/usr/local/bin/k3s-node-drain.sh
    TimeoutStartSec=300
    # Continue even if the script fails, we don't want to block a reboot
    SuccessExitStatus=0 1

    [Install]
    WantedBy=shutdown.target reboot.target halt.target
    ```

Finally set execution permissions and enable the service:

```bash
chmod +x /usr/local/bin/k3s-node-drain.sh
systemctl enable k3s-node-drain.service
systemctl daemon-reload
```

You can now reboot the server and see that pods are gracefully shutdown:

```bash
kubectl get po --watch
```
