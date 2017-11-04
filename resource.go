package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Mem_Usage() {
	now := time.Now()
	end := now.Unix()
	start := end - 300
	var metrics Resp_PrometheusMetrics
	var metric_report = MetricReport{
		Title: "mem_usage",
	}

	url := Config.General.Prometheus + "/api/v1/query_range?query=sum+by+(instance)+(100+-+node_memory_MemAvailable%7B%7D+%2F+node_memory_MemTotal%7B%7D+*+100+)+&start=" + strconv.FormatInt(start, 10) + "&end=" + strconv.FormatInt(end, 10) + "&step=301"
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
				if int(metric) > Config.Metrics["mem_usage"].Threshold {
					metric_report.Issue = append(metric_report.Issue, strings.Replace(value.Metric.Instance, ":9100", "", -1)+":\t\t "+fmt.Sprintf("%.2f", metric)+"%")
				}

			}
		}

	}
	UpdateDataReport(metric_report)
}

func Disk_Usage() {
	now := time.Now()
	end := now.Unix()
	start := end - 900
	var metrics Resp_PrometheusMetrics
	var metric_report = MetricReport{
		Title: "disk_usage",
	}

	url := Config.General.Prometheus + "/api/v1/query_range?query=sum((node_filesystem_size%7Bdevice%3D~%22.*xvda.*%22%2C+mountpoint%3D%22%2Fetc%2Fresolv.conf%22%7D+-+node_filesystem_avail%7Bdevice%3D~%22.*xvda.*%22%2C+mountpoint%3D%22%2Fetc%2Fresolv.conf%22%7D)+%2F+node_filesystem_size%7Bdevice%3D~%22.*xvda.*%22%2C+mountpoint%3D%22%2Fetc%2Fresolv.conf%22%7D+*+100+)+by+(instance)+&start=" + strconv.FormatInt(start, 10) + "&end=" + strconv.FormatInt(end, 10) + "&step=901"
	fmt.Println(url)

	metrics = Get_PrometheusMetrics(url)
	fmt.Println(metrics)

	// var totalValue float64
	for _, value := range metrics.Data.Result {
		for _, v := range value.Values {
			sliceValue := v.([]interface{})
			metric, err := strconv.ParseFloat(sliceValue[1].(string), 64)

			if err == nil {
				if int(metric) > Config.Metrics["disk_usage"].Threshold {
					metric_report.Issue = append(metric_report.Issue, strings.Replace(value.Metric.Instance, ":9100", "", -1)+":\t\t "+fmt.Sprintf("%.2f", metric)+"%")
				}

			}
		}

	}
	UpdateDataReport(metric_report)
}

func Load_Usage() {
	var metric_report = MetricReport{
		Title: "load_usage",
	}

	var load_metrics Resp_PrometheusMetrics
	var cpu_core_metrics Resp_PrometheusMetrics
	load_metrics = Get_PrometheusMetrics(Config.General.Prometheus + "/api/v1/query?query=node_load5{}")
	cpu_core_metrics = Get_PrometheusMetrics(Config.General.Prometheus + "/api/v1/query?query=machine_cpu_cores{}")

	for num, value := range load_metrics.Data.Result {
		load, err := strconv.ParseFloat(value.Value[1].(string), 64)
		core, err := strconv.ParseFloat(cpu_core_metrics.Data.Result[num].Value[1].(string), 64)

		if err == nil && int(load-core) > Config.Metrics["load_usage"].Threshold {
			metric_report.Issue = append(metric_report.Issue, strings.Replace(value.Metric.Instance, ":9100", "", -1)+":\t\t "+fmt.Sprintf("%.2f", load-core))
		}
	}

	UpdateDataReport(metric_report)
}
