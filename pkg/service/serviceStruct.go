package service

import (
	"github.com/ratheeshkumar25/opt_cut_material_service/config"
	inter "github.com/ratheeshkumar25/opt_cut_material_service/pkg/repo/interfaces"
	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/service/interfaces"
)

type MaterialService struct {
	Repo      inter.MaterialRepoInter
	StripePay *config.StripeClient
	redis     *config.RedisService
	//RazorPay  *config.RazorpayClient

}

func NewMaterialService(repo inter.MaterialRepoInter, strPay *config.StripeClient, redis *config.RedisService) interfaces.MaterialServiceInte {
	return &MaterialService{
		Repo:      repo,
		StripePay: strPay,
		redis:     redis,
	}
}

//razorpay *config.RazorpayClient
