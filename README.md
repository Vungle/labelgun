# labelgun
Insert AWS metadata as Kubernetes Labels

# Configure

Edit the labelgun.yml with approriate Environment Variable valuess for `KUBE_MASTER`, `AWS_REGION` and `LABELGUN_INTERVAL` in seconds

# Launch the DaemonSet

`kubectl create -f labelgun.yml`

Note: this requries you have DaemonSets enabled https://github.com/kubernetes/kubernetes/blob/master/docs/design/daemon.md
