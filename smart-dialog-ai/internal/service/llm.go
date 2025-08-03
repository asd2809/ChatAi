package service

import (
	"context"

	"github.com/sirupsen/logrus"
)

type SiliconFlowHandler struct {
	ctx     context.Context
	logger  *logrus.Logger
}