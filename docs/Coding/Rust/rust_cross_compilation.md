---
title: "Rust: cross compilation with CI"
slug: rust-cross-compilation/
description: "Cross compile Rust projects with GitHub Actions and Goreleaser."
categories: ["Rust", "CI", "Goreleaser", "Cross Compilation", "GitHub Actions"]
tags: ["Rust", "Cross Compilation", "CI", "Goreleaser", "GitHub Actions"]
---

[Rust](./rust_basics.md) is one of my favorite languages. [Go](https://golang.org/) is another one, and I particularly like it for it's easiness on cross compilation. But Rust is a bit more complex, and I found it a bit hard to cross compile it. So I decided to write a guide on how to cross compile Rust projects with GitHub Actions and Goreleaser.

Why [Goreleaser](https://goreleaser.com/)? Because it's a great tool for cross compilation, it's easy to use, has great features and I didn't find any other tool that was as good as it.

## Prerequisites

- A GitHub account
- A GitHub repository
- A Rust project (a Hello World is enough)

## Repository

Clone your repository. Then let's create a simple Rust project in case you have nothing:

```bash
cargo new myproject
mv myproject/* .
rm -rf myproject
```

Now add the rust toolchain version you want to use for your project. I recommend using the latest stable version (1.92.0 at the time of writing). Create a `rust-toolchain` file in the root of your repository:

=== "rust-toolchain"

    ```toml
    1.92.0
    ```

This will help to keep the same version of Rust for all your CI runs and your local development.

## GitHub Actions

We have to define all architectures, and the OS we want to build on. Then iterate over them to build the binaries. We'll store them as artifacts for later usage (with Goreleaser).

=== ".github/workflows/release.yml"

    {% raw %}
    ```yaml
    name: Release
    
    on:
      push:
        tags:
          - 'v*'
    
    env:
      # Rust build arguments
      BUILD_ARGS: "--release --all-features"
      # The binary name
      BIN_NAME: "myproject"
      # Docker token required to pull images from DockerHub
      DOCKER_LOGIN: ${{ secrets.DOCKER_LOGIN }}
      DOCKER_TOKEN: ${{ secrets.DOCKER_TOKEN }}
    
    jobs:
      build:
        name: Build - ${{ matrix.platform.name }}
        # By default, runs on Ubuntu, otherwise, override with the desired os
        runs-on: ${{ matrix.platform.os || 'ubuntu-24.04' }}
        strategy:
          fail-fast: false
          matrix:
            # Set platforms you want to build your binaries on
            platform:
              # Linux
                # The name is used for pretty print
              - name: Linux x86_64
                # The used Rust target architecture
                target: x86_64-unknown-linux-gnu
              - name: Linux aarch64
                target: aarch64-unknown-linux-gnu
    
              # Mac OS
              - name: MacOS x86_64
                target: x86_64-apple-darwin
              - name: MacOS aarch64
                target: aarch64-apple-darwin
    
              # Windows
              - name: Windows x86_64
                # Use another GitHub action OS
                os: windows-latest
                target: x86_64-pc-windows-msvc
    
        steps:
        - name: Checkout Git repo
          uses: actions/checkout@v4
    
        # Linux & Windows
        - name: Install Rust toolchain
          if: ${{ !contains(matrix.platform.target, 'apple') }}
          uses: actions-rust-lang/setup-rust-toolchain@v1
          with:
            target: ${{ matrix.platform.target }}
            components: rustfmt, clippy
        - name: Build ${{ matrix.platform.name }} binary
          if: ${{ !contains(matrix.platform.target, 'apple') }}
          uses: actions-rs/cargo@v1
          # We use cross-rs if not running on x86_64 architecture on Linux
          with:
            command: build
            use-cross: ${{ !contains(matrix.platform.target, 'x86_64') }}
            args: ${{ env.BUILD_ARGS }} --target ${{ matrix.platform.target }}
    
        # Mac OS
        - name: Login to DockerHub
          if: contains(matrix.platform.target, 'apple')
          # We log on DockerHub
          uses: docker/login-action@v3
          with:
            username: ${{ secrets.DOCKER_LOGIN }}
            password: ${{ secrets.DOCKER_TOKEN }}
        - name: Build ${{ matrix.platform.name }} binary
          if: contains(matrix.platform.target, 'apple')
          # We use a dedicated Rust image containing required Apple libraries to cross-compile on multiple archs
          run: |
            docker run --rm --volume "${PWD}":/root/src --workdir /root/src joseluisq/rust-linux-darwin-builder:2.0.0-beta.1 \
            sh -c "rustup target add ${{ matrix.platform.target }} && \
            if [ \"${{ matrix.platform.target }}\" = \"aarch64-apple-darwin\" ]; then \
              export CC=oa64-clang; \
              export CXX=oa64-clang++; \
            else \
              export CC=o64-clang; \
              export CXX=o64-clang++; \
            fi; \
            cargo build $BUILD_ARGS --target ${{ matrix.platform.target }}"
        
        - name: Store artifact
          uses: actions/upload-artifact@v4
          with:
            # Finally, we store the binary as GitHub artifact for later usage
            name: ${{ matrix.platform.target }}-${{ env.BIN_NAME }}
            path: target/${{ matrix.platform.target }}/release/${{ env.BIN_NAME }}${{ contains(matrix.platform.target, 'windows') && '.exe' || '' }}
            retention-days: 1
    
      release:
        name: Release
        needs: [build]
        # We run the release job only if a tag starts with 'v' letter
        if: startsWith( github.ref, 'refs/tags/v' )
        runs-on: ubuntu-22.04
        steps:
        - name: Checkout Git repo
          uses: actions/checkout@v4
          with:
            fetch-depth: 0
    
        # Download all artifacts
        - uses: actions/download-artifact@v4
          with:
            path: artifacts
    
        # Goreleaser  
        - name: Set up Go
          uses: actions/setup-go@v6
        - name: Run GoReleaser
          uses: goreleaser/goreleaser-action@v6
          with:
            distribution: goreleaser
            version: latest
            # Run goreleaser and ignore non-committed files (downloaded artifacts)
            args: release --clean --skip=validate
          env:
            GITHUB_TOKEN: ${{ secrets.GH_TOKEN }}
            # If you need homebrew to build a formula, you need to use a GitHub token
            # GITHUB_TOKEN: ${{ secrets.GH_TOKEN_GORELEASER }}
    ```
    {% endraw %}

### GitHub Token

Now you have to configure a fine-grained token to allow Goreleaser to update the homebrew-myproject repository (if you are using it).

Your `.goreleaser.yaml` configures a Homebrew Tap in a separate repository (me/homebrew-myproject). The default `secrets.GITHUB_TOKEN` only has permission to modify the current repository. To update the separate `homebrew-myproject` repository automatically, you need a token with permissions across both repositories.

#### Generate a Fine-grained PAT

* Go to GitHub Settings > Developer settings > Personal access tokens > [Fine-grained tokens](https://github.com/settings/personal-access-tokens).
* Name: myproject Release Token (or similar).
* Expiration: Set as desired.
* Resource owner: The owner of the repositories (me).
* Repository access: Select Only select repositories.
* Choose: `myproject` AND `homebrew-myproject`.
* Permissions (Repository permissions):
* Contents: Select Read and Write. (This allows pushing the Release to myproject and the Formula to homebrew-myproject)
* Click Generate token and copy the value.

#### Add Secret to GitHub Actions

* Navigate to your `myproject` repository on GitHub.
* Go to Settings > Secrets and variables > Actions.
* Click New repository secret.
* Name: `GH_TOKEN_GORELEASER` (Must match exactly what is in your YAML).
* Secret: Paste the token you copied.
* Click Add secret.

### DockerHub Token

You also have to configure a token to allow Goreleaser to pull the Docker image from DockerHub and avoid rate limiting.

#### Generate a DockerHub Token

* Go to DockerHub Settings > Security > [Personal Access Tokens](https://hub.docker.com/settings/security).
* Name: myproject Release Token (or similar).
* Expiration: Set as desired.
* Permissions: Select Read and Write.
* Click Generate token and copy the value.

#### Add Secret to GitHub Actions

* Navigate to your `myproject` repository on GitHub.
* Go to Settings > Secrets and variables > Actions.
* Click New repository secret.
* Name: `DOCKERHUB_TOKEN_GORELEASER` (Must match exactly what is in your YAML).
* Secret: Paste the token you copied.
* Click Add secret.

## Goreleaser

Now we have the cross compilation, let's add Goreleaser to our repository. In your repository, create a `goreleaser.yml` file at the root of your repository:

=== "goreleaser.yml"

    {% raw %}
    ```yaml
    archives:
      - formats: 
          - tar.gz
        # use zip for windows
        format_overrides:
          - formats: 
              - zip
            goos: windows
        # this name template makes the release distribution follow a common 
        # convention: ProjectName_Version_OS_Arch
        name_template: >-
          {{ .ProjectName }}_
          {{- .Version }}_
          {{- .Os }}_
          {{- .Arch }}
    
    brews:
      - name: myproject
        description: "My Project description"
        homepage: "https://github.com/xxx/myproject"
        license: "MIT"
        repository:
          owner: me
          name: homebrew-myproject
          token: "{{ .Env.GITHUB_TOKEN }}"
        commit_author:
          name: goreleaserbot
          email: bot@goreleaser.com
    
    builds:
      - id: myproject
        main: goreleaser.go
        binary: myproject
        goos:
          - linux
          - darwin
          - windows
        goarch:
          - amd64
          - arm64
        ignore:
          - goos: windows
            goarch: arm64
        hooks:
          post:
            - ./.goreleaser_hook.sh {{ .Arch }} {{ .Os }} {{ .ProjectName }}
    
    changelog:
      filters:
        exclude:
          - '^docs:'
          - '^test:'
      sort: asc
    
    checksum:
      name_template: 'checksums.txt'
    
    project_name: myproject
    
    snapshot:
      version_template: "{{ incpatch .Version }}-next"
    
    version: 2
    ```
    {% endraw %}

Here you'll be able to generate the releases for your project in those architectures:

- Linux (amd64, arm64)
- MacOS (amd64, arm64)
- Windows (amd64)

Now add a `.goreleaser_hook.sh` file at the root of your repository:

=== ".goreleaser_hook.sh"

    ```bash
    #!/usr/bin/env bash

    # Map amd64->x86_64 and arm64->aarch64
    arch=${1/amd64/x86_64}; arch=${arch/arm64/aarch64}
    name=$3; [[ $2 == "windows" ]] && name+=".exe"

    # Find source binary and destination directory
    src=$(find artifacts -path "*$arch*$2*/$name" -print -quit)
    dst=$(find dist -type d -name "${3}_${2}_${1}*" -print -quit)

    # Copy and chmod if both found
    [[ -f "$src" && -d "$dst" ]] && cp "$src" "$dst/$name" && chmod +x "$dst/$name"
    ```

And finally:

=== "goreleaser.go"

    ```go
    package main

    func main() {
    }
    ```

Now, when you push a tag to your repository, Goreleaser will generate the releases for your project in those architectures.

If you want on concrete example, you can check my [Vessel](https://github.com/deimosfr/vessel) project.