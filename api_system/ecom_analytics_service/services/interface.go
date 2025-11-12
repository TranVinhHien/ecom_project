package services

import (
	iservices "github.com/TranVinhHien/ecom_analytics_service/services/interface"
)

type ServiceUseCase interface {
	iservices.ServiceSITEUseCase
	iservices.FeedbackUseCase
}
