[![](http://badge-imagelayers.iron.io/vungle/labelgun:latest.svg)](http://imagelayers.iron.io/?images=vungle/labelgun:latest 'Get your own badge on imagelayers.iron.io')

# labelgun
Insert AWS metadata as Kubernetes Labels

### Supported:
* ec2tags
* availability zone
* instance type

# Configure

Edit the labelgun.yml with approriate Environment Variable valuess for `KUBE_MASTER`, `AWS_REGION` and `LABELGUN_INTERVAL` in seconds

# Launch the DaemonSet

`kubectl create -f labelgun.yml`

Note: this requries you have DaemonSets enabled https://github.com/kubernetes/kubernetes/blob/master/docs/design/daemon.md
