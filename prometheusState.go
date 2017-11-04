package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func PrometheusThrotle() {
	now := time.Now()
	end := now.Unix()
	start := end - 60
	var metrics Resp_PrometheusMetrics
	var metric_report = MetricReport{
		Title: "prometheus_throtle",
	}

	url := Config.General.Prometheus + "/api/v1/query_range?query=max_over_time(prometheus_local_storage_memory_series%7Bjob%3D%22prometheus%22%7D%5B6h%5D)+%2F+1000000+%2F+(prometheus_local_storage_target_heap_size_bytes%7B%7D+%2F+1024+%2F+1024+%2F+1024)&start=" + strconv.FormatInt(start, 10) + "&end=" + strconv.FormatInt(end, 10) + "&step=61"
	fmt.Println(url)

	metrics = Get_PrometheusMetrics(url)
	fmt.Println(metrics)

	// var totalValue float64
	for _, value := range metrics.Data.Result {
		for _, v := range value.Values {
			sliceValue := v.([]interface{})
			metric, err := strconv.ParseFloat(sliceValue[1].(string), 64)

			if err == nil {
				if int(metric) > Config.Metrics["prometheus_throtle"].Threshold {
					metric_report.Issue = append(metric_report.Issue, strings.Replace(value.Metric.Instance, ":9100", "", -1)+" Throtle: "+fmt.Sprintf("%.3f", metric)+"%")
				}

			}
		}

	}
	UpdateDataReport(metric_report)
}
