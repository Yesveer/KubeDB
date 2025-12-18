package models

import "time"

type MetalLBPool struct {
	Domain  string `bson:"domain" json:"domain"`
	Project string `bson:"project" json:"project"`
	Cluster string `bson:"cluster" json:"cluster"`

	PoolName          string   `bson:"pool_name" json:"pool_name"`
	Addresses         []string `bson:"addresses" json:"addresses"`
	AdvertisementName string   `bson:"advertisement_name" json:"advertisement_name"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
