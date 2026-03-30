---
title: "Cloud Native PostgreSQL: An Operator to run full featured PostgreSQL in Kubernetes"
slug: cloud-native-postgresql-operator/
description: "A guide to run a full featured PostgreSQL in Kubernetes using the PostgreSQL Operator."
categories: ["Database", "PostgreSQL", "Server", "Kubernetes"]
tags: ["PostgreSQL", "Database", "SQL", "Kubernetes", "Helm"]
---

![PostgreSQL](../../../static/images/postgresql-logo.avif)

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

And then set the cluster configuration with [MetalLB service](../../Containers/Kubernetes/metallb_lb_k8s.md):

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

## CLI

The CloudNativePG operator provides a CLI tool to manage your PostgreSQL clusters!!! This is a must have tool to manage your clusters. It allows you to manage its lifecycle, backup, restore, etc.

To install the CLI, you can use the following command:

```bash
kubectl krew install cnpg
```

Then, you can check your cluster status for example:

```bash
$ kubectl cnpg -n databases status pg-cluster

Cluster Summary
Name                     databases/pg-cluster
System ID:               7584493819456610323
PostgreSQL Image:        ghcr.io/cloudnative-pg/postgresql:18.1-system-trixie
Primary instance:        pg-cluster-1
Primary promotion time:  2025-12-16 16:29:45 +0000 UTC (126h11m10s)
Status:                  Cluster in healthy state 
Instances:               2
Ready instances:         2
Size:                    160M
Current Write LSN:       0/9000000 (Timeline: 1 - WAL File: 000000010000000000000009)

Continuous Backup status (Barman Cloud Plugin)
ObjectStore / Server name:      postgres-backup-s3/pg-cluster
First Point of Recoverability:  -
Last Successful Backup:         -
Last Failed Backup:             -
Working WAL archiving:          OK
WALs waiting to be archived:    0
Last Archived WAL:              000000010000000000000008   @   2025-12-21T21:05:07.539965Z
Last Failed WAL:                000000010000000000000006   @   2025-12-21T21:04:59.898354Z

Streaming Replication status
Replication Slots Enabled
Name          Sent LSN   Write LSN  Flush LSN  Replay LSN  Write Lag  Flush Lag  Replay Lag  State      Sync State  Sync Priority  Replication Slot
----          --------   ---------  ---------  ----------  ---------  ---------  ----------  -----      ----------  -------------  ----------------
pg-cluster-2  0/9000000  0/9000000  0/9000000  0/9000000   00:00:00   00:00:00   00:00:00    streaming  async       0              active

Instances status
Name          Current LSN  Replication role  Status  QoS        Manager Version  Node
----          -----------  ----------------  ------  ---        ---------------  ----
pg-cluster-1  0/9000000    Primary           OK      Burstable  1.28.0           node1.mycompany.com
pg-cluster-2  0/9000000    Standby (async)   OK      Burstable  1.28.0           node2.mycompany.com

Plugins status
Name                            Version  Status  Reported Operator Capabilities
----                            -------  ------  ------------------------------
barman-cloud.cloudnative-pg.io  0.9.0    N/A     Reconciler Hooks, Lifecycle Service
```

## Backups

With CloudNativePG, backups are managed by a third party tool called [Barman Cloud Plugin](https://cloudnative-pg.io/plugin-barman-cloud/docs/intro/). We'll see here how to perform backups with PITR (Point In Time Recovery) on an [Object Storage](../../File%20Sharing/Object%20Storage/index.md) and we'll take [Garage](../../File%20Sharing/Object%20Storage/garage.md) as an example.

Installation is simple:

```bash
helm repo add cloudnative-pg https://cloudnative-pg.io/charts/
helm install plugin-barman-cloud cloudnative-pg/plugin-barman-cloud
```

Then we'll start to configure it. So first, we'll add a configuration explaining how backups will be made with S3:

=== "postgres-backup-s3.yaml"

    ```yaml
    apiVersion: barmancloud.cnpg.io/v1
    kind: ObjectStore
    metadata:
      name: postgres-backup-s3
    spec:
      # the retention policy
      retentionPolicy: "30d"
      configuration:
        # the folder inside the bucket where backups will be stored
        destinationPath: s3://backups-postgres
        # the URL endpoint
        endpointURL: http://garage.<namespace>.svc:3900
        s3Credentials:
          # the secret containing the credentials
          accessKeyId:
            name: postgres-backup-s3-secrets
            key: ACCESS_KEY_ID
          secretAccessKey:
            name: postgres-backup-s3-secrets
            key: ACCESS_SECRET_KEY
        wal:
          compression: gzip
    ```

Now the secret containing the credentials:

=== "postgres-backup-s3-secrets.yaml"

    ```yaml
    apiVersion: v1
    kind: Secret
    metadata:
      name: postgres-backup-s3-secrets
    stringData:
      # set the access key ID
      ACCESS_KEY_ID: <access-key-id>
      # set the access secret key
      ACCESS_SECRET_KEY: <access-secret-key>
    ```

Now the configuration is correct. We need to apply it to an existing cluster in the `spec` section of the cluster configuration for WAL archiving:

=== "cluster.yaml"

    ```yaml
    apiVersion: postgresql.cnpg.io/v1
    kind: Cluster
    metadata:
      name: pg-cluster
    spec:
      ...
      # other configurations...
      plugins:
      - name: barman-cloud.cloudnative-pg.io
        isWALArchiver: true
        parameters:
          barmanObjectName: postgres-backup-s3
    ```

And finally, setup the backup scheduling:

=== "backup-schedule.yaml"

    ```yaml
    apiVersion: postgresql.cnpg.io/v1
    kind: ScheduledBackup
    metadata:
      name: postgres-backup-s3
    spec:
      cluster:
        name: pg-cluster
      schedule: '0 23 * * *'
      backupOwnerReference: self
      method: plugin
      pluginConfiguration:
        name: barman-cloud.cloudnative-pg.io
    ```

Now apply the configuration:

```bash
kubectl apply -f cluster.yaml -f postgres-backup-s3.yaml -f postgres-backup-s3-secrets.yaml -f backup-schedule.yaml
```

You should see the cluster WAL coming into the S3 bucket.

## Restore

If all your nodes are in a bad shape and you need to restore the cluster. Here is the procedure:

1. We'll start by deleting the existing cluster (if we want to keep the same name):
    ```bash
    kubectl delete cluster pg-cluster
    ```

2. Then, let's update our cluster configuration by adding the recovery section:

    === "cluster.yaml (with PITR)"

        ```yaml
        # Backup recovery bootstrap
        bootstrap:
          recovery:
            source: origin
            recoveryTarget:
              targetTime: "2026-03-30 11:14:21.000000+00"
        # Define the source of your recovery config
        externalClusters:
          - name: origin
            plugin:
              name: barman-cloud.cloudnative-pg.io
              parameters:
                barmanObjectName: postgres-backup-s3
                # add a new name to distinguich the new cluster and avoid restoration issues
                serverName: pg-cluster-restored
        ```

    === "cluster.yaml (without PITR)"

        ```yaml
        # Backup recovery bootstrap
        bootstrap:
          recovery:
            source: origin
        # Define the source of your recovery config
        externalClusters:
          - name: origin
            plugin:
              name: barman-cloud.cloudnative-pg.io
              parameters:
                barmanObjectName: postgres-backup-s3
                serverName: pg-cluster
        plugins:
          - name: barman-cloud.cloudnative-pg.io
            isWALArchiver: true
            parameters:
              barmanObjectName: postgres-backup-s3
              # add a new name to distinguich the new cluster and avoid restoration issues
              serverName: pg-cluster-restored
        ```

3. Finally, we'll apply the configuration:
    ```bash
    kubectl apply -f cluster.yaml
    ```

!!! note
    
    Once the cluster is restored, you need to update the cluster configuration to **remove or comment** the `bootstrap` section and the `externalClusters` section.