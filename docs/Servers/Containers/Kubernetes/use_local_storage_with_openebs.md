---
title: "Use local storage with OpenEBS"
slug: use-local-storage-with-openebs/
description: "Use your local storage with LVM/ZFS with OpenEBS in Kubernetes."
categories: ["Storage", "OpenEBS", "Kubernetes", "Home Lab", "LVM", "ZFS"]
tags: ["Storage", "OpenEBS", "Kubernetes", "Home Lab", "LVM", "ZFS"]
---

## Introduction

[OpenEBS](https://openebs.io/) is a storage solution for Kubernetes. It is a lightweight, open-source project that provides a storage solution for Kubernetes.

OpenEBS turns any storage available to Kubernetes worker nodes into Local or Replicated Kubernetes Persistent Volumes. OpenEBS helps application and platform teams easily deploy Kubernetes stateful workloads that require fast and highly durable, reliable, and scalable Container Native Storage. OpenEBS is also a leading choice for NVMe-based storage deployments.

!!! quote

    **In the case of Local Volumes:**

    - OpenEBS can create persistent volumes or use sub-directories on Hostpaths or use locally attached storage or sparse files or over existing [LVM](../../../Linux/Filesystems%20And%20Storage/LVM/index.md) or [ZFS](../../../Solaris/Filesystems/zfs_the_filesystem_par_excellence.md) stack.
    - The local volumes are directly mounted into the Stateful Pod, without any added overhead from OpenEBS in the data path, decreasing latency.
    - OpenEBS provides additional tooling for local volumes for monitoring, backup/restore, disaster recovery, snapshots when backed by [LVM](../../../Linux/Filesystems%20And%20Storage/LVM/index.md) or [ZFS](../../../Solaris/Filesystems/zfs_the_filesystem_par_excellence.md) stack, capacity-based scheduling, and more.

    **In the case of Replicated Volumes:**

    - OpenEBS Replicated Storage creates an NVMe target accessible over TCP, for each persistent volume.
    - The Stateful Pod writes the data to the NVMe-TCP target that synchronously replicates the data to multiple nodes in the cluster. The OpenEBS engine itself is deployed as a pod and orchestrated by Kubernetes. When the node running the Stateful pod fails, the pod will be rescheduled to another node in the cluster and OpenEBS provides access to the data using the available data copies on other nodes.
    - OpenEBS Replicated Storage is developed with durability and performance as design goals. It efficiently manages the compute (hugepages and cores) and storage (NVMe Drives) to provide fast block storage.

Here we'll see how to use local storage with LVM and OpenEBS. I won't cover the replicated volumes and ZFS because of the resources required to run them.

## Requirements

### Create a LVM PV/VG

If you're not familar, [refer to this LVM documentation](../../../Linux/Filesystems%20And%20Storage/LVM/lvm_working_with_logical_volume_management.md). First, you need to create a LVM PV/VG on your node.

```bash
pvcreate /dev/nvme0n1
vgcreate vg /dev/nvme0n1
```

Here my volume group name is `vg` and my physical volume is `/dev/nvme0n1`.

### Enable Thin Provisioning

Thin provisionning allows you to create logical volumes that are larger than the available extents. Using thin provisioning, you can manage a storage pool of free space, known as a thin pool, which can be allocated to an arbitrary number of devices when needed by applications. You can then create devices that can be bound to the thin pool for later allocation when an application actually writes to the logical volume. The thin pool can be expanded dynamically when needed for cost-effective allocation of storage space.

Before moving forward, you need to enable Thin Provisioning on your kernel:

```bash
modprobe dm_thin_pool
echo dm_thin_pool >> /etc/modules-load.d/dm_thin_pool.conf
```

## Deploy OpenEBS

To deploy OpenEBS, you need to install the OpenEBS operator. First let's create a custom override values file:

=== "openebs-override-values.yaml"

    ```yaml
    openebs-crds:
    csi:
        volumeSnapshots:
        enabled: false
        keep: false

    localpv-provisioner:
    rbac:
        create: true

    zfs-localpv:
    crds:
        zfsLocalPv:
        enabled: false
        csi:
        volumeSnapshots:
            enabled: false

    lvm-localpv:
    crds:
        lvmLocalPv:
        enabled: true
        csi:
        volumeSnapshots:
            enabled: true

    mayastor:
    csi:
        node:
        initContainers:
            enabled: false
    etcd:
        # -- Kubernetes Cluster Domain
        clusterDomain: cluster.local
    localpv-provisioner:
        enabled: false
    crds:
        enabled: false

    engines:
    local:
        lvm:
        enabled: true
        zfs:
        enabled: false
    replicated:
        mayastor:
        enabled: false
    ```

As you can see, everything but LVM is disabled.

Now, let's deploy OpenEBS with Helm:

```bash
helm repo add openebs https://openebs.github.io/openebs
helm repo update
helm upgrade --install openebs ./charts/openebs -n kube-system --wait -f openebs-override-values.yaml
```

You should now have OpenEBS deployed in your cluster:

```bash
$ kubectl get pods -n kube-system -l release=openebs
NAME                                              READY   STATUS        RESTARTS      AGE
openebs-localpv-provisioner-588759f89b-svhpn      1/1     Running       0             30s
openebs-lvm-localpv-controller-84844f9c47-fr5wd   5/5     Running       0             30s
openebs-lvm-localpv-controller-84844f9c47-wgpnp   5/5     Terminating   3 (19h ago)   19h
openebs-lvm-localpv-node-b258g                    2/2     Running       0             32h
openebs-lvm-localpv-node-c2882                    2/2     Running       0             2d3h
openebs-lvm-localpv-node-hbpds                    2/2     Running       0             2d3h
```

## Create an LVM StorageClass

Now, we can create a [StorageClass](https://kubernetes.io/docs/concepts/storage/storage-classes/) that uses the LVM PV/VG we created earlier:

=== "openebs-storageclass.yaml"

    ```yaml
    apiVersion: storage.k8s.io/v1
    kind: StorageClass
    metadata:
        # name of the StorageClass
        name: openebs-lvm
    parameters:
        # storage type
        storage: "lvm"
        # volume group
        volgroup: "vg"
        # enable thin provisioning
        thinProvision: "yes"
    provisioner: local.csi.openebs.io
    # allow volume expansion for later use
    allowVolumeExpansion: true
    # volume binding mode: WaitForFirstConsumer/Immediate
    volumeBindingMode: WaitForFirstConsumer
    # reclaim policy: Delete/Retain
    reclaimPolicy: Delete
    ```

Apply the `StorageClass`:

```bash
kubectl apply -f openebs-storageclass.yaml
```

## Use the LVM StorageClass

You can now use PVCs into your applications! For example:

=== "pvc.yaml"

    ```yaml
    kind: PersistentVolumeClaim
    apiVersion: v1
    metadata:
      name: csi-lvmpv
    spec:
      storageClassName: openebs-lvm
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 4Gi
    ```

=== "deployment.yaml"

    ```yaml
    apiVersion: v1
    kind: Pod
    metadata:
      name: fio
    spec:
     restartPolicy: Never
     containers:
     - name: perfrunner
       image: openebs/tests-fio
       command: ["/bin/bash"]
       args: ["-c", "while true ;do sleep 50; done"]
       volumeMounts:
          - mountPath: /datadir
            name: fio-vol
       tty: true
     volumes:
     - name: fio-vol
       persistentVolumeClaim:
         claimName: csi-lvmpv
    ```

Apply the PVC and the deployment:

```bash
kubectl apply -f pvc.yaml -f deployment.yaml
```

Let's check what we have now! If we look at the PVC, we should see something like this:

```bash
$ kubectl get pvc
NAME        STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   VOLUMEATTRIBUTESCLASS   AGE
csi-lvmpv   Bound    pvc-44e8c69d-b18f-427a-afd3-934fae47ba3e   4Gi        RWO            openebs-lvm    <unset>                 5m16s
```

You can look at the associated volume:

```bash
kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                         STORAGECLASS   VOLUMEATTRIBUTESCLASS   REASON   AGE
pvc-44e8c69d-b18f-427a-afd3-934fae47ba3e   4Gi        RWO            Delete           Bound    default/csi-lvmpv             openebs-lvm    <unset>                          4m3s
```

And finally validate on the host that the LVM volume has been created:

```bash
$ lvs
  LV                                       VG Attr       LSize   Pool        Origin Data%  Meta%  Move Log Cpy%Sync Convert
  pvc-44e8c69d-b18f-427a-afd3-934fae47ba3e vg Vwi-aotz--   4.00g vg_thinpool        3.22
```

## Resources

- [Kubernetes StorageClass](https://kubernetes.io/docs/concepts/storage/storage-classes/)
