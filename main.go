package main

import (
	"context"
	cluster "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	callbacks := Calls{}
	snapshotCache := cache.NewSnapshotCache(false, cache.IDHash{}, nil)

	server := xds.NewServer(context.Background(), snapshotCache, callbacks)

	if err := snapshotCache.SetSnapshot("node1", cache.NewSnapshot(
		"1.0",
		nil,
		BuildCds(),
		BuildRds(),
		nil,
		nil,
		nil)); err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	listen, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}

	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	endpoint.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	cluster.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	route.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	listener.RegisterListenerDiscoveryServiceServer(grpcServer, server)

	v := NewVHDSServer()
	route.RegisterVirtualHostDiscoveryServiceServer(grpcServer, v)

	panic(grpcServer.Serve(listen))
}
