---
title: "Kubernetes basic usage"
slug: kubernetes_basic_usage/
description: "Kubernetes basic usage guide"
categories: ["Kubernetes"]
tags: ["Kubernetes"]
---

## Introduction

[Kubernetes](https://kubernetes.io/) is an open-source container orchestration platform that automates the deployment, scaling, and management of containerized applications. This guide covers the essential concepts and commands for working with Kubernetes using kubectl, the command-line interface for Kubernetes.

If you want to install Kubernetes, there are many ways to do it. You can install it on a single machine, on a local machine, on a cloud provider, or on a home lab:

- Install it on a home lab with [K3s](./k3s_lightweight_k8s.md)
- Install it on your local machine with [Minikube](https://minikube.sigs.k8s.io/docs/), [Kind](https://kind.sigs.k8s.io/) or [K3s](./k3s_lightweight_k8s.md)
- Install it on a CI with [Kind](https://kind.sigs.k8s.io/)
- Install it on a bare metal/cloud provider with [Kubespray](https://github.com/kubernetes-sigs/kubespray) or [kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)
- Use a managed version with a cloud provider like [EKS](https://aws.amazon.com/eks/), [GKE](https://cloud.google.com/kubernetes-engine), [AKS](https://azure.microsoft.com/en-us/services/kubernetes-service/), [Kapsule](https://www.scaleway.com/en/kubernetes-kapsule/)...

## Kubectl Basics

### Installation and Configuration

Install kubectl using curl (check the architecture of your machine):

```bash
curl -LO https://dl.k8s.io/release/$(curl -Ls https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl
```

Verify the installation:

```bash
kubectl version --client
```

Configure kubectl to use your cluster:

```bash
kubectl config use-context your-cluster-name
```

View current context:

```bash
kubectl config current-context
```

List all contexts:

```bash
kubectl config get-contexts
```

Set namespace for current context:

```bash
kubectl config set-context --current --namespace=my-namespace
```

### Understanding kubectl Output Formats

Default output (wide format):

```bash
kubectl get pods
```

Detailed output in YAML format:

```bash
kubectl get pod <pod-name> -o yaml
```

JSON output format:

```bash
kubectl get pod <pod-name> -o json
```

Custom columns output:

```bash
kubectl get pods -o custom-columns=NAME:.metadata.name,STATUS:.status.phase,IP:.status.podIP
```

Save output to file:

```bash
kubectl get pods -o yaml > pods.yaml
```

### Resource Management Commands

#### Get Resources

List all resources of a specific type:

```bash
kubectl get <resource-type>
```

Common resource type shortcuts:

- pods (po)
- deployments (deploy)
- services (svc)
- configmaps (cm)
- secrets
- nodes (no)
- namespaces (ns)
- persistentvolumeclaims (pvc)
- persistentvolumes (pv)

Get resources with specific labels:

```bash
kubectl get pods -l app=nginx
```

Get resources in a specific namespace:

```bash
kubectl get pods -n my-namespace
```

Watch resources in real-time:

```bash
kubectl get pods -w
```

Get resources with field selectors:

```bash
kubectl get pods --field-selector status.phase=Running
```

#### Describe Resources

Get detailed information about a resource:

```bash
kubectl describe <resource-type> <resource-name>
```

Example commands:

```bash
kubectl describe pod my-pod
kubectl describe deployment my-deployment
kubectl describe service my-service
```

Get events related to a resource:

```bash
kubectl describe pod my-pod | grep -A 10 Events:
```

#### Create and Apply Resources

Create from file:

```bash
kubectl create -f resource.yaml
```

Apply changes (create or update):

```bash
kubectl apply -f resource.yaml
```

Create from URL:

```bash
kubectl create -f https://raw.githubusercontent.com/kubernetes/website/main/content/en/examples/application/nginx-app.yaml
```

Create with client-side dry-run:

```bash
kubectl create -f resource.yaml --dry-run=client
```

Create with server-side validation:

```bash
kubectl create -f resource.yaml --dry-run=server
```

#### Edit and Update Resources

Edit resource in default editor:

```bash
kubectl edit <resource-type> <resource-name>
```

Patch resource:

```bash
kubectl patch deployment my-deployment -p '{"spec":{"replicas":3}}'
```

Replace resource:

```bash
kubectl replace -f resource.yaml
```

Scale resources:

```bash
kubectl scale deployment my-deployment --replicas=3
```

#### Delete Resources

Delete by name:

```bash
kubectl delete <resource-type> <resource-name>
```

Delete by file:

```bash
kubectl delete -f resource.yaml
```

Delete by label:

```bash
kubectl delete pods -l app=nginx
```

Delete with grace period:

```bash
kubectl delete pod my-pod --grace-period=30
```

Force delete:

```bash
kubectl delete pod my-pod --force --grace-period=0
```

### Debugging and Troubleshooting

#### Logs

Get pod logs:

```bash
kubectl logs <pod-name>
```

Get logs from previous instance:

```bash
kubectl logs <pod-name> --previous
```

Follow logs in real-time:

```bash
kubectl logs -f <pod-name>
```

Get logs from specific container:

```bash
kubectl logs <pod-name> -c <container-name>
```

Get logs with timestamps:

```bash
kubectl logs <pod-name> --timestamps
```

Get logs from last N lines:

```bash
kubectl logs <pod-name> --tail=100
```

#### Exec and Port Forward

Execute command in container:

```bash
kubectl exec <pod-name> -- <command>
```

Start interactive shell:

```bash
kubectl exec -it <pod-name> -- /bin/bash
```

Port forward to pod:

```bash
kubectl port-forward <pod-name> 8080:80
```

Port forward to service:

```bash
kubectl port-forward svc/<service-name> 8080:80
```

#### Debugging Tools

Get resource events:

```bash
kubectl get events --sort-by='.lastTimestamp'
```

Get resource metrics:

```bash
kubectl top pods
kubectl top nodes
```

Debug pod issues:

```bash
kubectl debug <pod-name>
```

Copy files to container:

```bash
kubectl cp <pod-name>:/path/to/file ./local-file
```

Copy files from container:

```bash
kubectl cp ./local-file <pod-name>:/path/to/file
```

### Advanced kubectl Usage

#### Resource Queries

Get resources with custom output:

```bash
kubectl get pods -o jsonpath='{.items[*].metadata.name}'
```

Get specific fields:

```bash
kubectl get pod <pod-name> -o jsonpath='{.status.podIP}'
```

Get multiple fields:

```bash
kubectl get pod <pod-name> -o jsonpath='{.metadata.name}{"\t"}{.status.podIP}{"\n"}'
```

Get resources with label selectors:

```bash
kubectl get pods -l 'environment in (production,staging)'
```

#### Configuration Management

View current config:

```bash
kubectl config view
```

View specific context:

```bash
kubectl config view --minify
```

Set namespace for current context:

```bash
kubectl config set-context --current --namespace=my-namespace
```

Create new context:

```bash
kubectl config set-context my-context --cluster=my-cluster --user=my-user --namespace=my-namespace
```

## Standard Kubernetes Objects

Kubernetes provides several built-in objects to manage containerized applications:

### Pods

The smallest deployable unit in Kubernetes.

Create a new pod:

```bash
kubectl run nginx --image=nginx
```

Get pod details:

```bash
kubectl get pods
kubectl describe pod <pod-name>
```

Delete a pod:

```bash
kubectl delete pod <pod-name>
```

### Deployments

Manages the desired state for Pods and ReplicaSets.

Create a new deployment:

```bash
kubectl create deployment nginx --image=nginx
```

Scale a deployment:

```bash
kubectl scale deployment nginx --replicas=3
```

Update a deployment:

```bash
kubectl set image deployment/nginx nginx=nginx:1.19
```

### Services

Exposes applications running on Pods to the network.

Create a new service:

```bash
kubectl expose deployment nginx --port=80 --type=LoadBalancer
```

Get service details:

```bash
kubectl get services
```

### ConfigMaps and Secrets

Store configuration data and sensitive information.

Create a new ConfigMap:

```bash
kubectl create configmap my-config --from-literal=key1=value1
```

Create a new Secret:

```bash
kubectl create secret generic my-secret --from-literal=username=admin
```

### Namespaces

Provide scope for resources and enable resource isolation.

Create a new namespace:

```bash
kubectl create namespace my-namespace
```

List resources in a namespace:

```bash
kubectl get all -n my-namespace
```

### PersistentVolumes and PersistentVolumeClaims

Manage storage resources.

Create a new PVC:

```bash
kubectl create -f pvc.yaml
```

List PVCs:

```bash
kubectl get pvc
```

## Custom Resource Definitions (CRDs)

CRDs extend Kubernetes' API to create custom resources.

### Creating a CRD

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: myresources.example.com
spec:
  group: example.com
  names:
    kind: MyResource
    plural: myresources
    singular: myresource
    shortNames:
      - mr
  scope: Namespaced
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                field1:
                  type: string
```

### Using Custom Resources

Create a custom resource from file:

```bash
kubectl create -f my-resource.yaml
```

List all custom resources:

```bash
kubectl get myresources
```

Get detailed information about a custom resource:

```bash
kubectl describe myresource <name>
```

## Operators and Controllers

### Understanding Operators

Operators are a method of packaging, deploying, and managing a Kubernetes application. They extend the Kubernetes API to create, configure, and manage instances of complex applications like databases, message brokers, etc.

### Working with Operator CRDs

#### Finding Operator CRDs

List all CRDs in the cluster:

```bash
kubectl get crds
```

Get details about a specific CRD:

```bash
kubectl describe crd <crd-name>
```

Get the API version and kind:

```bash
kubectl api-resources | grep <operator-name>
```

#### Understanding CRD Relationships

Find related resources using owner references:

```bash
kubectl get <resource-type> -o jsonpath='{.items[*].metadata.ownerReferences}'
```

Find resources managed by an operator:

```bash
kubectl get all -l app.kubernetes.io/managed-by=<operator-name>
```

Find resources belonging to an instance:

```bash
kubectl get all -l app.kubernetes.io/instance=<instance-name>
```

Check the status of operator-managed resources:

```bash
kubectl get <crd-kind> <name> -o jsonpath='{.status.conditions}'
```

#### Troubleshooting Operators

Get operator pod status:

```bash
kubectl get pods -n <operator-namespace> -l app.kubernetes.io/name=<operator-name>
```

Check operator logs:

```bash
kubectl logs -n <operator-namespace> -l app.kubernetes.io/name=<operator-name>
```

Check CRD status:

```bash
kubectl get crd <crd-name> -o jsonpath='{.status.conditions}'
```

Verify CRD schema:

```bash
kubectl get crd <crd-name> -o yaml
```

List all resources managed by the operator:

```bash
kubectl get all -l app.kubernetes.io/managed-by=<operator-name>
```

Check for failed resources:

```bash
kubectl get <crd-kind> -o jsonpath='{.items[?(@.status.conditions[0].status=="False")].metadata.name}'
```
