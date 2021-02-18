package main

import (
	"context"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
)

type Calls struct{}

func (c Calls) OnStreamOpen(ctx context.Context, id int64, str string) error {
	return nil
}
func (c Calls) OnStreamClosed(id int64) {}
func (c Calls) OnStreamRequest(id int64, req *discovery.DiscoveryRequest) error {
	return nil
}
func (c Calls) OnStreamResponse(id int64, req *discovery.DiscoveryRequest, res *discovery.DiscoveryResponse) {

}

func (c Calls) OnFetchRequest(ctx context.Context, req *discovery.DiscoveryRequest) error {
	return nil
}

func (c Calls) OnFetchResponse(req *discovery.DiscoveryRequest, res *discovery.DiscoveryResponse) {
}
