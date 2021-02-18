package main

import (
	"fmt"
	config "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	"github.com/golang/protobuf/ptypes"

	"github.com/golang/protobuf/ptypes/any"
	"strconv"
	"time"
)

var (
	version = 1
)

type vhds struct {}

func NewVHDSServer() *vhds {
	return &vhds{}
}

func buildResponse(targetPrefix string) *discovery.DeltaDiscoveryResponse {
	resources := []*discovery.Resource{
		{
			Name:     ResourceName,
			Version:  strconv.Itoa(version),
			Resource: buildNewVHDSResources(targetPrefix),
		},
	}
	version++
	return &discovery.DeltaDiscoveryResponse{
		TypeUrl:           VHDSTypeUrl,
		Resources:         resources,
		RemovedResources:  []string{},
	}
}

func buildNewVHDSResources(targetPrefix string) *any.Any {
	var clusterName = ClusterName
	var virtualHostName = "service"
	vhost := &config.VirtualHost{
		Name:    virtualHostName,
		Domains: []string{"*"},
		Routes: []*config.Route{{
			Match: &config.RouteMatch{
				PathSpecifier: &config.RouteMatch_Prefix{
					Prefix: targetPrefix,
				},
			},

			Action: &config.Route_Route{
				Route: &config.RouteAction{
					ClusterSpecifier: &config.RouteAction_Cluster{
						Cluster: clusterName,
					},
				},
			},
		}}}
	vh, _ := ptypes.MarshalAny(vhost)
	return vh
}

func (v *vhds) DeltaVirtualHosts(dvhs route.VirtualHostDiscoveryService_DeltaVirtualHostsServer) error {
	fmt.Println("已经连接到 vhds ")
	go func() {
		for {
			req, err := dvhs.Recv()
			if err != nil {
				break
			}
			// TODO 补充上根据 envoy 的上报来下发资源
			//  InitialResourceVersions 就是当前 envoy 的资源名和版本
			// 	比如某一时刻断开了 envoy 和控制面的连接，后重连上去
			//  会返回一个 map[rds_routes:version]
			fmt.Println("InitialResourceVersions: ", req.GetInitialResourceVersions())
			// GetResourceNamesSubscribe 是当前 envoy 持有的资源名
			// 类似 [rds_routes]
			fmt.Println("ResourceNamesSubscribe: ", req.GetResourceNamesSubscribe())
			// GetResourceNamesUnsubscribe 暂时还没有搞清楚，应该是要删除的？
			fmt.Println("ResourceNamesUnsubscribe: ", req.GetResourceNamesUnsubscribe())
		}
	}()

	for range time.Tick(10 * time.Second) {
		// 通过 vhds 下发来变更 targetPrefix 的 VirtualHost
		// 现在可以通过 127.0.0.1:80/route(version%3 + 1) 来访问了
		prefix := "/route" + strconv.Itoa(version%3 + 1)
		fmt.Println("通过 vhds 来变更 VirtualHost，现在可以访问 127.0.0.1" + prefix)
		if err := dvhs.Send(buildResponse(prefix)); err != nil {
			return err
		}
	}

	return nil
}
