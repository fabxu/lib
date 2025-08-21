package util

import "strings"

type ConfigMap map[string]string

func (m ConfigMap) TrimSpace() ConfigMap {
	for k, v := range m {
		m[k] = strings.TrimSpace(v)
	}

	return m
}

func (m ConfigMap) AnyEmpty() (bool, string) {
	if len(m) == 0 {
		return false, ""
	}

	for k, v := range m {
		if strings.TrimSpace(v) == "" {
			return true, k
		}
	}

	return false, ""
}
