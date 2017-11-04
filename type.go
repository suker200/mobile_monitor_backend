package main

import (
	"sync"
)

type DataReport struct {
	sync.Mutex
	Infos []MetricReport `json:"infos"`
}

type MetricReport struct {
	Title string   `json:"title"`
	Alias string   `json:"alias"`
	Issue []string `json:"issue"`
}

type ConfigInfo struct {
	General struct {
		Prometheus string `json:"prometheus"`
	} `json:"general"`
	Metrics map[string]MetricInfo `json:"metrics"`
}

type MetricInfo struct {
	Name      string `json:"name"`
	Interval  int    `json:"interval"`
	Alias     string `json:"alias"`
	Project   string `json:"project"`
	Label     string `json:"label"`
	Threshold int    `json:"threshold"`
}

type Resp_PrometheusMetrics struct {
	Status string `json:"status"`
	Data   struct {
		Result []struct {
			Metric struct {
				Instance  string `json:"instance"`
				Node      string `json:"node"`
				Namespace string `json:"namespace"`
				Pod       string `json:"pod"`
				Upstream  string `json:"upstream"`
				Server    string `json:"server"`
			} `json:"metric"`
			Value  []interface{} `json:"value"`
			Values []interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}
