package etcd

import (
	"context"
	"fmt"
	"github.com/LERSONG/beetle/registry"
	"github.com/LERSONG/beetle/util/addr"
	"github.com/coreos/etcd/clientv3"
	"os"
	"testing"
)

var conf clientv3.Config

func TestMain(m *testing.M) {
	conf.Endpoints = []string{
		"10.200.100.200:42379",
	}
	os.Exit(m.Run())
}

func TestRegister(t *testing.T) {
	registry, err := NewRegistry(
		registry.Addrs([]string{"10.200.100.200:42379"}),
		registry.ServiceName("lc-test"),
		registry.ServiceAddr(":8080"),
	)
	if err != nil {
		t.Failed()
		return
	}
	err = registry.Register()
	if err != nil {
		t.Fail()
		return
	}

	addr, err := addr.Extract("")
	if err != nil {
		t.Fail()
		return
	}
	expectAddr := fmt.Sprintf("%s:%s", addr, "8080")
	client, err := clientv3.New(conf)
	if err != nil {
		t.Fail()
		return
	}

	path := fmt.Sprintf("%s_", "lc-test")
	resp, err := client.Get(context.Background(), path, clientv3.WithPrefix())
	if err != nil {
		t.Fail()
		return
	}
	exist := false
	for _, kv := range resp.Kvs {
		if string(kv.Value) == expectAddr {
			exist = true
		}
	}
	if !exist {
		t.Fatal("test registry fail")
	}

}
