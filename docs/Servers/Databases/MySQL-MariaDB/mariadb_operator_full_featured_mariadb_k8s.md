---
title: "MariaDB Operator: Run full featured MariaDB in Kubernetes"
slug: mariadb-operator-full-featured-mariadb-k8s/
description: "A guide to run a full featured MariaDB in Kubernetes using the MariaDB Operator."
categories: ["Database", "MariaDB", "Server", "Kubernetes"]
tags: ["MariaDB", "Database", "SQL", "Kubernetes", "Helm"]
---

## Introduction

[MariaDB](./mariadb_migration_from_mysql.md) is a community-developed, commercially supported fork of the MySQL relational database management system. MariaDB is a drop-in replacement for MySQL and is designed to be a more robust, scalable, and secure database solution.

The [MariaDB Operator](https://github.com/mariadb-operator/mariadb-operator) is a tool that allows you to run a full featured MariaDB in Kubernetes. It is a great way to run a MariaDB server in Kubernetes and it is a great way to manage your MariaDB server in Kubernetes.

The benefits are:

- It is a full featured MariaDB server solution
- It is a great way to run a MariaDB server in Kubernetes
- You don't have to manage the MariaDB server manually for backups, upgrades, etc.
- There are build-in ways to scale the MariaDB server
- Easily extensible with MaxScale, [Replications with auto-failover](./replication_master_to_master.md), [GaleraDB](./mariadb_galera_cluster_multimaster_replication.md)...and easily customizable

Before starting, you need to have a Kubernetes cluster.

## Installation

First, you need to install the MariaDB Operator with Helm:

```bash
helm repo add mariadb-operator https://helm.mariadb.com/mariadb-operator
helm install mariadb-operator-crds mariadb-operator/mariadb-operator-crds
helm install mariadb-operator mariadb-operator/mariadb-operator
```

For best practices and later usage, I advise to set the metrics and resources like:

=== "mariadb-operator-override-values.yaml"

    ```yaml
    metrics:
      enabled: true
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1
          memory: 256Mi
    ```

You can update the `mariadb-operator` chart with the following values file (simply add `-f mariadb-operator-override-values.yaml` to the helm install command).

If you check your pods, you should see something like this:

```bash
$ kubectl get pods -n mariadb-operator
NAME                                                READY   STATUS    RESTARTS   AGE
mariadb-operator-855f9bbdfd-4ggmh                   1/1     Running   0          8h
mariadb-operator-cert-controller-67f78fc9f4-5ttd7   1/1     Running   0          8h
mariadb-operator-webhook-5d8c997f76-tg5rj           1/1     Running   0          8h
```

## Single instance configuration

You can configure a single instance of MariaDB with the following configuration:

=== "mariadb-instance.yaml"

    ```yaml
    apiVersion: k8s.mariadb.com/v1alpha1
    kind: MariaDB
    metadata:
      # The name of the MariaDB instance
      name: mariadb
    spec:
      timeZone: "UTC"
      # Optional: set the desired version of MariaDB
      # image:
      #   tag: "11.0.3"
      # Root password
      rootPasswordSecretKeyRef:
        name: mariadb
        key: root-password
      # MariaDB resources
      resources:
        requests:
          cpu: 100m
          memory: 512Mi
        limits:
          cpu: 1
          memory: 512Mi
      # Storage configuration
      storage:
        size: 2Gi
        # Prefer a local path storage class for better performance
        storageClassName: openebs-lvm
        resizeInUseVolumes: true
        waitForVolumeResize: true
      myCnf: |
        [mariadb]
        bind-address=*
        default_storage_engine=InnoDB
        binlog_format=row
        innodb_autoinc_lock_mode=2
        innodb_buffer_pool_size=256M
        max_allowed_packet=128M
      # Optional: set the service type to LoadBalancer if you want to expose the MariaDB instance.
      service:
        type: LoadBalancer
        metadata:
          annotations:
            # Here we use metallb.
            metallb.universe.tf/loadBalancerIPs: 192.168.0.1
        externalTrafficPolicy: Local
        sessionAffinity: None
      # Enable metrics if you have a prometheus operator installed
      metrics:
        enabled: false
    ```

Regarding the StorageClass, I'm using [OpenEBS LVM](../../Containers/Kubernetes/use_local_storage_with_openebs.md) for the best performance. It's not mandatory and you can use any other storage class. However, for maximum performances, you should use local storage.

=== "mariadb-secret.yaml"

    ```yaml
      apiVersion: v1
      kind: Secret
      metadata:
        name: mariadb
        labels:
          k8s.mariadb.com/watch: ""
      stringData:
        # Set the root password for the MariaDB instance
        root-password: MariaDB11!
    ```

Apply this configuration with kubectl:

```bash
kubectl apply -f mariadb-instance.yaml -f mariadb-secret.yaml
```

And you should see the MariaDB instance running:

```bash
$ kubectl get sts mariadb
NAME      READY   AGE
mariadb   1/1     18m
```

## Create a user

You can decide to create a dedicated secret for this new user or use the same secret as the MariaDB instance. Here we're going to update the existing secret:

=== "mariadb-secret.yaml"

    ```yaml
      apiVersion: v1
      kind: Secret
      metadata:
        name: mariadb
        labels:
          k8s.mariadb.com/watch: ""
      stringData:
        # Set the root password for the MariaDB instance
        root-password: MariaDB11!
        mariadb-user: passord
    ```

You can now create a user on the MariaDB instance with the following configuration:

=== "mariadb-user.yaml"

    ```yaml
    apiVersion: k8s.mariadb.com/v1alpha1
    kind: User
    metadata:
      # The name of the user in MariaDB
      name: mariadb-user
    spec:
      # Select the MariaDB instance to create the user on
      mariaDbRef:
        name: mariadb
      # The login and password for the user
      passwordSecretKeyRef:
        name: mariadb # name of the secret
        key: mariadb-user # name of the key in the secret
      # This field defaults to 10
      maxUserConnections: 20
      # The host to allow the user to connect from
      host: "%"
      # Delete the resource in the database whenever the CR gets deleted.
      # Alternatively, you can specify Skip in order to omit deletion.
      cleanupPolicy: Delete
      requeueInterval: 30s
      retryInterval: 5s
    ```

You can now apply this configuration with kubectl:

```bash
kubectl apply -f mariadb-user.yaml -f mariadb-secret.yaml
```

## Grant privileges

You can grant privileges to a user with the following configuration:

=== "mariadb-grant.yaml"

    ```yaml
    apiVersion: k8s.mariadb.com/v1alpha1
    kind: Grant
    metadata:
      name: grant
    spec:
      mariaDbRef:
        name: mariadb
      privileges:
        - "SELECT"
        - "INSERT"
        - "UPDATE"
        - "ALL PRIVILEGES"
      database: "*"
      table: "*"
      username: mariadb-user # name of the user in MariaDB
      grantOption: true
      host: "%"
      # Delete the resource in the database whenever the CR gets deleted.
      # Alternatively, you can specify Skip in order to omit deletion.
      cleanupPolicy: Delete
      requeueInterval: 30s
      retryInterval: 5s
    ```

## High availability

The MariaDB Operator also supports multiple High Availability solutions like:

- [Replications with auto-failover](./replication_master_to_master.md): a master-master replication solution with auto-failover capabilities, also called SemiSync replication.
- [Galera](./mariadb_galera_cluster_multimaster_replication.md): a multi-master replication solution with auto-failover capabilities.
- MaxScale: a proxy solution with failover and load balancing capabilities

We'll take a look here at the Replications with auto-failover solution. I like this one because it's built-in and it's very easy to configure for a home lab. Prefer the Galera solution if you need a multi-master solution. And use MaxScale if you need a more complex solution with failover and load balancing capabilities.

### Replications with auto-failover

Here is a configuration with only 2 nodes (not perfect because of the missing Qorum) but it's a good start for a home lab.

=== "mariadb-replication.yaml"

    ```yaml
    apiVersion: k8s.mariadb.com/v1alpha1
    kind: MariaDB
    metadata:
      name: mariadb
    spec:
      timeZone: "UTC"
      # Root password
      rootPasswordSecretKeyRef:
        name: mariadb
        key: root-password
      # MariaDB resources
      resources:
        requests:
          cpu: 100m
          memory: 768Mi
        limits:
          cpu: 1
          memory: 768Mi
      # Storage configuration
      storage:
        size: 2Gi
        storageClassName: openebs-lvm
        resizeInUseVolumes: true
        waitForVolumeResize: true
      # MariaDB configuration
      myCnf: |
        [mariadb]
        bind-address=*
        default_storage_engine=InnoDB
        binlog_format=row
        innodb_autoinc_lock_mode=2
        innodb_buffer_pool_size=256M
        max_allowed_packet=128M
      # Enable replication with auto-failover
      replication:
        enabled: true
        probesEnabled: true
        primary:
          automaticFailover: true
      # Number of replicas with even number of replicas (here 2)
      replicas: 2
      replicasAllowEvenNumber: true
      # On update,the replicas will be updated first and then the primary
      updateStrategy:
        type: ReplicasFirstPrimaryLast
      # To be used for read requests. It will point to all nodes
      service:
        type: LoadBalancer
        metadata:
          annotations:
            metallb.universe.tf/loadBalancerIPs: 192.168.0.1
        externalTrafficPolicy: Local
        sessionAffinity: None
      # To be used for write requests. It will point to a single node, the primary.
      primaryService:
        type: LoadBalancer
        metadata:
          annotations:
            metallb.universe.tf/loadBalancerIPs: 192.168.0.2
      # To be used for read requests. It will point to all nodes, except the primary.
      secondaryService:
        type: LoadBalancer
        metadata:
          annotations:
            metallb.universe.tf/loadBalancerIPs: 192.168.0.3
      # Ensure no instance is scheduled on the same node
      affinity:
        antiAffinityEnabled: true
      # Ensure no more than 1 instance is unavailable
      podDisruptionBudget:
        maxUnavailable: 1
      # Enable metrics if you have a prometheus operator installed
      metrics:
        enabled: false
    ```

Once applied, you can check replication status with:

```bash hl_lines="25 27 28 32 34 39 40"
$ kubectl describe MariaDB mariadb
Name:         mariadb
...
Spec:
  Affinity:
  Replicas:                                   2
  Replicas Allow Even Number:                 true
  Replication:
    Enabled:  true
    Primary:
      Automatic Failover:  true
      Pod Index:           0
    Probes Enabled:        true
    Replica:
      Connection Retries:  10
      Connection Timeout:  10s
      Gtid:                CurrentPos
      Sync Timeout:        10s
      Wait Point:          AfterSync
    Sync Binlog:           true
...
Status:
  Conditions:
    Last Transition Time:     2025-05-20T22:16:34Z
    Message:                  Running
    Reason:                   StatefulSetReady
    Status:                   True
    Type:                     Ready
    Last Transition Time:     2025-05-20T22:15:33Z
    Message:                  Updated
    Reason:                   Updated
    Status:                   True
    Type:                     Updated
  Current Primary:            mariadb-0
  Current Primary Pod Index:  0
  Default Version:            11.4
  Replicas:                   2
  Replication Status:
    mariadb-0:  Master
    mariadb-1:  Slave
...
```

## Troubleshooting

If you encounter issues, you can check the logs of the MariaDB instance with the following command:

```bash
kubectl logs -n mariadb-operator -l app.kubernetes.io/name=mariadb-operator
```

## Resources

- [MariaDB Operator: all available configuration options](https://github.com/mariadb-operator/mariadb-operator/blob/main/docs/API_REFERENCE.md)
- [Configurations examples](https://github.com/mariadb-operator/mariadb-operator/tree/main/examples)
