[![](https://images.microbadger.com/badges/image/dailyhotel/labelgun.svg)](https://microbadger.com/images/dailyhotel/labelgun "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/dailyhotel/labelgun.svg)](https://microbadger.com/images/dailyhotel/labelgun "Get your own version badge on microbadger.com")

# labelgun
Insert AWS metadata as Kubernetes Labels

### Supported:
* ec2tags
* availability zone
* instance type

# Configure

Edit the `labelgun.yml` with appropriate Environment Variable values for [`LABELGUN_ERR_THRESHOLD`](https://godoc.org/github.com/golang/glog) and `LABELGUN_INTERVAL` in seconds.

# Launch the DaemonSet

```yaml
apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  namespace: kube-system
  name: labelgun
spec:
  template:
    metadata:
      labels:
        app: labelgun
      name: labelgun
    spec:
      containers:
      - image: dailyhotel/labelgun
        name: labelgun
        env:
          - name: LABELGUN_ERR_THRESHOLD
            value: "ERROR"
```

`kubectl create -f labelgun.yml`

Note: this requires you have DaemonSets enabled https://github.com/kubernetes/kubernetes/blob/master/docs/design/daemon.md

