[![](https://images.microbadger.com/badges/image/dailyhotel/labelgun.svg)](https://microbadger.com/images/dailyhotel/labelgun "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/dailyhotel/labelgun.svg)](https://microbadger.com/images/dailyhotel/labelgun "Get your own version badge on microbadger.com")
[![Build Status](https://travis-ci.org/DailyHotel/labelgun.svg?branch=master)](https://travis-ci.org/DailyHotel/labelgun)

# labelgun

**Insert AWS EC2 Tags as Kubernetes Node Labels.** 

This is the improved version of [Vungle/labelgun](https://github.com/Vungle/labelgun) in several aspects:

* [DaemonSet](https://github.com/kubernetes/kubernetes/blob/master/docs/design/daemon.md) is not required. Just launch a single pod and save the rest of your computational resources.
* Kubernetes version v1.5.x is supported
* Fine-grained logging
* Private base image Vungle/kubectl is removed
* Better developer support using `Makefile` and `glide.yaml`

## Supported:

* ec2tags
* ~~availability zone~~ and ~~instance type~~ are not supported any more since [Kubernetes itself provides the both](https://kubernetes.io/docs/admin/multiple-zones/).

## Configure

Edit the `labelgun.yml` with appropriate Environment Variable values for [`LABELGUN_ERR_THRESHOLD`](https://godoc.org/github.com/golang/glog) and `LABELGUN_INTERVAL` in seconds.

## Launch the DaemonSet

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: labelgun
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: labelgun 
    spec:
      containers:
        - image: dailyhotel/labelgun:latest
          imagePullPolicy: Always
          name: labelgun
          env:
            - name: LABELGUN_ERR_THRESHOLD
              value: "INFO"
            - name: LABELGUN_INTERVAL
              value: "60"
```

`kubectl create -f labelgun.yml`

## Develop

``` bash
go get github.com/dailyhotel/labelgun

glide install --strip-vendor --strip-vcs

make
```
