package config

import (
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

func NewMidtransSnapClient() snap.Client {
	var s snap.Client

	s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	return s
}

func NewMidtransCoreAPIClient() coreapi.Client {
	var c coreapi.Client

	c.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	return c
}
