// proxy serve
package main

import (
	"log"
	"net/rpc"
)

type ProxyStore struct {
	urls   *URLStore // urls cache in slave server
	client *rpc.Client
}

func NewProxyStore(addr string) *ProxyStore {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Println("ProxyStore:", err)
	}
	// do not create file in slave server
	return &ProxyStore{urls: NewURLStore(""), client: client}
}

// forward calls to the master server
func (s *ProxyStore) Get(key, url *string) error {
	if err := s.urls.Get(key, url); err == nil { // url found in local map
		return nil
	}
	// url not found in local map, make rpc-call:
	if err := s.client.Call("Store.Get", key, url); err != nil {
		return err
	}
	s.urls.Set(key, url) // update cache
	return nil
}

func (s *ProxyStore) Put(url, key *string) error {
	if err := s.client.Call("Store.Put", url, key); err != nil {
		return err
	}
	s.urls.Set(key, url) // update cache
	return nil
}
