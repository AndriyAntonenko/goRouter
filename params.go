package router

import (
	"strconv"
)

type RouterParams struct {
	Params map[string]string
}

func NewRouterParams() *RouterParams {
	return &RouterParams{
		Params: make(map[string]string),
	}
}

func (rp *RouterParams) addParam(name string, value string) {
	rp.Params[name] = value
}

func (rp *RouterParams) GetString(name string) string {
	value, ok := rp.Params[name]
	if !ok {
		return ""
	}

	return value
}

func (rp *RouterParams) ParseInt(name string) int {
	value := rp.GetString(name)
	if value == "" {
		return 0
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return intValue
}
