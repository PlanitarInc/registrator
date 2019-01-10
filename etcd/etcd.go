package etcd

import (
	"log"
	"net"
	"net/url"
	"strconv"

	"github.com/gliderlabs/registrator/bridge"
)

func init() {
	bridge.Register(new(Factory), "etcd")
}

type Factory struct{}

func (f *Factory) New(uri *url.URL) bridge.RegistryAdapter {
	client, err := NewEtcdClient(uri.Host)
	if err != nil {
		log.Fatal("etcd: can't allocate client:", err)
	}

	if client.client != nil {
		log.Println("etcd: using v0 client")
	}
	return &EtcdAdapter{client: client, path: uri.Path}
}

type EtcdAdapter struct {
	client *EtcdClient

	path string
}

func (r *EtcdAdapter) Ping() error {
	r.syncEtcdCluster()

	return r.client.Ping()
}

func (r *EtcdAdapter) syncEtcdCluster() {
	if !r.client.SyncEtcdCluster() {
		log.Println("etcd: sync cluster was unsuccessful")
	}
}

func (r *EtcdAdapter) Register(service *bridge.Service) error {
	r.syncEtcdCluster()

	path := r.path + "/" + service.Name + "/" + service.ID
	port := strconv.Itoa(service.Port)
	addr := net.JoinHostPort(service.IP, port)

	err := r.client.Set(path, addr, uint64(service.TTL))
	if err != nil {
		log.Println("etcd: failed to register service:", err)
	}
	return err
}

func (r *EtcdAdapter) Deregister(service *bridge.Service) error {
	r.syncEtcdCluster()

	path := r.path + "/" + service.Name + "/" + service.ID

	err := r.client.Delete(path, false)
	if err != nil {
		log.Println("etcd: failed to deregister service:", err)
	}
	return err
}

func (r *EtcdAdapter) Refresh(service *bridge.Service) error {
	return r.Register(service)
}

func (r *EtcdAdapter) Services() ([]*bridge.Service, error) {
	return []*bridge.Service{}, nil
}
