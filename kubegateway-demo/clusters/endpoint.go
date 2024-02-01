package clusters

import (
	"context"
	"k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"sync"
)

type EndpointInfo struct {
	ctx    context.Context
	cancel context.CancelFunc

	Cluster  string
	Endpoint string

	proxyConfig        *rest.Config
	proxyUpgradeConfig *rest.Config
	// http2 proxy round tripper
	ProxyTransport http.RoundTripper
	// http1 proxy round tripper for websocket
	PorxyUpgradeTransport proxy.UpgradeRequestRoundTripper

	clientset kubernetes.Interface

	sync.Mutex
}

func UnwrapUpgradeRequestRoundTripper(rt http.RoundTripper) (proxy.UpgradeRequestRoundTripper, bool) {
	urrt, ok := rt.(proxy.UpgradeRequestRoundTripper)
	if ok {
		return urrt, ok
	}

	var rtw net.RoundTripperWrapper
	var isWrapper bool
	rtw, isWrapper = rt.(net.RoundTripperWrapper)
	for isWrapper {
		rt = rtw.WrappedRoundTripper()
		urrt, found := rt.(proxy.UpgradeRequestRoundTripper)
		if found {
			return urrt, true
		}
		rtw, isWrapper = rt.(net.RoundTripperWrapper)
	}

	return nil, false
}
