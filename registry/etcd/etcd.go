package etcd

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/LERSONG/beetle/registry"
	"github.com/LERSONG/beetle/util/addr"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sync"

	"net"
	"strings"
	"time"
)

var (
	Prefix             = "%s_%s"
	DefaultServiceAddr = ":9999"
)

type etcdRegistry struct {
	client    *clientv3.Client
	options   registry.Options
	leaseId   clientv3.LeaseID
	node      registry.Node
	closeCh   chan struct{}
	closeOnce sync.Once
}

var _ registry.Registry = &etcdRegistry{}

func NewRegistry(opts ...registry.Option) (registry.Registry, error) {
	e := &etcdRegistry{
		options: registry.Options{},
		closeCh: make(chan struct{}),
	}
	err := configure(e, opts...)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func configure(e *etcdRegistry, opts ...registry.Option) error {
	config := clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}

	for _, o := range opts {
		o(&e.options)
	}

	if e.options.ServiceName == "" {
		panic("service name required")
	}

	if e.options.ServiceAddr == "" {
		e.options.ServiceAddr = DefaultServiceAddr
	}

	if e.options.Timeout == 0 {
		e.options.Timeout = 5 * time.Second
	}

	if e.options.Secure || e.options.TLSConfig != nil {
		tlsConfig := e.options.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		config.TLS = tlsConfig
	}

	if e.options.Context != nil {
		u, ok := e.options.Context.Value(authKey{}).(*authCreds)
		if ok {
			config.Username = u.Username
			config.Password = u.Password
		}
		cfg, ok := e.options.Context.Value(logConfigKey{}).(*zap.Config)
		if ok && cfg != nil {
			config.LogConfig = cfg
		}
	}

	var cAddrs []string

	for _, address := range e.options.Addrs {
		if len(address) == 0 {
			continue
		}
		addr, port, err := net.SplitHostPort(address)
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			port = "2379"
			addr = address
			cAddrs = append(cAddrs, net.JoinHostPort(addr, port))
		} else if err == nil {
			cAddrs = append(cAddrs, net.JoinHostPort(addr, port))
		}
	}

	// if we got addrs then we'll update
	if len(cAddrs) > 0 {
		config.Endpoints = cAddrs
	}

	cli, err := clientv3.New(config)
	if err != nil {
		return err
	}
	e.client = cli
	node, err := e.genNode()
	if err != nil {
		return err
	}
	e.node = node
	return nil
}

func (e *etcdRegistry) genNode() (registry.Node, error) {
	node := registry.Node{
		Id:       uuid.New().String(),
		Metadata: e.options.Metadata,
	}
	host, port, err := net.SplitHostPort(e.options.ServiceAddr)
	if err != nil {
		return node, err
	}
	addr, err := addr.Extract(host)
	if err != nil {
		return node, err
	}
	if strings.Count(addr, ":") > 0 {
		addr = "[" + addr + "]"
	}
	node.Address = fmt.Sprintf("%s:%s", addr, port)
	return node, nil
}

func (e *etcdRegistry) Register(ops ...registry.RegisterOption) error {
	registerOptions := registry.RegisterOptions{}
	for _, op := range ops {
		op(&registerOptions)
	}

	if registerOptions.RegisterTTL == 0 {
		registerOptions.RegisterTTL = time.Second * 30
	}

	if registerOptions.RegisterInterval > registerOptions.RegisterTTL || registerOptions.RegisterInterval < registerOptions.RegisterTTL/3 {
		registerOptions.RegisterInterval = registerOptions.RegisterTTL / 3
	}

	e.doRegister(registerOptions)
	go func() {
		ticker := time.NewTicker(registerOptions.RegisterInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				e.doRegister(registerOptions)
			case <-e.closeCh:
				return
			}
		}
	}()

	return nil
}

func (e *etcdRegistry) doRegister(opts registry.RegisterOptions) error {
	if e.leaseId > 0 {
		_, err := e.client.KeepAliveOnce(context.TODO(), e.leaseId)
		if err == nil {
			return nil
		} else if err != rpctypes.ErrLeaseNotFound {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()
	lgr, err := e.client.Grant(ctx, int64(opts.RegisterTTL.Seconds()))
	if err != nil {
		return err
	}
	if lgr != nil {
		_, err = e.client.Put(ctx, e.getPath(), e.node.Address, clientv3.WithLease(lgr.ID))
	} else {
		_, err = e.client.Put(ctx, e.getPath(), e.node.Address)
	}
	if err != nil {
		return err
	}
	e.leaseId = lgr.ID
	return nil
}

func (e *etcdRegistry) Deregister() error {
	e.closeOnce.Do(func() {
		close(e.closeCh)
	})

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	_, err := e.client.Delete(ctx, e.getPath())
	return err
}

func (e *etcdRegistry) String() string {
	return "etcd"
}

func (e etcdRegistry) getPath() string {
	return fmt.Sprintf(Prefix, e.options.ServiceName, e.node.Id)
}
