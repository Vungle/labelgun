FROM vungle/kubectl
RUN apk update && apk add curl
ADD bin/labelgun-linux-amd64 /usr/bin/labelgun
CMD "labelgun"
