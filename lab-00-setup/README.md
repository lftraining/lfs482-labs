# Lab 00: Setup

# Introduction

This preliminary lab will download the required binaries and container images for the other labs in this course and 
improve the attendees experience, especially for those on limited or slow network connection.

The required binaries will be stored in the [bin](../bin) directory. You can add these to the PATH environment variable
by running the following command, which will work from all the lab directories and ensure the version that the labs have
been tested against will be used.

```shell
export PATH=../bin:$PATH
```

A number of the labs using [Kind](https://kind.sigs.k8s.io/) to provide a Kubernetes cluster to run the exercises on.
Each lab is configured to load the required container images into the Kind cluster when the cluster is launched at the
start of each lab. This prevents attendees from having to download the required container images in the cluster every
time; improving launch times for the Pods and reducing the network usage.

# Download the Binaries and Container Images

From the [lab-00-setup](.) directory, run the following command.

```shell
make all
```

Depending on your network speeds, this will take a few a minutes. While you wait, you can have a look at the
[Course Overview](../README.md) to familiarise yourself with the various lab exercises you will be completing.
