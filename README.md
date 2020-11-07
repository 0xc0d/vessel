# vessel
A tiny educational-purpose project to create containers, written in Go.

It basically is a tiny version of docker, it uses neither [containerd](https://containerd.io/) nor [runc](https://github.com/opencontainers/runc). Only a set of the Linux features.

## Features
Vessel supports:
* __Control Groups__ for resource restriction (CPU, Memory, Swap, PIDs)
* __Namespace__ for global system resources isolation (Mount, UTS, Network, IPS, PID)
* __Union File System__ for branches to be overlaid in a single coherent file system. (OverlayFS)

## Read more
Here is the list of blog posts I've written about vessel:

1. [Build Containers From Scratch in Go (Part 1: Namespaces)](https://alijosie.medium.com/build-containers-from-scratch-in-go-part-1-namespaces-c07d2291038b)
2. To be continued...

## Install

    go get -u github.com/0xc0d/vessel
    
## Usage

    Usage:
      vessel [command]
    
    Available Commands:
      exec        Run a command inside a existing Container.
      help        Help about any command
      images      List local images
      ps          List Containers
      run         Run a command inside a new Container.

## Examples

Run `/bin/sh` in `alpine:latest`

    vessel run alpine /bin/sh
    vessel run alpine # same as above due to alpine default command

Restart Nginx service inside a container with ID: 123456789123

    vessel exec 1234567879123 systemctrl restart nginx
    
List running containers

    vessel ps
    
List local images

    vessel images

## Notice
vessel, obviously, is not a production ready container manager tool. 
