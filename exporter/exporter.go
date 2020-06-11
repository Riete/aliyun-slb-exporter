package exporter

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *SlbExporter) NewClient() {
	client, err := cms.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}
	s.client = client
}

func (s *SlbExporter) GetInstance() {
	client, err := slb.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}
	request := slb.CreateDescribeLoadBalancersRequest()
	request.PageSize = requests.NewInteger(100)
	response, err := client.DescribeLoadBalancers(request)
	if err != nil {
		panic(err)
	}
	instances := make(map[string]string)
	for _, v := range response.LoadBalancers.LoadBalancer {
		if v.LoadBalancerName != "" {
			instances[v.LoadBalancerId] = v.LoadBalancerName
		} else {
			instances[v.LoadBalancerId] = v.LoadBalancerId
		}
	}
	s.instances = instances
}

func (s *SlbExporter) InitGauge() {
	s.NewClient()
	s.GetInstance()
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			s.GetInstance()
		}
	}()
	s.metrics = map[string]*prometheus.GaugeVec{}
	for k, v := range Layer4And7Metrics {
		s.metrics[k] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "aliyun_slb",
			Name:      v,
		}, []string{"instance_id", "instance_name", "vip", "protocol", "port"})
	}
}

func (s *SlbExporter) GetMetric(metricName string) {
	var dimensions []map[string]string
	if instanceIds != "" {
		for _, v := range strings.Split(instanceIds, ",") {
			d := map[string]string{"instanceId": v}
			dimensions = append(dimensions, d)
		}
	} else {
		for k := range s.instances {
			d := map[string]string{"instanceId": k}
			dimensions = append(dimensions, d)
		}
	}
	dimension, err := json.Marshal(dimensions)
	if err != nil {
		log.Println(err)
	}
	request := cms.CreateDescribeMetricLastRequest()
	request.MetricName = metricName
	request.Namespace = PROJECT
	request.Dimensions = string(dimension)
	request.Period = "120"
	response, err := s.client.DescribeMetricLast(request)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal([]byte(response.Datapoints), &s.DataPoints)
	if err != nil {
		log.Println(err)
	}
}

func (s SlbExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, v := range s.metrics {
		v.Describe(ch)
	}
}

func (s SlbExporter) Collect(ch chan<- prometheus.Metric) {
	for m := range Layer4And7Metrics {
		s.GetMetric(m)
		for _, v := range s.DataPoints {
			s.metrics[m].With(prometheus.Labels{
				"instance_id":   v.InstanceId,
				"instance_name": s.instances[v.InstanceId],
				"vip":           v.Vip,
				"protocol":      v.Protocol,
				"port":          v.Port,
			}).Set(v.Average)
		}
		time.Sleep(34 * time.Millisecond)
	}
	for _, m := range s.metrics {
		m.Collect(ch)
	}
}
