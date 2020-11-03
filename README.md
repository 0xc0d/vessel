# vessel
A tiny tool for managing OCI containers written in Go.
It basically is a tiny version of docker, but without using containerd or runc.

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

## Example

run `/bin/sh` in `alpine:latest`

    vessel run alpine /bin/sh
    vessel run alpine # same as above due to alpine default command