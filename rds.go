package main

import (
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

func BuildRds() []types.Resource {
	// 通过 RDS 下发 vhds ，而且只能通过 RDS 来下发
	r := &route.RouteConfiguration{
		Name: ResourceName,
		Vhds: &route.Vhds{
			ConfigSource: &core.ConfigSource{
				ResourceApiVersion: core.ApiVersion_V3,
				ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
					ApiConfigSource: &core.ApiConfigSource{
						TransportApiVersion: core.ApiVersion_V3,
						ApiType:             core.ApiConfigSource_DELTA_GRPC,
						GrpcServices: []*core.GrpcService{
							{
								TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
									EnvoyGrpc: &core.GrpcService_EnvoyGrpc{
										ClusterName: "xds_cluster",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return []types.Resource{r}
}
