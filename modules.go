package main

import (
	_ "github.com/PlanitarInc/registrator/consul"
	_ "github.com/PlanitarInc/registrator/consulkv"
	_ "github.com/PlanitarInc/registrator/etcd"
	_ "github.com/PlanitarInc/registrator/skydns2"
	_ "github.com/PlanitarInc/registrator/zookeeper"
)
