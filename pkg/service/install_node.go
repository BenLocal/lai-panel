package service

import (
	"context"

	"github.com/benlocal/lai-panel/pkg/ctx"
)

type InstallNodeService struct {
	appCtx *ctx.AppCtx
}

func NewInstallNodeService(appCtx *ctx.AppCtx) *InstallNodeService {
	return &InstallNodeService{
		appCtx: appCtx,
	}
}

func (s *InstallNodeService) Name() string {
	return "install-node-service"
}

func (s *InstallNodeService) Start(ctx context.Context) error {
	return nil
}

func (s *InstallNodeService) Shutdown() error {
	return nil
}
