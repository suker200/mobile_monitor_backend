package main

import (
	"fmt"
	"strconv"
	// "strings"
	// "time"
	"regexp"
)

func PodRestart() {
	var metrics Resp_PrometheusMetrics
	var metric_report = MetricReport{
		Title: "pod_restart",
	}

	url := Config.General.Prometheus + "/api/v1/query?query=sum(max_over_time(kube_pod_container_status_restarts%7B%7D%5B5m%5D)+-+min_over_time(kube_pod_container_status_restarts%7B%7D%5B5m%5D))+by+(namespace%2C+pod)+%3E+0"
	fmt.Println(url)

	metrics = Get_PrometheusMetrics(url)
	fmt.Println(metrics)

	// var totalValue float64
	for _, value := range metrics.Data.Result {
		metric, err := strconv.ParseFloat(value.Value[1].(string), 64)

		if err == nil {
			metric_report.Issue = append(metric_report.Issue, value.Metric.Namespace+"--"+value.Metric.Pod+":\t\t "+fmt.Sprintf("%.0f", metric))
		} else {
			fmt.Println(err.Error())
		}
		// }

	}
	UpdateDataReport(metric_report)
}

func PodPending() {
	var metrics Resp_PrometheusMetrics
	var metric_report = MetricReport{
		Title: "pod_pending",
	}

	url := Config.General.Prometheus + "/api/v1/query?query=kube_pod_status_phase%7Bphase%3D%22Pending%22%7D"
	fmt.Println(url)

	metrics = Get_PrometheusMetrics(url)
	fmt.Println(metrics)

	// var totalValue float64
	for _, value := range metrics.Data.Result {
		metric, err := strconv.ParseFloat(value.Value[1].(string), 64)

		if err == nil {
			metric_report.Issue = append(metric_report.Issue, value.Metric.Namespace+"--"+value.Metric.Pod+":\t\t\t "+fmt.Sprintf("%.0f", metric))
		} else {
			fmt.Println(err.Error())
		}
		// }
	}

	UpdateDataReport(metric_report)
}

func UpstreamResponseTime() {
	// now := time.Now()
	// end := now.Unix()
	// start := end - 60
	re := regexp.MustCompile("[0-9]+.[0-9]+.[0-9]+.[0-9]+")
	var metrics Resp_PrometheusMetrics
	var metric_report = MetricReport{
		Title: "upstream_response_time",
	}

	url := Config.General.Prometheus + "/api/v1/query?query=sum+by+(server%2C+upstream)+(nginx_upstream_response_msecs_avg%7Bupstream!%3D%22devops-prometheus-prometheus-server-80%22%7D+%2F1000)+%3E+" + fmt.Sprint(Config.Metrics["upstream_response_time"].Threshold)
	fmt.Println(url)

	metrics = Get_PrometheusMetrics(url)
	fmt.Println(metrics)

	// var totalValue float64
	for _, value := range metrics.Data.Result {
		metric, err := strconv.ParseFloat(value.Value[1].(string), 64)

		if err == nil {
			metric_report.Issue = append(metric_report.Issue, value.Metric.Upstream+"--"+fmt.Sprint(re.FindString(value.Metric.Server))+":\t\t "+fmt.Sprintf("%.3f", metric)+" sec")
		} else {
			fmt.Println(err.Error())
		}
		// }
	}
	UpdateDataReport(metric_report)
}
