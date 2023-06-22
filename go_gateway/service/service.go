package service

import (
	eoscContext "github.com/baker-yuan/go-gateway/context"
	"github.com/baker-yuan/go-gateway/eosc"
)

const (
	ServiceSkill = "github.com/baker-yuan/go-gateway/service.service.IService"
)

type IService interface {
	eosc.IWorker
	eoscContext.EoApp
	eoscContext.BalanceHandler
	eoscContext.UpstreamHostHandler
}

// CheckSkill 检查目标技能是否符合
func CheckSkill(skill string) bool {
	return skill == ServiceSkill
}
