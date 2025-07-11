# [WIP] Kubernetes Lab
A full local Kubernetes setup with monitoring, logging and demo apps - Using k3d + Tilt

> Note: this is a work in progress

## Requirements

1. [Install Docker](https://docs.docker.com/engine/install/)
2. [Setup Docker with non-root user](https://docs.docker.com/engine/install/linux-postinstall/)
3. [Install k3d](https://k3d.io/stable/#installation)
4. [Install kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
5. [Install Tilt](https://docs.tilt.dev/)


## Setup

### Setup Demo apps in your hosts file

```text
127.0.0.1 nginx-hello.kubelab.lan
127.0.0.1 argus-mcp.kubelab.lan
127.0.0.1 prometheus.kubelab.lan
```

* **On macOS/Linux:** Edit `/etc/hosts` (e.g., `sudo nano /etc/hosts`).
* **On Windows:** Open Notepad **as an Administrator** and edit `C:\Windows\System32\drivers\etc\hosts`.



## Start

Run `./start-lab.sh` to start everything

## Stop / Start

After the initial setup, you can use `tilt up` or `tilt down`


## Access

### Tilt UI

You can access the Tilt UI [HERE](http://localhost:10350/)


### The deployed pods

Based on the hosts entry you have added, you can access them like this

`http://nginx-hello.kubelab.lan:8081/`

> Note: The port number 8081 comes from the fact that we set it up to use port 8081
> during k3d installation. Check the script `start-lab.sh`