package api

import (
	"github.com/enbility/eebus-go/api"
)

type Solution struct {
	Service api.ServiceInterface
}

func NewSolution(service api.ServiceInterface) *Solution {
	return &Solution{
		Service: service,
	}
}
