package service

import "github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/vm"

type QueryInput struct {
	VMURL string

	Expr  string
	Start string
	End   string
	Step  string
}

type QueryService struct {
	client *vm.Client
}

func NewQueryService(vmURL string) *QueryService {
	return &QueryService{
		client: vm.NewClient(vmURL),
	}
}
