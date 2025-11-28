---
title: "Cilium: CNI and network configuration for Kubernetes"
slug: cilium/
description: "A guide to run Cilium for Kubernetes"
categories: ["Network", "Kubernetes", "Server", "Home Lab"]
tags: ["Kubernetes", "Cilium", "CNI", "Network", "Home Lab"]
---

## Introduction

[Cilium](https://cilium.io/) is an open-source project that provides networking, security, and observability for cloud-native environments such as Kubernetes clusters and other container orchestration platforms.

At the foundation of Cilium is a new Linux kernel technology called **eBPF**, which enables the dynamic insertion of powerful security visibility and control logic within Linux itself. Because eBPF runs inside the Linux kernel, Cilium security policies can be applied and updated without any changes to the application code or container configuration.

Cilium acts as a CNI (Container Network Interface) plugin for Kubernetes, but goes far beyond basic connectivity:

*   **Networking**: It provides a highly scalable networking plane, supporting multi-cluster connectivity (Cluster Mesh) and replacing `kube-proxy` for service load balancing.
*   **Security**: It implements identity-aware network policies (L3-L7) that are decoupled from network addressing, allowing for more flexible and secure communication controls.
*   **Observability**: Through its component **Hubble**, it offers deep visibility into network traffic, service dependencies, and operational metrics.

Cilium's performance advantage comes from its use of eBPF to bypass the legacy `iptables`-based networking stack traditionally used in Kubernetes.

*   **Efficient Lookups**: Unlike `iptables` which uses a linear list of rules (O(n) complexity), Cilium uses eBPF maps (hash tables) which provide O(1) lookups. This means performance remains stable even as the number of services and rules grows massively.
*   **Reduced Overhead**: eBPF programs run directly in the kernel, minimizing the overhead of context switching between user space and kernel space.
*   **Direct Routing**: Cilium can perform more efficient routing and load balancing decisions earlier in the packet processing path.

## Installation

Installing Cilium is pretty straightforward with Helm:

```bash
helm repo add cilium https://helm.cilium.io/
helm install cilium cilium/cilium --version 1.18.4 -n kube-system --wait
```

## Configuration

Here is a custom configuration I'm using to support:

1. Load Balancing Layer 7 with envoy
2. Gateway API
3. Hubble
4. L2 announcements
5. Hubble UI

=== "values-overrides.yaml"

    ```yaml
    terminationGracePeriodSeconds: 10

    # replace kube-proxy with cilium
    kubeProxyReplacement: "true"
    socketLB:
      enabled: true

    # Operator
    operator:
      replicas: 2
      podDisruptionBudget:
        enabled: true
        minAvailable: 1
        maxUnavailable: null
      endpointGCInterval: "2m0s"
      nodeGCInterval: "2m0s"
      identityGCInterval: "10m0s"
      setNodeTaints: true

    # Gateway API
    ingressController:
      enabled: false
    gatewayAPI:
      enabled: true
      enableProxyProtocol: true

    # L7 support
    envoy:
      terminationGracePeriodSeconds: 10
    l7Proxy: true
    l7:
      backend: envoy
      algorithm: least_request
      acceleration: native

    # L2 announcements
    K8sClientRateLimit:
      qps: 20
      burst: 40

    lbExternalClusterIP: true
    l2announcements:
      enabled: true
      leaseDuration: 10s
      leaseRenewDeadline: 5s
      leaseRetryPeriod: 1s

    # BPF
    bpf:
      masquerade: true

    # Hubble
    hubble:
      relay:
        enabled: true
      ui:
        enabled: true
    ```

You can apply the configuration with:

```bash
helm install --upgrade -n kube-system --values values-overrides.yaml cilium cilium/cilium
```

### Layer 2 Announcements and IPAM

[L2 Announcements](https://docs.cilium.io/en/stable/network/l2-announcements/) is a feature which makes services visible and reachable on the local area network. This feature is primarily intended for on-premises deployments within networks without BGP based routing such as office or campus networks. It's a good alternative to [MetalLB](./metallb_lb_k8s.md).

When used, this feature will respond to ARP queries for ExternalIPs and/or LoadBalancer IPs. These IPs are Virtual IPs (not installed on network devices) on multiple nodes, so for each service one node at a time will respond to the ARP queries and respond with its MAC address. This node will perform load balancing with the service load balancing feature, thus acting as a north/south load balancer.

The advantage of this feature over NodePort services is that each service can use a unique IP so multiple services can use the same port numbers. When using NodePorts, it is up to the client to decide to which host to send traffic, and if a node goes down, the IP+Port combo becomes unusable. With L2 announcements the service VIP simply migrates to another node and will continue to work.

You can control the IP pool with [LoadBalancer IP Address Management (IPAM)](https://docs.cilium.io/en/stable/network/lb-ipam/#loadbalancer-ip-address-management-lb-ipam):

=== "ippool.yaml"

    ```yaml
    apiVersion: "cilium.io/v2"
    kind: CiliumLoadBalancerIPPool
    metadata:
      name: "cilium-lb-pool"
    spec:
      blocks:
      - start: "192.168.0.201"
        stop: "192.168.0.240"
    ```

Then create a layer 2 configuration with the interfaces you want IPs to be announced on (from the hosts network interfaces):

=== "l2advertisement.yaml"

    ```yaml
    apiVersion: "cilium.io/v2alpha1"
    kind: CiliumL2AnnouncementPolicy
    metadata:
      name: default-l2-announcement-policy
      namespace: kube-system
    spec:
      externalIPs: true
      loadBalancerIPs: true
      interfaces:
      - ^bond[0-9]+$
      - ^enp0s[0-9]+$
      - ^ens[0-9]+$
    ```

In the Helm chart values, you need to have the following configuration(same as the `values-overrides.yaml` above):

=== "values-overrides.yaml"

    ```yaml
    lbExternalClusterIP: true
    l2announcements:
      enabled: true
      leaseDuration: 10s
      leaseRenewDeadline: 5s
      leaseRetryPeriod: 1s
    ```

Then apply the configuration:

```bash
kubectl apply -n kube-system -f ippool.yaml -f l2advertisement.yaml
```

Now if you deploy a service with as a load balancer, it will use the IP pool you defined above:

=== "service.yaml"

    ```yaml
    apiVersion: v1
    kind: Service
    metadata:
      name: my-service 
      annotations:
        # Optional: Request a specific IP from the pool
        lbipam.cilium.io/ips: 192.168.0.210
    spec:
      type: LoadBalancer
      externalTrafficPolicy: Cluster
      internalTrafficPolicy: Cluster
      ports:
      - name: http
        port: 80
        protocol: TCP
        targetPort: http
      selector:
        app: my-app
    ```

Once deployed, you can check the service with:

```bash
$ kubectl get svc my-service
NAME             TYPE           CLUSTER-IP      EXTERNAL-IP      PORT(S)          AGE
my-service       LoadBalancer   10.43.132.174   192.168.0.210    80:31043/TCP     15m
```

!!! note

    ICMP is not supported by Cilium L2 announcements. You have to use arpping to check if the IP is reachable or use netcat to check if the port is open.

## Troubleshooting

### L2 Announcements

First, on the official site, you'll find [useful information](https://docs.cilium.io/en/stable/network/l2-announcements/#troubleshooting).
