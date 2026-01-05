package server

import (
	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
	"review-service/internal/conf"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewRegister, NewGRPCServer, NewHTTPServer)

func NewRegister(conf *conf.Registry) registry.Registrar {
	//new Consul client
	c := api.DefaultConfig()
	c.Address = conf.Consul.Address //自己的配置
	c.Scheme = conf.Consul.PassDeregister
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	// new reg with consul client
	reg := consul.New(client, consul.WithHealthCheck(true))
	return reg
}
