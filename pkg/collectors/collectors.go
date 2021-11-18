package collectors

import (
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricLabel struct {
	Name   string
	Labels []string
}

type SystemInfo struct {
	Name    string
	Vendor  string
	Model   string
	Version string
}
type AppInfo struct {
	Name    string
}

type PerfCollector struct {
	sysInfoDescriptors map[string]*prometheus.Desc
	sysPerfDescriptors map[string]*prometheus.Desc
	appDescriptors map[string]*prometheus.Desc
	up prometheus.Gauge
	sequenceNumber uint64
}
const (
	// Storage Metrics
	StorageReadIOPS     = "fusion_storage_perf_rd_iops"
	StorageWriteIOPS    = "fusion_storage_perf_wr_iops"
	StorageReadBytes    = "fusion_storage_perf_rd_bytes"
	StorageWriteBytes   = "fusion_storage_perf_wr_bytes"
	StorageLatency      = "fusion_storage_perf_latency_seconds"
	StorageReadLatency  = "fusion_storage_perf_rd_latency_seconds"
	StorageWriteLatency = "fusion_storage_perf_wr_latency_seconds"

	StorageAvailableBytes = "fusion_storage_capacity_available_bytes"
	StorageUsedBytes = "fusion_storage_capacity_used_bytes"

	// System info and state
	SystemMetadata = "fusion_subsystem_metadata"
	SystemHealth   = "fusion_subsystem_health_state"

	// Application metrics
	AppCapacityAvailableBytes = "fusion_application_capacity_available_bytes"
	AppCapacityUsedBytes = "fusion_application_capacity_used_bytes"

	AppBackupDuration = "fusion_application_backup_duration"
	AppBackupUsedBytes = "fusion_application_backup_used_bytes"
	AppBackupAvailableBytes = "fusion_application_backup_available_bytes"
	AppBackupJobRunningCount = "fusion_application_backup_job_running_count"
	AppBackupJobSuccessCount = "fusion_application_backup_job_success_count"
	AppBackupJobFailedCount = "fusion_application_backup_job_failed_count"
)

var (
	// Metadata label
	subsystemMetadataLabel = []string{"subsystem_name", "vendor", "model", "version"}

	// Other label
	subsystemCommonLabel = []string{"subsystem_name"}

	// Application label
	appLabels = []string{"subsystem_name", "app",}

	systemMetricsMap = map[string]MetricLabel{
		SystemMetadata: {"System information", subsystemMetadataLabel},
		SystemHealth:   {"System health", subsystemCommonLabel},
	}

	perfMetricsMap = map[string]MetricLabel{
		StorageReadIOPS:     {"overall performance - read IOPS", subsystemCommonLabel},
		StorageWriteIOPS:    {"overall performance - write IOPS", subsystemCommonLabel},
		StorageReadBytes:    {"overall performance - read throughput bytes/s", subsystemCommonLabel},
		StorageWriteBytes:   {"overall performance - write throughput bytes/s", subsystemCommonLabel},
		StorageLatency:      {"overall performance - average latency seconds", subsystemCommonLabel},
		StorageReadLatency:  {"overall performance - read latency seconds", subsystemCommonLabel},
		StorageWriteLatency: {"overall performance - write latency seconds", subsystemCommonLabel},
		StorageAvailableBytes: {"overall capacity - available bytes", subsystemCommonLabel},
		StorageUsedBytes: {"overall capacity - used bytes", subsystemCommonLabel},
	}

	appMetricsMap = map[string]MetricLabel{
		AppCapacityAvailableBytes: {"application capacity - available bytes", appLabels},
		AppCapacityUsedBytes: {"application capacity - used bytes", appLabels},
		AppBackupAvailableBytes: {"application backup capacity - available bytes", appLabels},
		AppBackupUsedBytes: {"application backup capacity - used bytes", appLabels},
		AppBackupDuration: {"application capacity - available bytes", appLabels},
		AppBackupJobFailedCount: {"application capacity - available bytes", appLabels},
		AppBackupJobRunningCount: {"application capacity - available bytes", appLabels},
		AppBackupJobSuccessCount: {"application capacity - available bytes", appLabels},
	}
)

func NewPerfCollector() (*PerfCollector, error) {
	f := &PerfCollector{
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "up",
			Help: "Was the last scrape successful.",
		}),
	}
	f.initSubsystemDescs()

	return f, nil
}

func (f *PerfCollector) Describe(ch chan<- *prometheus.Desc) {

	for _, v := range f.sysInfoDescriptors {
		ch <- v
	}

	for _, v := range f.sysPerfDescriptors {
		ch <- v
	}

	ch <- f.up.Desc()

}

func (f *PerfCollector) Collect(ch chan<- prometheus.Metric) {
	f.sequenceNumber++
	f.collectSystemMetrics(ch)

	ch <- f.up

}

func (f *PerfCollector) initSubsystemDescs() {
	f.sysInfoDescriptors = make(map[string]*prometheus.Desc)
	f.sysPerfDescriptors = make(map[string]*prometheus.Desc)
	f.appDescriptors = make(map[string]*prometheus.Desc)

	for metricName, metricLabel := range systemMetricsMap {
		f.sysInfoDescriptors[metricName] = prometheus.NewDesc(
			metricName,
			metricLabel.Name, metricLabel.Labels, nil,
		)
	}

	for metricName, metricLabel := range perfMetricsMap {
		f.sysPerfDescriptors[metricName] = prometheus.NewDesc(
			metricName,
			metricLabel.Name, metricLabel.Labels, nil,
		)
	}

	for metricName, metricLabel := range appMetricsMap {
		f.appDescriptors[metricName] = prometheus.NewDesc(
			metricName,
			metricLabel.Name, metricLabel.Labels, nil,
		)
	}
}

func (f *PerfCollector) collectSystemMetrics(ch chan<- prometheus.Metric) bool {
	// system metadata
	systemInfo := SystemInfo{
		Name: "ibm-fusion",
		Vendor: "IBM",
		Model: "SDS",
		Version: "0.1.0",
	};

	newSystemMetrics(ch, f.sysInfoDescriptors[SystemMetadata], 0, &systemInfo)
	// Determine the health 0 = OK, 1 = warning, 2 = error
	status := 0.0
	t := time.Now().Minute() 
	if t%3 == 0 {
		status = 1
	}
	newPerfMetrics(ch, f.sysInfoDescriptors[SystemHealth], status, &systemInfo)

	// Parse Perf Results
	metricValue:= 30*1024 + float64(rand.Intn(2*1024)*10) // read iops 30k - 50k
	newPerfMetrics(ch, f.sysPerfDescriptors[StorageReadIOPS], metricValue , &systemInfo)
	metricValue = 10*1024 + float64(rand.Intn(2*1024)*10) // write iops 10k - 30k
	newPerfMetrics(ch, f.sysPerfDescriptors[StorageWriteIOPS], metricValue , &systemInfo)
	metricValue = 70*1024*1024 + float64(rand.Intn(2*1024)*10*1024) // read bytes 70M - 90M
	newPerfMetrics(ch, f.sysPerfDescriptors[StorageReadBytes], metricValue , &systemInfo)
	metricValue = 20*1024*1024 + float64(rand.Intn(2*1024)*10*1024) // write bytes 20M - 40M
	newPerfMetrics(ch, f.sysPerfDescriptors[StorageWriteBytes], metricValue , &systemInfo)
	metricValue = 0.0001 + float64(rand.Intn(10))*0.0001 // read latency 0.0001 - 0.001
	newPerfMetrics(ch, f.sysPerfDescriptors[StorageReadLatency], metricValue , &systemInfo)
	metricValue = 0.0001 + float64(rand.Intn(10))*0.0001 // write latency 0.0001 - 0.001
	newPerfMetrics(ch, f.sysPerfDescriptors[StorageWriteLatency], metricValue , &systemInfo)
	const totalBytes = 50*1024*1024*1024*1024 // 50T
	metricValue = 30*1024*1024*1024*1024 + float64(rand.Intn(2*1024)*10*1024*1024*1024) // used bytes 30T - 50T
	newPerfMetrics(ch, f.sysPerfDescriptors[StorageUsedBytes], metricValue , &systemInfo)
	newPerfMetrics(ch, f.sysPerfDescriptors[StorageAvailableBytes], totalBytes - metricValue , &systemInfo)
	/*
	for _, m := range f.sysPerfDescriptors {
		unixT := time.Now().Unix()
		metricValue:= float64(unixT)* float64(rand.Intn(100))
		newPerfMetrics(ch, m, metricValue , &systemInfo)
	}
	*/

	// Parse Application Results
	appSamples := []string{"book-store", "wordpress", "devOps"}
	const appProvisionedBytes = 100*1024*1024*1024*1024 // 100T
	const appBackupTotalBytes = 100*1024*1024*1024*1024 // 100T
	for _, appName:= range appSamples {
		appInfo := AppInfo{
			Name: appName,
		}
		metricValue = 10*1024*1024*1024*1024 + float64(rand.Intn(2*1024)*10*1024*1024*1024) // app used bytes 10T - 30T
		newAppMetrics(ch, f.appDescriptors[AppCapacityUsedBytes], metricValue , &systemInfo, &appInfo)
		newAppMetrics(ch, f.appDescriptors[AppCapacityAvailableBytes], appProvisionedBytes - metricValue , &systemInfo, &appInfo)
		metricValue = 60*1024*1024*1024*1024 + float64(rand.Intn(2*1024)*10*1024*1024*1024) // backup used bytes 60T - 80T
		newAppMetrics(ch, f.appDescriptors[AppBackupUsedBytes], metricValue , &systemInfo, &appInfo)
		newAppMetrics(ch, f.appDescriptors[AppBackupAvailableBytes], appBackupTotalBytes - metricValue , &systemInfo, &appInfo)
		metricValue = 59 + float64(rand.Intn(60)*60) // backup duration 1 minute - 60 minutes
		newAppMetrics(ch, f.appDescriptors[AppBackupDuration], metricValue , &systemInfo, &appInfo)
		metricValue = float64(rand.Intn(60)) // backup failed count 1 - 60 
		newAppMetrics(ch, f.appDescriptors[AppBackupJobFailedCount], metricValue , &systemInfo, &appInfo)
		metricValue = float64(rand.Intn(60)) // backup running count 1 - 60
		newAppMetrics(ch, f.appDescriptors[AppBackupJobRunningCount], metricValue , &systemInfo, &appInfo)
		metricValue = float64(rand.Intn(60)) // backup success count 1 - 60
		newAppMetrics(ch, f.appDescriptors[AppBackupJobSuccessCount], metricValue , &systemInfo, &appInfo)
		/*
		for _, m := range f.appDescriptors {
			unixT := time.Now().Unix()
			metricValue:= float64(unixT)* float64(rand.Intn(100))
			newAppMetrics(ch, m, metricValue , &systemInfo, &appInfo)
		}
		*/
	}
	return true
}

func newSystemMetrics(ch chan<- prometheus.Metric, desc *prometheus.Desc, value float64, info *SystemInfo) {
	ch <- prometheus.MustNewConstMetric(
		desc,
		prometheus.GaugeValue,
		value,
		info.Name,
		info.Vendor,
		info.Model,
		info.Version,
	)
}

func newPerfMetrics(ch chan<- prometheus.Metric, desc *prometheus.Desc, value float64, info *SystemInfo) {
	ch <- prometheus.MustNewConstMetric(
		desc,
		prometheus.GaugeValue,
		value,
		info.Name,
	)
}

func newAppMetrics(ch chan<- prometheus.Metric, desc *prometheus.Desc, value float64, info *SystemInfo, appInfo *AppInfo) {
	ch <- prometheus.MustNewConstMetric(
		desc,
		prometheus.GaugeValue,
		value,
		info.Name,
		appInfo.Name,
	)
}
