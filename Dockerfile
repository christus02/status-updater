FROM golang:1.15
COPY ./src/status-patcher.go ./
RUN go get "k8s.io/klog"
RUN go get "k8s.io/api/core/v1"
RUN go get "k8s.io/apimachinery/pkg/util/runtime"
RUN go get "k8s.io/client-go/informers"
RUN go get "k8s.io/client-go/kubernetes"
RUN go get "k8s.io/client-go/tools/cache"
RUN go get "k8s.io/client-go/rest"
RUN go get "k8s.io/client-go/tools/clientcmd"
RUN go get "k8s.io/apimachinery/pkg/apis/meta/v1"
RUN go get "k8s.io/apimachinery/pkg/types"

RUN go build -o status-updater
RUN rm status-patcher.go

CMD ["./status-updater"]
