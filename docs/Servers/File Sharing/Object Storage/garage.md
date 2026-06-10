---
title: "Garage: a lightweight S3 compatible object storage"
slug: garage/
description: "Garage is a lightweight S3 compatible object storage"
categories: ["Servers", "Object Storage"]
tags: ["Servers", "Object Storage", "Garage", "S3"]
---

![Garage](../../../static/images/garage-logo.svg)

[Garage](https://garagehq.deuxfleurs.fr/) is a lightweight S3 compatible object storage. It is a self-hosted, open-source, and easy-to-deploy solution for storing and managing large amounts of data. It is written in Rust and is available as a Docker image.

I used to look into an S3 storage solution for my home lab. Something really simple, I don't need replication, high availability, or any of the other features that come with it. I just needed a simple, easy-to-use, and easy-to-manage solution. Garage was a good candidate for me.

For information, [Minio](https://github.com/minio/minio) has been declared as maintenance mode and will not receive any new features or updates.
The other interesting solution is [RustFS](https://github.com/rustfs/rustfs), but minimum required resources looks to be higher than Garage.

## Garage Installation

The installation I chose to use is [Kubernetes](../../Containers/Kubernetes/index.md). Unfortunatelly, the Helm Chart maturity is not perfect yet, mostly if you want to use it with the webui. So I had to manually bypass some configuration to get it to work.

### Prepare configuration

First of all, we're not going to use the configuration inside the chart, we're going to use a custom ConfigMap instead:

=== "configmap.yaml"

    ```yaml
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: garage-config
    data:
      garage.toml: |-
        metadata_dir = "/mnt/meta"
        data_dir = "/mnt/data"

        db_engine = "lmdb"

        block_size = "10M"

        replication_factor = 1
        consistency_mode = "consistent"

        compression_level = 5

        rpc_bind_addr = "[::]:3901"
        # Generate one a token with: openssl rand -hex 32
        rpc_secret = "xxx"

        bootstrap_peers = []

        [kubernetes_discovery]
        namespace = "storage"
        service_name = "garage"
        skip_crd = false

        [s3_api]
        s3_region = "us-east-1"
        api_bind_addr = "[::]:3900"
        root_domain = ".s3.mydomain.com"

        [s3_web]
        bind_addr = "[::]:3902"
        root_domain = ".web.mydomain.com"
        index = "index.html"

        [admin]
        api_bind_addr = "[::]:3903"
        # Generate one a token with: openssl rand -hex 32
        admin_token = "xxx"
    ```

The important things here are:

- `replication_factor`: set to 1 if you don't want replication
- `kubernetes_discovery`: set to true to use kubernetes discovery
- `rpc_secret` and `admin_token`: generate a random token with `openssl rand -hex 32` for each of them

You can find the [full configuration here](https://garagehq.deuxfleurs.fr/documentation/reference-manual/configuration/).

The other thing to deploy is a service for the admin interface because the chart doesn't provide it for some reasons:

!!! quote

    NOTE: The admin API is excluded for now as it is not consistent across nodes

In our case, as there is a single instance, there won't be any issues, so we'll add this service:

=== "service.yaml"

    ```yaml
    apiVersion: v1
    kind: Service
    metadata:
      name: garage-admin
    spec:
      type: ClusterIP
      ports:
      - name: s3-admin
        port: 3903
        protocol: TCP
        targetPort: 3903
      selector:
        app.kubernetes.io/instance: garage
        app.kubernetes.io/name: garage
    ```

Now we can deploy those resources (in the same namespace where the garage chart will be deployed):

```bash
kubectl apply -f configmap.yaml -f service.yaml
```

### Deploy Garage chart

Unfortunately the chart is only available on the official Garage repository. So we'll need to clone the repo and then being able to use it:

```bash
git clone https://git.deuxfleurs.fr/Deuxfleurs/garage
cd garage/scripts/helm
```

You can create an override file to override the default valuees. Here is the one I used:

=== "values-override.yaml"

    ```yaml
    garage:
      rpcSecret: "<Use the RPC secret from the configmap>"
      existingConfigMap: "garage-config"

    persistence:
      enabled: true
      meta:
        storageClass: "<Use the storageclass from the configmap>"
        size: 2Gi
      data:
        storageClass: "<Use the storageclass from the configmap>"
        size: 50Gi

    deployment:
      replicaCount: 1

    service:
      s3:
        api:
          port: 3900
        web:
          port: 3902
    ```

- `persistence.meta.size`: the size of the metadata volume. Start with small size and grow the `pvc` size when needed.
- `persistence.data.size`: the size of the data volume
- `garage.rpcSecret`: the RPC secret from the configmap above
- `garage.existingConfigMap`: the name of the configmap we deployed above

Now we can deploy:

```bash
helm install garage ./garage -f values-override.yaml
```

!!! note

    CRD will be deployed by the chart, but you can use `garage.kubernetesSkipCrd: true` to skip it. But as it's more convenient, we'll keep them like this.

## Garage Web UI

We're almost ready to use Garage but we will finish with the [Garage Web UI](https://github.com/khairul169/garage-webui) installation. Helm chart is not ready yet but a [PR is in progress](https://github.com/khairul169/garage-webui/pull/46).

Let's clone the PR:

```bash
git clone https://github.com/rojinebrahimi/garage-webui.git
cd deploy/helm
```

And prepare an override file:

=== "values-override.yaml"

    ```yaml
    httproute:
      enabled: true
      parentRefs:
        - name: gateway
          namespace: kube-system
      hostnames:
        - "garage.mydomain.com"

    config:
      enabled: true
      name: "garage-config"

    garageConfig:
      s3Endpoint: "http://garage:3900"
      region: "us-east-1"
      adminApiUrl: "http://garage-admin:3903"
      adminApiKey: "<Use the admin token from the Garage admin config>"

    auth:
      # Generate this using: htpasswd -nbBC 10 "admin" "yourpassword"
      userPassHash: "user:password"
    ```

Here is what you have to update:

- `garageConfig.adminApiKey`: the admin token from the Garage admin config (in the `[admin]` section)
- `auth.userPassHash`: the user password hash (using `htpasswd -nbBC 10 "user" "password"`)
- `config.name`: the name of the configmap we deployed for garage

I personnaly use Gateway API to expose the web ui, but you can use Ingress or NodePort if you prefer. Now you can look at the web ui at [http://garage.mydomain.com](http://garage.mydomain.com):

![Garage Web UI dashboard](../../../static/images/garage_dashboard_unavailable.avif)

You should see "Unavailable" status, and one connected node. This is normal, as we only have one node.

### Create a partition

Now go into the "Cluster" tab, you should see the node ID in "Active" status. Now we can assign a new partition:

![Garage Web UI assign node capacity](../../../static/images/garage_assign_node.avif)

1. Click on "Assign"
2. Set the Zone
3. Set the partition capacity, **do not got beyong the data disk capacity you've set**
4. Click on "Save"
5. Click on "Apply"

In the end, you should see a summary like:

```
==== COMPUTATION OF A NEW PARTITION ASSIGNATION ====

Partitions are replicated 1 times on at least 1 distinct zones.

Optimal partition size:                     175.8 MB
Usable capacity / total cluster capacity:   45.0 GB / 45.0 GB (100.0 %)
Effective capacity (replication factor 1):  45.0 GB

us-east-1           Tags  Partitions        Capacity  Usable capacity
  be004718df343c99  []    256 (256 new)     45.0 GB   45.0 GB (100.0%)
  TOTAL                   256 (256 unique)  45.0 GB   45.0 GB (100.0%)
```

If you look into the dashboard, you should see "Healthy" status:

![Garage Web UI dashboard](../../../static/images/garage_dashboard_ok.avif)

### Create a bucket

Now go into the "Buckets" tab and click on "Create bucket", then name it:

![Garage Web UI create bucket](../../../static/images/garage_create_bucket.avif)

### Create a key

Go into "Keys" tab and click on "Create key", then name it:

![Garage Web UI create key](../../../static/images/garage_create_key.avif)

You can then see the Key ID and Secret key.

### Manage bucket permissions

If you click on "Manage" on a bucket, you'll be able to set quotas, website access and permissions. Click on the "Permissions" tab and "Allow key". Then add the permissions you wish and click on "Submit".

![Garage Web UI manage bucket permissions](../../../static/images/garage_bucket_permissions.avif)

## API exposure

To be able to use the API inside the cluster, you first need to decide how you want to expose it:

=== "In cluster"

     To access the API from inside the cluster, you can use the service name `garage` and the port `3900`. In case you need to access from another namespace, the internal FQDN is:

     * `garage.<namespace>.svc` or
     * `garage.<namespace>.svc.cluster.local`

=== "LoadBalancer"

    A simple solution is to update the `service` part in the `values.yaml` file of the garage chart and change the service type:

    ```yaml
    service:
      type: LoadBalancer
      s3:
        api:
          port: 3900
        web:
          port: 3902
    ```

    Then you'll have a LoadBalancer IP assigned to the service. You can use it to access the API from outside the cluster.

=== "Gateway API"

    You can use Gateway API to expose the API. Useful if you want to expose it publicly:

    ```yaml
    apiVersion: gateway.networking.k8s.io/v1alpha2
    kind: TCPRoute
    metadata:
      name: garage-storage
    spec:
      parentRefs:
      - name: external-gateway
        sectionName: s3-api
      rules:
      - backendRefs:
        - name: garage
          port: 3900
    ```

    And finally ensure you have the correct listener properly set on your Gateway:

    ```yaml
    apiVersion: gateway.networking.k8s.io/v1alpha2
    kind: Gateway
    metadata:
      name: external-gateway
    spec:
    listeners:
      - name: s3-api
        protocol: TCP
        port: 3900
        allowedRoutes:
          namespaces:
            from: Same
    ```

## CLI usage

You can easily use AWS CLI to interact with Garage. First, you need to set the environment variables:

```bash
export AWS_ACCESS_KEY_ID="<Your Key ID>"
export AWS_SECRET_ACCESS_KEY="<Your Secret Key>"
export AWS_DEFAULT_REGION="us-east-1"
export AWS_ENDPOINT_URL="http://garage:3900"
```

- `AWS_DEFAULT_REGION`: set to the region set in the [garage config](#prepare-configuration).
- `AWS_ENDPOINT_URL`: use the one from the [API exposure](#api-exposure) method you've chosen.

Then you can use the AWS CLI to interact with Garage. For example, to list all buckets (with recent AWS CLI, `--endpoint-url` is required):

```bash
aws s3 ls --endpoint-url http://garage:3900
```

If you want to have more examples, you can look at the [official documentation](https://garagehq.deuxfleurs.fr/documentation/connect/cli/).

# Troubleshooting

## Garage S3 Metadata Recovery (LMDB Corruption)

Here is how to recover a Garage node from metadata database corruption (e.g., `LMDB: MDB_CORRUPTED`) **while preserving the node's unique ID/keys** in a replicated Garage cluster. I had to do it for my Docker registry pointing to Garage.

---

### Symptoms

You may see errors in your client applications (like a Docker registry or log collectors) such as:

- `503 Service Unavailable`
- `unknown: blob upload unknown to registry` (due to interrupted multi-part uploads)

When checking the Docker Registry or Garage logs, you will see the following error indicating a database page corruption on one of the nodes:

```text
Could not reach quorum of 1 (sets=Some(1)). 0 of 1 request succeeded, others returned errors:
["ef427134c5a3bc85: Remote error: DB error: LMDB: MDB_CORRUPTED: Located page was wrong type"]
```

_Note: In the example above, `ef427134c5a3bc85` corresponds to node `garage-1`._

---

### Recovery Procedure

#### Step 1: Scale Down the StatefulSet

To release the Persistent Volume Claim (PVC) of the corrupted node so we can mount it to a recovery pod, scale down the Garage StatefulSet to `1` replica:

```bash
kubectl scale statefulset garage --replicas=1
```

Verify that the corrupted pod (e.g., `garage-1`) has terminated:

```bash
kubectl get pods
```

---

#### Step 2: Spin Up a Recovery Helper Pod

Create a helper pod manifest (`garage-recovery.yaml`) that mounts the metadata PVC of the corrupted node (`meta-garage-1`):

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: garage-recovery
spec:
  containers:
    - name: helper
      image: busybox
      command: ["sh", "-c", "sleep 3600"]
      volumeMounts:
        - mountPath: /volume
          name: meta-volume
  volumes:
    - name: meta-volume
      persistentVolumeClaim:
        claimName: meta-garage-1
```

Apply the manifest:

```bash
kubectl apply -f garage-recovery.yaml
```

Wait until the helper pod is running:

```bash
kubectl get pod garage-recovery
```

---

#### Step 3: Remove the Corrupted Database

Exec into the helper pod to view the metadata folder.

!!! Warning

    Do **NOT** wipe the entire folder. You must preserve the `node_key` and `node_key.pub` files to prevent the node from generating a new cryptographic identity (which would break cluster layouts). Only delete the database folder.

Delete only the `db.lmdb` directory:

```bash
kubectl exec garage-recovery -- rm -rf /volume/db.lmdb
```

---

#### Step 4: Scale Up Garage

Delete the recovery helper pod:

```bash
kubectl delete -f garage-recovery.yaml
```

Scale the Garage StatefulSet back to its original replica count:

```bash
kubectl scale statefulset garage --replicas=2
```

---

#### Step 5: Trigger Metadata Resync

Since the database on the restored node is now empty, it will report `403 Forbidden` errors for S3 keys until it receives the database tables.

Trigger a repair run from the healthy node (`garage-0`) to push the synchronized tables and blocks to the restored node:

```bash
# Sync metadata tables
kubectl exec garage-0 -- /garage repair -a --yes tables

# Sync storage blocks
kubectl exec garage-0 -- /garage repair -a --yes blocks
```

---

#### Step 6: Verify Cluster Health

Verify that the restored node has successfully joined as a healthy member:

```bash
kubectl exec garage-0 -- /garage status
```

You should see both nodes listed under `==== HEALTHY NODES ====`:

```text
==== HEALTHY NODES ====
ID                Hostname  Address          Tags  Zone       Capacity  DataAvail         Version
a521e2d1bff9b439  garage-0  10.0.0.222:3901  []    us-east-1  46.6 GiB  75.5 GiB (96.3%)  v2.3.0
ef427134c5a3bc85  garage-1  10.0.4.17:3901   []    us-east-1  45.5 GiB  69.4 GiB (88.4%)  v2.3.0
```
