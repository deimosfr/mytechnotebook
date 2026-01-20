---
title: "Cilium: Advanced CNI and network security for Kubernetes"
slug: cilium/
description: "A guide to run Cilium for Kubernetes"
categories: ["Network", "Kubernetes", "Server", "Home Lab"]
tags: ["Kubernetes", "Cilium", "CNI", "Network", "Home Lab"]
---

![Cilium](../../../static/images/cilium_logo.avif)

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

We can install Cilium easily with Helm:

```bash
helm repo add cilium https://helm.cilium.io/
helm install cilium cilium/cilium --version 1.18.4 -n kube-system --wait
```

## Configuration

Let's create a custom configuration to support:

1.  Load Balancing Layer 7 with Envoy
2.  Gateway API
3.  Hubble
4.  L2 announcements
5.  Hubble UI

=== "values-overrides.yaml"

    ```yaml
    terminationGracePeriodSeconds: 10

    # replace kube-proxy with cilium
    kubeProxyReplacement: "true"
    socketLB:
      enabled: true

    # Optional: Enable Prometheus metrics collection
    prometheus:
      enabled: true
      metricsService: true
      serviceMonitor:
        enabled: true
        trustCRDsExist: true

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
      # Optional: Enable if you have a Grafana Operator to deploy dashboards
      dashboards:
        enabled: true
      # Optional: Enable if you have a Prometheus Operator to scrape metrics
      prometheus:
        serviceMonitor:
          enabled: true

    # Gateway API
    ingressController:
      enabled: false
    gatewayAPI:
      enabled: true
      enableProxyProtocol: true

    # L7 support
    envoy:
      terminationGracePeriodSeconds: 10
      securityContext:
        capabilities:
          keepCapNetBindService: true
          envoy:
          - NET_ADMIN
          - SYS_ADMIN
          - NET_BIND_SERVICE
      # Optional: Enable if you have a Prometheus Operator to scrape metrics
      prometheus:
        serviceMonitor:
          enabled: true
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

We apply the configuration with:

```bash
helm install --upgrade -n kube-system --values values-overrides.yaml cilium cilium/cilium
```

### Layer 2 Announcements and IPAM

[L2 Announcements](https://docs.cilium.io/en/stable/network/l2-announcements/) makes services visible and reachable on the local area network. This is primarily for on-premises deployments without BGP, serving as a solid alternative to [MetalLB](./metallb_lb_k8s.md).

!!! warning

    L2 announcements are not stable yet. I encountered many various issues, IP unreachable, recover requiring a restart of cilium, etc... **DO NOT USE IN PRODUCTION YET**

When we use this feature, it responds to ARP queries for ExternalIPs and/or LoadBalancer IPs. These IPs are Virtual IPs (not installed on network devices) on multiple nodes. One node at a time will respond to the ARP queries with its MAC address, performing load balancing as a north/south load balancer.

The advantage over NodePort services is that each service can use a unique IP, allowing multiple services to use the same port numbers. With L2 announcements, if a node goes down, the service VIP simply migrates to another node and continues to work.

We can control the IP pool with [LoadBalancer IP Address Management (IPAM)](https://docs.cilium.io/en/stable/network/lb-ipam/#loadbalancer-ip-address-management-lb-ipam):

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

Then create a layer 2 configuration with the interfaces we want IPs to be announced on (from the hosts network interfaces):

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

In the Helm chart values, we need to have the following configuration (same as the `values-overrides.yaml` above):

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

Now if we deploy a service with as a load balancer, it will use the IP pool we defined above:

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

Once deployed, we can check the service with:

```bash
$ kubectl get svc my-service
NAME             TYPE           CLUSTER-IP      EXTERNAL-IP      PORT(S)          AGE
my-service       LoadBalancer   10.43.132.174   192.168.0.210    80:31043/TCP     15m
```

!!! note

    ICMP is not supported by Cilium L2 announcements. You have to use arpping to check if the IP is reachable or use netcat to check if the port is open.

### Network Policies

Cilium implements [Network Policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/) using [CiliumNetworkPolicy](https://docs.cilium.io/en/stable/security/policy/) resources. Cilium Network Policies are more permissive than Kubernetes Network Policies, as they can apply to all pods in a cluster, not just those in a namespace. It allows you to fine grain filter traffic to and from pods.

Here is an example on a Gateway API, to allow only local traffic to a service:

=== "gateway.yaml"

    ```yaml
    apiVersion: "cilium.io/v2"
    kind: CiliumNetworkPolicy
    metadata:
      name: allow-internal-to-gateway
      namespace: default
    spec:
      endpointSelector:
        matchLabels:
          # Match the internal gateway
          gateway.networking.k8s.io/gateway-name: internal-gateway
      ingress:
      - fromCIDR:
        # Your local network
        - 192.168.0.0/24
      egress:
      - toEntities:
        # Allow all outgoing traffic
        - all
    ```

!!! tip

    If you're not familiar with NetworkPolicies and need assistance, you can use the [Online Network Policy Editor](https://editor.networkpolicy.io/) to generate a policy based on your requirements. You can generate native Kubernetes `NetworkPolicy` or `CiliumNetworkPolicy` objects.

## Troubleshooting

### L2 Announcements

On the official site, you'll find [useful information](https://docs.cilium.io/en/stable/network/l2-announcements/#troubleshooting). It's a complete procedure.
