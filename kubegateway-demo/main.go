package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/endpoints/responsewriter"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"kubegatewaydemo/clusters"
	"kubegatewaydemo/proxy"
	"kubegatewaydemo/request"
	"kubegatewaydemo/transport"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	retryAfter = 1
)

type dispatcher struct {
	codecs          serializer.CodecFactory
	enableAccessLog bool
}

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	_, config := GetClientSet()
	http2configCopy := *config
	http2configCopy.WrapTransport = transport.NewDynamicImpersonatingRoundTripper
	ts, err := rest.TransportFor(&http2configCopy)
	if err != nil {
		panic(err)
	}
	urrt, ok := clusters.UnwrapUpgradeRequestRoundTripper(ts)
	if !ok {
		klog.Errorf("failed to convert transport to proxy.UpgradeRequestRoundTripper for <cluster:%s,endpoint:%s>", "tjkj", http2configCopy.Host)
	}
	endpoint := &clusters.EndpointInfo{
		Cluster:               "tjkj",
		Endpoint:              http2configCopy.Host,
		ProxyTransport:        ts,
		PorxyUpgradeTransport: urrt,
		Mutex:                 sync.Mutex{},
	}

	user, ok := genericapirequest.UserFrom(ctx)
	if !ok {
		d.responseError(errors.NewInternalError(fmt.Errorf("no user info found in request context")), w, req, proxy.StatusReasonInvalidRequestContext)
		return
	}
	requestInfo, ok := genericapirequest.RequestInfoFrom(ctx)
	if !ok {
		d.responseError(errors.NewInternalError(fmt.Errorf("no request info found in request context")), w, req, proxy.StatusReasonInvalidRequestContext)
		return
	}

	ep := &url.URL{
		Scheme: "https",
		Host:   "tjyw-k8s-api:6443",
		Path:   "/",
	}

	// mark this proxy request forwarded
	if err := request.SetProxyForwarded(req.Context(), endpoint.Endpoint); err != nil {
		d.responseError(errors.NewInternalError(err), w, req, proxy.StatusReasonInvalidRequestContext)
		return
	}

	location := &url.URL{}
	location.Scheme = ep.Scheme
	location.Host = ep.Host
	location.Path = req.URL.Path
	location.RawQuery = req.URL.Query().Encode()

	newReq, _ := newRequestForProxy(location, req, "")
	// close this request if endpoint is stoped
	go func() {
		select {
		case <-newReq.Context().Done():
			// this context comes from incoming server requests, and then we use
			// it as proxy client context to control cancellation
			//
			// For incoming server requests, the context is canceled when the
			// client's connection closes, the request is canceled (with HTTP/2),
			// or when the ServeHTTP method returns.
		}
	}()

	genericapirequest.RequestInfoFrom(ctx)

	delegate := proxy.DecorateResponseWriter(req, w, true, requestInfo, "hostname", endpoint.Endpoint, user, nil)
	delegate.MonitorBeforeProxy()
	defer delegate.MonitorAfterProxy()

	rw := responsewriter.WrapForHTTP1Or2(delegate)

	proxyHandler := proxy.NewUpgradeAwareHandler(location, endpoint.ProxyTransport, endpoint.PorxyUpgradeTransport, false, false, d, endpoint)
	proxyHandler.ServeHTTP(rw, newReq)
}

func (d *dispatcher) responseError(err *errors.StatusError, w http.ResponseWriter, req *http.Request, reason string) {
	gv := schema.GroupVersion{Group: "", Version: "v1"}

	switch {
	case errors.IsTooManyRequests(err), utilnet.IsProbableEOF(err):
		w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
	case errors.IsServiceUnavailable(err):
		w.Header().Set("Retry-After", strconv.Itoa(retryAfter*30))
	}

	code := int(err.Status().Code)
	if proxy.CaptureErrorReason(reason) {
		var urlHost string
		if req.URL != nil {
			// url.Host is different from req.Host when caller is reverse proxy.
			// we need this host to determine which endpoint it is if possible.
			urlHost = req.URL.Host
		}
		klog.Errorf("[proxy termination] method=%q host=%q uri=%q url.host=%v resp=%v reason=%q message=[%v]", req.Method, HostWithoutPort(req.Host), req.RequestURI, urlHost, code, reason, err.Error())
	}

	runtime.Must(request.SetProxyTerminated(req.Context(), reason))

	responsewriters.ErrorNegotiated(err, d.codecs, gv, w, req)
}

func HostWithoutPort(hostport string) string {
	hostport = strings.ToLower(hostport)
	noport, _, err := net.SplitHostPort(hostport)
	if err != nil {
		return hostport
	}
	return noport
}

// newRequestForProxy returns a shallow copy of the original request with a context that may include a timeout for discovery requests
func newRequestForProxy(location *url.URL, req *http.Request, _ string) (*http.Request, context.CancelFunc) {
	ctx := req.Context()
	newCtx, cancel := context.WithCancel(ctx)

	// WithContext creates a shallow clone of the request with the same context.
	newReq := req.WithContext(newCtx)
	newReq.Header = utilnet.CloneHeader(req.Header)
	newReq.URL = location

	return newReq, cancel
}

// implements k8s.io/apimachinery/pkg/util/proxy.ErrorResponder interface
func (d *dispatcher) Error(w http.ResponseWriter, req *http.Request, err error) {
	status := proxy.ErrorToProxyStatus(err)
	reason := proxy.StatusReasonUpgradeAwareHandlerError
	if status.Code == http.StatusBadGateway {
		reason = proxy.StatusReasonReverseProxyError
	}
	d.responseError(&errors.StatusError{ErrStatus: *status}, w, req, reason)
}

func main() {
	d := &dispatcher{
		enableAccessLog: true,
	}
	log.Fatal(http.ListenAndServe(":8082", d))
}

func GetClientSet() (*kubernetes.Clientset, *rest.Config) {
	var err error
	var config *rest.Config
	var kubeConfig *string
	if home := homeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// 使用 ServiceAccount 创建集群配置（InCluster模式）
	if config, err = rest.InClusterConfig(); err != nil {
		// 使用 KubeConfig 文件创建集群配置
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig); err != nil {
			panic(err.Error())
		}
	}
	// 创建 clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientSet, config
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
