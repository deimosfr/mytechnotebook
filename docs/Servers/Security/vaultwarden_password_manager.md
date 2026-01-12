---
title: "Vaultwarden: The great Password Manager"
slug: vaultwarden-password-manager/
description: "Deploying Vaultwarden, a lightweight password manager compatible with Bitwarden clients."
categories: ["Servers", "Security", "Password Manager"]
tags: ["Vaultwarden", "Bitwarden", "Password Manager"]
---

![Vaultwarden](../../static/images/vaultwarden_logo.avif){width=600}

[Vaultwarden](https://github.com/dani-garcia/vaultwarden) is an unofficial Bitwarden server implementation written in [Rust](../../Coding/Rust/index.md). It's compatible with the official [Bitwarden](https://bitwarden.com/) clients and is perfect for self-hosted environments where the official resource-heavy containers might be overkill. If you consider using a SaaS password manager, then look at [Bitwarden](https://bitwarden.com/).

We'll deploy Vaultwarden using the Gissilabs Helm chart, configured with a PostgreSQL database, SMTP for emails, and secured with network policies.

## Preparation

First, let's add the Gissilabs Helm repository:

```bash
helm repo add gissilabs https://gissilabs.github.io/charts/
helm repo update
```

### Database Connection

!!! note

    By default Vaultwarden uses SQLite as a database. However, we are using an external PostgreSQL database as it is more reliable and provides better performance.

We are using an external [PostgreSQL](../Databases/PostgreSQL/cloud_native_postgresql_operator.md) database. We need to create a Secret to hold the connection string.

Create a file named `secret.yaml`:

=== "secret.yaml"

    ```yaml
    apiVersion: v1
    kind: Secret
    metadata:
      name: vaultwarden-database
    type: Opaque
    stringData:
      # Replace with your actual database URL
      database-url: postgres://user:password@host:5432/dbname
    ```

Update the `database-url` value with your actual database connection string.

Apply it to your cluster:

```bash
kubectl apply -f secret.yaml
```

## Configuration

Now, let's configure the Vaultwarden instance. We'll set up the domain, database connection, SMTP settings, and resource limits.

Create a `values-overrides.yaml` file with the following configuration:

=== "values-overrides.yaml"

    ```yaml
    # Database Configuration
    database:
      type: postgresql
      existingSecret: "vaultwarden-database"
      existingSecretKey: "database-url"
      maxConnections: 20
      retries: 15
    
    vaultwarden:
      # Domain should be https absolutely
      domain: "https://vaultwarden.mydomain.com"
      # Registration
      allowSignups: false
      signupDomains: []
      verifySignup: true
      requireEmail: true
      emailAttempts: 3
      emailTokenExpiration: 600
      allowInvitation: true
      invitationExpiration: 120
      passwordHintsAllowed: true
      showPasswordHint: false
      enableWebVault: true
      # send feature
      enableSends: true
      orgCreationUsers: all
      attachmentLimitOrg: ""
      attachmentLimitUser: "1000"
      sendLimitUser: "1000"
      hibpApiKey: ""
      autoDeleteDays: "7"
      orgEvents: false
      orgEventsRetention: ""
      emailChangeAllowed: true
      extraEnv: {}
   
      # admin feature available on `/admin` path
      admin:
        enabled: false
        disableAdminToken: false
        # Generate a token: echo -n "MySecretPassword" | argon2 "$(openssl rand -base64 32)" -e -id -k 65540 -t 3 -p 4
        token: "set a very long token here"
    
      # emergency access
      emergency:
        enabled: true
        reminder: "0 3 * * * *"
        timeout: "0 7 * * * *"
   
      # Optional: smtp configuration. Gmail used here
      # More info: https://support.google.com/accounts/answer/185833?hl=en&ref_topic=7189145
      smtp:
        enabled: true
        host: "smtp.gmail.com"
        from: "my email address"
        fromName: "Vaultwarden"
        security: starttls
        port: "587"
        authMechanism: Plain
        heloName: ""
        timeout: 15
        invalidHostname: false
        invalidCertificate: false
        user: "my email address"
        password: "my email password"
        embedImages: true
   
      # enable physical security keys
      yubico:
        enabled: false
    
      log:
        level: "info"
    
      icons:
        service: google
        disableDownload: false
        cache: "2592000"
        cacheFailed: "259200"
        redirectCode: 302
   
      # enable mobile push notifications
      # Create here key and id: https://bitwarden.com/host/
      push:
        enabled: true
        installationId: "your installation id"
        installationKey: "your installation key"
        relayUri: "https://api.bitwarden.eu"
        identityUri: "https://identity.bitwarden.eu/"
    
    service:
      type: ClusterIP
      httpPort: 80
      externalTrafficPolicy: Cluster
    
    ingress:
      enabled: false
    
    ingressRoute:
      enabled: false
    
    persistence:
      enabled: true
      size: 1Gi
      accessMode: ReadWriteOnce
      storageClass: "openebs-lvm" # define here the storage class you want to use
    
    replicaCount: 1
    
    # resources limitation
    resources:
      requests:
        memory: "256Mi"
        cpu: "10m"
      limits:
        memory: "256Mi"
        cpu: "500m"
    
    # As we're using a single replica, we need to set the strategy to Recreate because of the PVC
    strategy:
      type: Recreate
    ```

We can now deploy the Vaultwarden instance:

```bash
helm install vaultwarden gissilabs/vaultwarden -f values-overrides.yaml
```

## Exposing the Service

Finally, we'll expose Vaultwarden using the [Gateway API](../Containers/Kubernetes/gateway-api.md). This `HTTPRoute` directs traffic from `vaultwarden.mydomain.com` to our service.

=== "httproute.yaml"

    ```yaml
    apiVersion: gateway.networking.k8s.io/v1
    kind: HTTPRoute
    metadata:
      name: vaultwarden
    spec:
      hostnames:
      - vaultwarden.mydomain.com
      parentRefs:
      - name: mygateway
        namespace: gateway
      rules:
      - backendRefs:
        - kind: Service
          name: vaultwarden
          port: 80
        matches:
        - path:
            type: PathPrefix
            value: / 
    ```

## Security


### Configuration update

Once you've deployed the service and conencted with your first use, it's recommended to:

1. Update the configuration to disable the registration
2. Enable the emergency access
3. Disable the admin feature

## Optional: Network Security

Because Vaultwarden is better when accessible online, it's better to limit the pod from being able to access anything in your cluster but the minimum. We'll use a `CiliumNetworkPolicy` to restrict traffic. This policy allows:

*   **Egress** to CoreDNS (UDP 53).
*   **Egress** to the PostgreSQL cluster (TCP 5432).
*   **Egress** to the world (necessary for SMTP, mobile push notifications, and icon fetching).

=== "network-policy.yaml"

    ```yaml
    apiVersion: cilium.io/v2
    kind: CiliumNetworkPolicy
    metadata:
      name: vaultwarden
    spec:
      endpointSelector:
        matchLabels:
          app.kubernetes.io/name: vaultwarden
      egress:
        # allow dns queries to coredns
        - toEndpoints:
            - matchLabels:
                io.kubernetes.pod.namespace: kube-system
                k8s-app: kube-dns
          toPorts:
            - ports:
                - port: "53"
                  protocol: UDP
              rules:
                dns:
                  - matchPattern: "*"
        # allow access to the postgresql cluster
        - toEndpoints:
            - matchLabels:
                cnpg.io/cluster: pg-cluster
                io.kubernetes.pod.namespace: databases
              matchExpressions:
                - key: cnpg.io/cluster
                  operator: In
                  values:
                    - pg-cluster
          toPorts:
            - ports:
                - port: "5432"
        # allow access to the world for smtp, mobile push notifications, and icon fetching
        - toEntities:
            - world
    ```

### Cloudflared

Being able to access the service from the internet is a must. And being more protected is also a must. You can use [Cloudflared](../Security/VPN%20and%20Tunnels/cloudflared.md) to expose the service.

