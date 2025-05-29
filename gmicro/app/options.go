package app

import (
	"minifast/gmicro/registry"
	"net/url"
	"os"
	"time"
)

type Option func(o *options)

type options struct {
	id        string
	endpoints []url.URL
	name      string

	sigs []os.Signal

	// 允许用户传入自己的实现
	registrar        registry.Registrar
	registrarTimeout time.Duration

	// stop 超时时间
	stopTimeout time.Duration

	restServer *restserver.Server
	rpcServer  *rpcserver.Server
}

func WithRegistrar(registrar registry.Registrar) Option {
	return func(o *options) {
		o.registrar = registrar
	}
}

func WithEndpoints(endpoints ...url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func WithRpcServer(rpcServer *rpcserver.Server) Option {
	return func(o *options) {
		o.rpcServer = rpcServer
	}
}

func WithRestServer(restServer *restserver.Server) Option {
	return func(o *options) {
		o.restServer = restServer
	}
}

func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithSigs(sigs ...os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}
