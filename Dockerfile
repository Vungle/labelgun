FROM vungle/kubectl
ADD bin/labelgun-linux-amd64 /usr/bin/labelgun
CMD "labelgun"
