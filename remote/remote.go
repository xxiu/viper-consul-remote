package remote

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/consul/api"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

const (
	defaultWaitTime = 10 * time.Minute
)

type remoteConfigProvider struct{}

type Consul struct {
	Address    string
	Token      string
	Datacenter string
	Client     *consulapi.Client
	kv         *consulapi.KV
}

func (rc remoteConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {

	consul, err := getConfigManager(rp)
	if err != nil {
		return nil, err
	}
	pair, _, err := consul.kv.Get(rp.Path(), nil)

	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(pair.Value), nil

}

func (rc remoteConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	consul, err := getConfigManager(rp)
	if err != nil {
		return nil, err
	}
	pair, _, err := consul.kv.Get(rp.Path(), nil)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(pair.Value), nil
}

func (rc remoteConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {

	consul, err := getConfigManager(rp)
	if err != nil {
		return nil, nil
	}

	quitwc := make(chan bool)
	viperResponsCh := make(chan *viper.RemoteResponse)

	go func(rp viper.RemoteProvider, c *Consul, vr chan<- *viper.RemoteResponse) {
		withindex := uint64(0)
		for {
			opts := api.QueryOptions{
				WaitIndex: withindex,
				WaitTime:  defaultWaitTime,
			}
			keypair, meta, err := c.Client.KV().Get(rp.Path(), &opts)
			if keypair == nil && err == nil {
				err = fmt.Errorf("Key ( %s ) was not found.", rp.Path())
			}

			vr <- &viper.RemoteResponse{
				Error: err,
				Value: keypair.Value,
			}
			withindex = meta.LastIndex
		}
	}(rp, consul, viperResponsCh)

	return viperResponsCh, quitwc

}

func getConfigManager(rp viper.RemoteProvider) (*Consul, error) {

	if rp.Provider() != "consul" {
		return nil, errors.New("remote is consul!")
	}

	consul := &Consul{
		Address:    rp.Endpoint(),
		Token:      rp.SecretKeyring(),
		Datacenter: "dc1",
	}

	config := consulapi.DefaultConfig()
	config.Address = consul.Address
	config.Token = consul.Token
	config.Datacenter = consul.Datacenter

	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	consul.Client = client
	consul.kv = client.KV()

	return consul, nil
}

func init() {
	viper.RemoteConfig = &remoteConfigProvider{}
}
