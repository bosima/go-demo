package consul

import (
	consulapi "github.com/hashicorp/consul/api"
)

var consulClient *consulapi.Client

func init() {
	consulClient, _ = consulapi.NewClient(consulapi.DefaultConfig())
}
