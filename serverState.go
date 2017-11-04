package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ServerState() {
	now := time.Now()
	end := now.Unix()
	start := end - 60
	var metrics Resp_PrometheusMetrics
	var metric_report = MetricReport{
		Title: "server_state",
	}

	url := Config.General.Prometheus + "/api/v1/query_range?query=sum(up%7Bkubernetes_io_role%3D~%22.*%22%2C+job%3D%22kubernetes-nodes%22%7D)+by+(instance)&start=" + strconv.FormatInt(start, 10) + "&end=" + strconv.FormatInt(end, 10) + "&step=61"
	fmt.Println(url)

	metrics = Get_PrometheusMetrics(url)
	fmt.Println(metrics)

	// var totalValue float64
	for _, value := range metrics.Data.Result {
		for _, v := range value.Values {
			sliceValue := v.([]interface{})
			// fmt.Println(sliceValue)
			metric, err := strconv.ParseFloat(sliceValue[1].(string), 64)

			if err == nil {
				if int(metric) < Config.Metrics["server_state"].Threshold {
					metric_report.Issue = append(metric_report.Issue, strings.Replace(value.Metric.Instance, ":9100", "", -1)+": Down")
				}

			}
		}

	}
	UpdateDataReport(metric_report)
}
