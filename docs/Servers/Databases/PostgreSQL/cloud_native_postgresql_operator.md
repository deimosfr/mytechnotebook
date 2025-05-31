---
title: "Cloud Native PostgreSQL: An Operator to run full featured PostgreSQL in Kubernetes"
slug: cloud-native-postgresql-operator/
description: "A guide to run a full featured PostgreSQL in Kubernetes using the PostgreSQL Operator."
categories: ["Database", "PostgreSQL", "Server", "Kubernetes"]
tags: ["PostgreSQL", "Database", "SQL", "Kubernetes", "Helm"]
---

## Introduction

[CloudNativePG](https://cloudnative-pg.io/) is the Kubernetes operator that covers the full lifecycle of a highly available PostgreSQL database cluster with a primary/standby architecture, using native streaming replication.

## Installation

To install the CloudNativePG operator, you can use the helm chart:

```bash
helm repo add cloudnative-pg https://cloudnative-pg.io/charts
helm repo update
helm install cloudnative-pg cloudnative-pg/cloudnative-pg
```

## Configuration

The CloudNativePG operator can be configured using the following parameters:

```yaml
replicaCount: 1 # Number of replicas to run
resources: {}
# If you're using Prometheus, you can enable the pod monitor
monitoring:
  podMonitorEnabled: false
```

## Deploy a cluster

First of all, set the login/password of the superuser. Let's define a secret for this:

=== "secret.yaml"

    ```yaml
    apiVersion: v1
    kind: Secret
    type: kubernetes.io/basic-auth
    metadata:
      name: pg-cluster-superuser
    data:
      password: cG9zdGdyZXM
      username: cG9zdGdyZXM
    ```

As you can see, the password and username are base64 encoded. You can use the following command to encode them:

```bash
echo -n "postgres" | base64
```

And then set the cluster configuration with [MetalLB service](../../Network/metallb_lb_k8s.md):

=== "cluster.yaml"

    ```yaml
    apiVersion: postgresql.cnpg.io/v1
    kind: Cluster
    metadata:
      name: pg-cluster
    spec:
      # Number of replicas to run
      instances: 2
      # Storage class to use
      storage:
        storageClass: openebs-lvm
        # Size of the storage
        size: 2Gi

      # Enable superuser access
      enableSuperuserAccess: true
      # Secret containing the superuser credentials we previously created
      superuserSecret:
        name: cluster-example-superuser

      # Resources for those PostgreSQL pods
      resources:
        requests:
          memory: "512Mi"
          cpu: "256m"
        limits:
          memory: "512Mi"
          cpu: "1"

      # Managed services
      managed:
        services:
          additional:
            # Point to the read/write load balancer service
            - selectorType: rw
              serviceTemplate:
                metadata:
                  name: pg-cluster-rw-lb
                  annotations:
                    # Set the IP address of the load balancer if you're using MetalLB.
                    metallb.universe.tf/loadBalancerIPs: 192.168.0.206
                spec:
                  type: LoadBalancer

      # Optional: Set the affinity for the pods
      affinity:
        enablePodAntiAffinity: true
        topologyKey: kubernetes.io/hostname
        podAntiAffinityType: required
        tolerations:
          - key: xxx
            operator: "Equal"
            value: "yyy"
            effect: NoSchedule
    ```

Then, deploy the cluster and the secret:

```bash
kubectl apply -f secret.yaml -f cluster.yaml
```

You can validate the deployment pods and services this way:

=== "services"

    ```bash
    $ kubectl get services
    NAME               TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
    pg-cluster-r       ClusterIP      10.43.71.9       <none>        5432/TCP       2d18h
    pg-cluster-ro      ClusterIP      10.43.186.251    <none>        5432/TCP       2d18h
    pg-cluster-rw      ClusterIP      10.43.65.15      <none>        5432/TCP       2d18h
    pg-cluster-rw-lb   LoadBalancer   10.43.169.20     192.168.0.206 5432:32088/TCP 2d17h
    ```

=== "pods"

    ```bash
    $ kubectl get pods
    NAME            READY   STATUS    RESTARTS   AGE
    pg-cluster-1    1/1     Running   0          2d17h
    pg-cluster-2    1/1     Running   0          2d17h
    ```

As you can see, multiple services:

- `pg-cluster-r`: Points to any PostgreSQL instance in the cluster (read).
- `pg-cluster-ro` Points to the replicas, where available (read-only).
- `pg-cluster-rw` Points to the primary instance of the cluster (read/write).
- `pg-cluster-rw-lb` is the load balancer service we've set earlier
