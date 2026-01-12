---
title: "Immich: The Self-Hosted Google Photos Alternative"
slug: immich-google-photos-alternative/
description: "Discover Immich, the high-performance self-hosted photo and video solution. Learn how to deploy it on Kubernetes with a dedicated Vector Database and VPA optimization."
categories: ["Linux", "Multimedia", "Kubernetes", "PostgreSQL"]
tags: ["Immich", "Photos", "Backup", "Self-Hosted", "VectorDB", "Postgres", "Kubernetes"]
---

![Immich Logo](../../static/images/immich_logo.svg){ width=400 }

Start taking control of your memories. [Immich](https://immich.app/) is a high-performance, self-hosted backup solution for photos and videos. It is designed to be a drop-in replacement for Google Photos, offering features like automatic backup, facial recognition, and a stunning mobile app, all without the privacy concerns of big tech.

We love Immich because it's fast, actively developed, and gives us complete ownership of our data. Let's walk through a robust Kubernetes deployment.

## Architecture & Requirements

Immich is composed of several microservices (Server, Microservices, Machine Learning, Postgres, Redis). Because of its smart search and facial recognition capabilities, it requires a **Vector Database**.

We will deploy:

*  **Immich Server & Microservices**: The core application.
*  **Vector Database**: A CloudNativePG cluster with `pgvector` extension for machine learning data.
*  **Gateway API**: To expose the service securely.

## Vector Database Configuration

Immich relies heavily on vector embeddings for its semantic search and facial recognition. Instead of a generic database, we will provision a dedicated PostgreSQL cluster utilizing the [Cloudnative Postgres Operator](../../Servers/Databases/PostgreSQL/cloud_native_postgresql_operator.md) and the [vectorchord](https://github.com/tensorchord/cloudnative-vectorchord) image which comes with `pgvector` pre-installed.

This configuration creates a specialized "Vector DB" cluster:

=== "cluster-vectordb.yaml"

    ```yaml
    apiVersion: postgresql.cnpg.io/v1
    kind: Cluster
    metadata:
      # We name it specifically for vector workloads
      name: pg-cluster-vectordb
    spec:
      instances: 3
      # Specialized image with pgvector support
      imageName: ghcr.io/tensorchord/cloudnative-vectorchord:17.5-0.4.3

      resources:
        requests:
          memory: "768Mi"
          cpu: "256m"
        limits:
          memory: "768Mi"
          cpu: "1"

      # Storage configuration
      storage:
        storageClass: openebs-lvm
        size: 5Gi

      postgresql:
        shared_preload_libraries:
          - "vchord.so"

      bootstrap:
        initdb:
          postInitSQL:
            - CREATE EXTENSION IF NOT EXISTS vchord CASCADE;

      # Optional: enable WAL archiving
      plugins:
        - name: barman-cloud.cloudnative-pg.io
          isWALArchiver: true
          parameters:
            barmanObjectName: postgres-backup-s3
    ```

!!! danger

    Do not try to update image version of Postgres Vector DB to the latest version. Rollback is impossible and Immich may not support it.
    Instead, check the supported versions on the [Immich Documentation](https://docs.immich.app/install/upgrading#do-i-still-need-pgvectors-installed-after-migrating-to-vectorchord).

After a few minutes, the cluster should be ready. You can check this way:

```bash
$ kubectl cnpg -n databases status pg-cluster-vectordb
Cluster Summary
Name                     databases/pg-cluster-vectordb
System ID:               7593739157184380948
PostgreSQL Image:        ghcr.io/tensorchord/cloudnative-vectorchord:17.5-0.4.3
Primary instance:        pg-cluster-vectordb-1
Primary promotion time:  2026-01-11 14:26:31 +0000 UTC (33h18m49s)
Status:                  Cluster in healthy state 
Instances:               3
Ready instances:         3
Size:                    747M
Current Write LSN:       0/6F000000 (Timeline: 1 - WAL File: 00000001000000000000006F)

Continuous Backup status (Barman Cloud Plugin)
ObjectStore / Server name:      postgres-backup-s3/pg-cluster-vectordb
First Point of Recoverability:  -
Last Successful Backup:         -
Last Failed Backup:             -
Working WAL archiving:          OK
WALs waiting to be archived:    0
Last Archived WAL:              00000001000000000000006E   @   2026-01-11T23:43:12.040003Z
Last Failed WAL:                -

Streaming Replication status
Replication Slots Enabled
Name                   Sent LSN    Write LSN   Flush LSN   Replay LSN  Write Lag  Flush Lag  Replay Lag  State      Sync State  Sync Priority  Replication Slot
----                   --------    ---------   ---------   ----------  ---------  ---------  ----------  -----      ----------  -------------  ----------------
pg-cluster-vectordb-2  0/6F000000  0/6F000000  0/6F000000  0/6F000000  00:00:00   00:00:00   00:00:00    streaming  async       0              active
pg-cluster-vectordb-3  0/6F000000  0/6F000000  0/6F000000  0/6F000000  00:00:00   00:00:00   00:00:00    streaming  async       0              active

Instances status
Name                   Current LSN  Replication role  Status  QoS        Manager Version  Node
----                   -----------  ----------------  ------  ---        ---------------  ----
pg-cluster-vectordb-1  0/6F000000   Primary           OK      Burstable  1.28.0           k8s-node-1.deimos.local
pg-cluster-vectordb-2  0/6F000000   Standby (async)   OK      Burstable  1.28.0           k8s-node-2.deimos.local
pg-cluster-vectordb-3  0/6F000000   Standby (async)   OK      Burstable  1.28.0           k8s-node-3.deimos.local

Plugins status
Name                            Version  Status  Reported Operator Capabilities
----                            -------  ------  ------------------------------
barman-cloud.cloudnative-pg.io  0.10.0   N/A     Reconciler Hooks, Lifecycle Service
```

## Immich Configuration

Now we configure Immich itself.

### Chart configuration

We create a `values-overrides.yaml` file to configure the chart:

=== "values-overrides.yaml"

    ```yaml
    controllers:
      main:
        containers:
          main:
            # set immich version (should be compatible with Postgres Vector DB image)
            image:
              tag: v2.4.1
            env:
              # use your valkey service name or leave this to use embedded valkey
              {% raw -%}
              REDIS_HOSTNAME: "{{ printf "%s-valkey" .Release.Name }}"
              {% endraw %}
              # use information from a secret
              DB_HOSTNAME:
                valueFrom:
                  secretKeyRef:
                    name: immich-database-app
                    key: host
              DB_USERNAME:
                valueFrom:
                  secretKeyRef:
                    name: immich-database-app
                    key: user
              DB_PASSWORD:
                valueFrom:
                  secretKeyRef:
                    name: immich-database-app
                    key: password
              DB_DATABASE_NAME:
                valueFrom:
                  secretKeyRef:
                    name: immich-database-app
                    key: db_name
    
    immich:
      # enable Prometheus metrics
      metrics:
        enabled: true
      # point to an existing pvc
      persistence:
        library:
          existingClaim: immich-library
      # configuration is immich-config.json converted to yaml
      # ref: https://immich.app/docs/install/config-file/
      configuration: {}
        # trash:
        #   enabled: false
        #   days: 30
        # storageTemplate:
        #   enabled: true
    
    # use the embedded valkey
    valkey:
      enabled: true
    
    # Immich components
    server:
      enabled: true
      ingress:
        main:
          enabled: false
   
    # disable machine learning if you don't need it
    machine-learning:
      enabled: false
    ```

### Connecting the Database

We create a Kubernetes Secret to manage the connection details to our Vector DB:

=== "secret.yaml"

    ```yaml
    apiVersion: v1
    kind: Secret
    metadata:
      name: immich-database-app
    type: Opaque
    stringData:
      # update all those fields to match your setup
      host: pg-cluster-vectordb
      user: immich
      password: immich
      db_name: immich
    ```

### Storage for Photos

Photos take up space! We need a specific Persistent Volume Claim (PVC) for the library. In my case, I'm using an SMB share here (likely from a NAS) to provide enough storage:

=== "pvc.yaml"

    ```yaml
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: immich-library
    spec:
      accessModes:
        - ReadWriteMany
      storageClassName: smb-hdd
      resources:
        requests:
          storage: 100Gi
    ```

### Exposing with Gateway API

To access Immich from the web (and for the mobile app to connect), we use the Kubernetes Gateway API:

=== "httproute.yaml"

    ```yaml
    apiVersion: gateway.networking.k8s.io/v1
    kind: HTTPRoute
    metadata:
      name: immich
    spec:
      hostnames:
      - immich.mydomain.com
      parentRefs:
      # my local gateway
      - name: local
        namespace: kube-system
        sectionName: https
      rules:
      - backendRefs:
        # Target the Immich service name
        - kind: Service
          name: immich-server
          port: 2283
        matches:
        - path:
            type: PathPrefix
            value: /
    ```

# Immich Installation

To install Immich, we'll first deploy configuration we setup earlier:

```
kubectl apply -f secret.yaml -f pvc.yaml -f httproute.yaml
```

Then we use the [Helm chart](https://github.com/immich-app/immich-helm-chart) provided by the Immich team:

```
helm repo add immich https://immich-app.github.io/immich-charts
helm install immich immich/immich --version 0.10.3 -f values-overrides.yaml
```

You immich instance should now be accessible at `immich.mydomain.com`.

## Best Practices

1.  **Dedicated Database**: Separating your Vector DB allows you to tune it specifically for vector search workloads without affecting other postgres instances.
2.  **Backups**: Immich data is precious. Ensure your `pvc` and your Postgres clusters have automated backup policies.
