---
title: "SOPS: Mastering Secrets inside git"
slug: sops-mastering-secrets-inside-git/
description: "A comprehensive guide to managing secrets securely using SOPS."
categories: ["Linux", "Security", "DevOps"]
tags:
  ["sops", "security", "encryption", "gitops", "devops"]
---

## Why Managing Secrets is Hard

In the modern DevOps landscape, "Secrets Sprawl" is a real issue. We have API keys, database credentials, and certificates that need to be used by our applications, but storing them safely is a constant struggle.

- **Environment Variables**: Hard to manage at scale.
- **Git Repositories**: Dangerous if secrets are committed in plain text.
- **Dedicated Vaults** (like HashiCorp Vault): Powerful but complex to maintain.

[SOPS](https://github.com/getsops/sops) (Secrets OPerationS) is a very good solution made by Mozilla.

## What is SOPS?

SOPS is an open-source tool, originally by Mozilla and now a CNCF project, that acts as an encrypted file editor. It solves the problem of "secrets in git" by encrypting the *values* of your configuration files while keeping the *keys* readable.

This means a file like this:

```yaml
database:
  host: db.example.com
  password: SECRET_PASSWORD
```

Becomes this in your repo:

```yaml
database:
  host: db.example.com
  password: ENC[AES256_GCM,data:...]
```

You can still see the structure (`database.host`), but the sensitive `password` is secure. When you need to edit it, SOPS handles the decryption and re-encryption transparently.

## Getting Started

### Installation

=== "MacOS"

    ```bash
    brew install sops
    ```

=== "Arch Linux"

    ```bash
    sudo pacman -S sops
    ```

=== "Linux"

    Download the latest release from the [GitHub repository](https://github.com/getsops/sops/releases).

    ```bash
    wget https://github.com/getsops/sops/releases/download/v3.7.3/sops-v3.7.3.linux.amd64
    sudo mv sops-v3.7.3.linux.amd64 /usr/local/bin/sops
    sudo chmod +x /usr/local/bin/sops
    ```

### Concepts: Key Management Services (KMS)

SOPS doesn't manage keys itself; it relies on external Key Management Services. It supports:

- **age**: A modern, simple file encryption tool (becoming the preferred default for many).
- **PGP/GPG**: Great for personal use or small teams.
- **AWS KMS**: Ideal for AWS-centric infrastructure.
- **Google Cloud KMS** & **Azure Key Vault**: For their respective clouds.

## Using SOPS with age

Let's look at a robust workflow using **age**.

### Install age

You'll need the `age` tool to generate keys.

=== "MacOS"

    ```bash
    brew install age
    ```

=== "Debian"

    ```bash
    sudo apt install age
    ```

=== "Arch Linux"

    ```bash
    sudo pacman -S age
    ```

### Generate a private Key

Generate a new key pair and store it in a key file.

```bash
mkdir -p ~/.age
chmod 700 ~/.age
age-keygen -o ~/.age/keys.txt
```

This will create a key file and output your **public key** (starting with `age1...`).  The public key is present as comment in the private key file.

### Configure SOPS

While you can pass flags to every command, using a `.sops.yaml` configuration file is best practice. Create this file in the root of your project:

```yaml
creation_rules:
  # Encrypt all files ending in .enc.yaml with our age public key
  - path_regex: .*\.enc\.yaml$
    age: "age1..." # Replace with your actual public key
```

### Create an Encrypted File

Now, simply create a new file matching your regex. SOPS will automatically look up the config.

Create a file called `common_secrets.yaml` with the following content:

```yaml
api_key: "hidden-value"
debug_mode: true
```

Then you can encrypt it with:

```bash
sops -e --age=age1... common_secrets.yaml > common_secrets.enc.yaml
```

Now you can remove the non encrypted file and git commit your encrypted file:

```bash
rm common_secrets.yaml
```

### Editing and Viewing

To edit the file again, just run:
```bash
sops secrets.enc.yaml
```
SOPS decrypts it in memory, opens your editor, and re-encrypts on save.

To output the decrypted content to stdout:
```bash
sops -d secrets.enc.yaml
```

## Advanced Features

### Key Groups and Redundancy

SOPS allows you to use multiple keys for the same file. This is critical for redundancy. For example, you could encrypt a file with both an AWS KMS key *and* a backup Age key. If AWS is down or you lose access, the Age key can still decrypt the data.

```yaml
creation_rules:
  - path_regex: .*\.yaml$
    key_groups:
      - 
        age:
          - "age1primarykey..."
        kms:
          - arn: "arn:aws:kms:region:account:key/123"
```


### GitDiff Integration

Since the encrypted files change every time they are saved (due to new nounces/IVs), git diffs can be noisy. You can configure git to use sops to diff the *decrypted* contents:

```bash
# .gitattributes
*.enc.yaml diff=sops

# Global git config
git config --global diff.sops.textconv "sops -d"
```

Now `git diff` shows you the actual changes to your data, not just scrambled ciphertext!

## DevOps Integration

### In CI/CD

In your CI pipeline (e.g., GitHub Actions, GitLab CI), you need to provide the `keys.txt` content to the worker.

1.  Store the content of your `keys.txt` as a secret in your CI system (e.g., `SOPS_AGE_KEY`).
2.  Set the environment variable `SOPS_AGE_KEY` in your pipeline job.

You can then run:


```bash
sops -d secrets.enc.yaml > secrets.yaml
```
And load `secrets.yaml` into your application.

### With Kubernetes

Tools like **FluxCD** have native SOPS integration. You can commit valid Kubernetes Secrets encrypted with SOPS to your git repo. The Flux controller in the cluster, having access to the decryption key, will apply them as standard (decrypted) Secrets inside the cluster.

You can also look at [Sops Operator](https://github.com/isindir/sops-secrets-operator) for a more Kubernetes-native approach.

## Conclusion

SOPS strikes a perfect balance between security and usability. It fits naturally into the "Infrastructure as Code" philosophy, allowing you to version control your secrets safely without the headache of managing a complex Vault infrastructure for every project.
