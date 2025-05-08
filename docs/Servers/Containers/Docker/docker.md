---
weight: 999
title: "Docker: containerize your applications"
description: "A comprehensive guide to Docker containers"
toc: true
---

## Introduction

Hey there! Let's talk about Docker - it's this cool open-source platform that makes deploying apps a breeze through something called containerization. Think of containers as neat little packages that bundle up your application with everything it needs - the code, runtime, tools, and libraries. This way, your app runs the same way no matter where you deploy it.

What makes Docker special? Unlike traditional virtual machines that simulate an entire operating system, Docker containers share the host's kernel and only isolate the app processes. This makes them super lightweight and quick to start up - perfect when you need to get things running fast!

Here's why so many developers love Docker:

- **Consistency**: Your app works the same way everywhere - from your laptop to the production server
- **Isolation**: Apps and their dependencies stay in their own sandbox, away from your host system
- **Portability**: Got Docker? Then you can run the containers - on Mac, Windows, Linux, in the cloud, wherever!
- **Efficiency**: Containers are lightweight and share the host kernel, so they use way fewer resources than VMs
- **Scalability**: Need more power? Just spin up more containers - it's that simple
- **Versioning**: Made a mistake? No problem! You can track changes and roll back to previous versions

## Installation

{{< tabs tabTotal="2">}}
{{% tab tabName="Debian" %}}

Update your package index:

Install required packages to allow apt to use a repository over HTTPS:

```bash
sudo apt-get update && apt-get install -y ca-certificates curl gnupg lsb-release
```

Add Docker's official GPG key:

```bash
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
```

Set up the repository:

```bash
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

Update apt again and install Docker:

```bash
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
```

{{% /tab %}}
{{% tab tabName="RedHat" %}}

Install required packages:

```bash
sudo yum install -y yum-utils
```

Add the Docker repository:

```bash
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
```

Install Docker:

```bash
sudo yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
```

Start and enable Docker service:

```bash
sudo systemctl start docker
sudo systemctl enable docker
```

{{% /tab %}}
{{< /tabs >}}

Verify the installation:

```bash
sudo docker run hello-world
```

## Configuration

### Memory Management

For optimal container performance, especially when running memory-intensive applications, you'll need to enable proper memory management by modifying your boot parameters:

```bash
# Edit GRUB configuration
sudo nano /etc/default/grub

# Add these parameters to GRUB_CMDLINE_LINUX_DEFAULT
GRUB_CMDLINE_LINUX_DEFAULT="quiet cgroup_enable=memory swapaccount=1"

# Update GRUB and reboot
sudo update-grub
sudo reboot
```

### Storage Location

By default, Docker stores its images and containers in `/var/lib/docker`. If you want to change this location (for example, to a larger drive), you can modify Docker's configuration:

```bash
# Edit Docker daemon configuration
sudo nano /etc/docker/daemon.json

# Add or modify the storage-driver and data-root settings
{
  "storage-driver": "overlay2",
  "data-root": "/path/to/your/docker/data"
}

# Restart Docker service
sudo systemctl restart docker
```

Alternatively, if you're using Docker installed via package manager, you can edit:

```bash
sudo nano /etc/default/docker

# Add or modify DOCKER_OPTS
DOCKER_OPTS="-g /path/to/your/docker/data"

# Restart Docker
sudo systemctl restart docker
```

## Usage

### Basic Commands

#### Pulling Images

Download Docker images from Docker Hub:

```bash
docker pull ubuntu:latest
docker pull nginx:1.21
```

#### List Images

View all downloaded images:

```bash
docker images
```

#### Run Containers

Start a container from an image:

```bash
docker run -it ubuntu bash    # Interactive terminal
docker run -d nginx           # Detached mode
```

#### List Containers

View running containers:

```bash
docker ps           # Running containers
docker ps -a        # All containers (including stopped)
```

#### Stop and Remove Containers

```bash
docker stop <container_id>
docker rm <container_id>
```

### Expose Ports

Map container ports to host ports:

```bash
docker run -d -p 8080:80 nginx
```

This maps port 80 in the container to port 8080 on the host.

### Mount Volumes

Share data between host and container:

```bash
docker run -d -v /host/path:/container/path nginx
```

For named volumes:

```bash
docker volume create my_volume
docker run -d -v my_volume:/container/path nginx
```

### Network Configuration

#### Create a Network

```bash
docker network create my_network
```

#### Connect Containers

```bash
docker run -d --name db --network my_network mysql
docker run -d --name web --network my_network nginx
```

### Build Custom Images

1. Create a Dockerfile:

```dockerfile
FROM ubuntu:20.04
RUN apt-get update && apt-get install -y nginx
COPY ./app /var/www/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

2. Build the image:

```bash
docker build -t my-app:1.0 .
```

### Advanced Container Management

#### Committing Changes

After making changes to a container, you can create a new image that preserves those changes:

```bash
# Make changes to a container
docker exec -it <container_id> bash
# ... make your changes ...

# Commit changes to a new image
docker commit <container_id> username/my-custom-image:tag

# View your new image
docker images
```

This is particularly useful for creating custom images based on your modifications.

#### Viewing Container Differences

To see what files have been changed, added, or deleted in a container compared to its base image:

```bash
docker diff <container_id>
```

The output uses these prefixes:

- `A`: Added
- `D`: Deleted
- `C`: Changed

This is extremely helpful for troubleshooting and understanding what modifications have occurred inside a container.

### Container Management

#### View Container Logs

```bash
docker logs <container_id>
```

#### Execute Commands in Running Containers

```bash
docker exec -it <container_id> bash
```

#### Check Container Resource Usage

```bash
docker stats <container_id>
```

## References

- [Docker Official Documentation](https://docs.docker.com/)
- [Docker Hub](https://hub.docker.com/)
