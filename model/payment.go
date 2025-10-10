package model

import "github.com/midtrans/midtrans-go/snap"

type CreatePaymentResponse struct {
	SnapURL *snap.Response `json:"snap_url"`
}
