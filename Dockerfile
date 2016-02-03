FROM vungle/kubectl
ADD ./labelgun /usr/bin/labelgun
ENV NODE=$(kubectl describe pod $HOSTNAME | grep Node | awk '{print $2}’) | sed 's@/.*@@'
RUN labelgun
