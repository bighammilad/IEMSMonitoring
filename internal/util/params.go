package util

import (
	"net/url"
)

type UrlParams map[string][]string

func NewUrlParams(vals url.Values) UrlParams {
	params := UrlParams{}
	for k, v := range vals {
		params[k] = v
	}
	//delete(params, "")
	//add vals to params
	return params
}

func (p UrlParams) Get(key string) string {
	// 1. key exist or not and greater or equal than 1
	if len(p[key]) == 0 {
		delete(p, key)
		return ""
	}
	value := p[key][0]
	delete(p, key)
	return value
}

func (p *UrlParams) IsEmpty() bool {
	return len(*p) == 0
}
