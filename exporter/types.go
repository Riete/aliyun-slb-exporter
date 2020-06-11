package exporter

import (
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	accessKeyId     = os.Getenv("ACCESS_KEY_ID")
	accessKeySecret = os.Getenv("ACCESS_KEY_SECRET")
	regionId        = os.Getenv("REGION_ID")
	instanceIds     = os.Getenv("INSTANCE_ID")
)

const PROJECT string = "acs_slb_dashboard"

var Layer4And7Metrics = map[string]string{
	// layer 4
	"HeathyServerCount":          "heathy_server_count",          //后端健康ECS实例个数
	"UnhealthyServerCount":       "unhealthy_server_count",       //后端异常ECS实例个数
	"PacketTX":                   "packet_tx",                    //端口每秒流出数据包数
	"PacketRX":                   "packet_rx",                    //端口每秒流入数据包数
	"TrafficRXNew":               "traffic_rx_new",               //端口每秒流入数据量
	"TrafficTXNew":               "traffic_tx_new",               //端口每秒流出数据量
	"ActiveConnection":           "active_connection",            //端口当前活跃连接数，既客户端正在访问SLB产生的连接
	"InactiveConnection":         "inactive_connection",          //端口当前非活跃连接数，既访问SLB后未断开的空闲的连接
	"NewConnection":              "new_connection",               //端口当前新建连接数
	"MaxConnection":              "max_connection",               //端口当前并发连接数
	"DropConnection":             "drop_connection",              //监听每秒丢失连接数
	"DropPacketRX":               "drop_packet_rx",               //监听每秒丢失入包数
	"DropPacketTX":               "drop_packet_tx",               //监听每秒丢失出包数
	"DropTrafficRX":              "drop_traffic_rx",              //监听每秒丢失入bit数
	"DropTrafficTX":              "drop_traffic_tx",              //监听每秒丢失出bit数
	"InstanceActiveConnection":   "instance_active_connection",   //实例每秒活跃连接数
	"InstanceDropConnection":     "instance_drop_connection",     //实例每秒丢失连接数
	"InstanceDropPacketRX":       "instance_drop_packet_rx",      //实例每秒丢失入包数
	"InstanceDropPacketTX":       "instance_drop_packet_tx",      //实例每秒丢失出包数
	"InstanceDropTrafficRX":      "instance_drop_traffic_rx",     //实例每秒丢失入bit数
	"InstanceDropTrafficTX":      "instance_drop_traffic_tx",     //实例每秒丢失出bit数
	"InstanceInactiveConnection": "instance_inactive_connection", //实例每秒非活跃连接数
	"InstanceMaxConnection":      "instance_max_connection",      //实例每秒最大并发连接数
	"InstanceNewConnection":      "instance_new_connection",      //实例每秒新建连接数
	"InstancePacketRX":           "instance_packet_rx",           //实例每秒入包数
	"InstancePacketTX":           "instance_packet_tx",           //实例每秒出包数
	"InstanceTrafficRX":          "instance_traffic_rx",          //实例每秒入bit数
	"InstanceTrafficTX":          "instacne_traffic_tx",          //实例每秒出bit数
	// layer 7
	"Qps":                     "qps",                        //端口维度的QPS
	"Rt":                      "rt",                         //端口维度的请求平均延时
	"StatusCode2xx":           "status_code_2xx",            //端口维度的slb返回给client的2xx状态码统计
	"StatusCode3xx":           "status_code_3xx",            //端口维度的slb返回给client的3xx状态码统计
	"StatusCode4xx":           "status_code_4xx",            //端口维度的slb返回给client的4xx状态码统计
	"StatusCode5xx":           "status_code_5xx",            //端口维度的slb返回给client的5xx状态码统计
	"StatusCodeOther":         "status_code_other",          //端口维度的slb返回给client的other状态码统计
	"UpstreamCode4xx":         "upstream_code_4xx",          //端口维度的rs返回给slb的4xx状态码统计
	"UpstreamCode5xx":         "upstream_code_5xx",          //端口维度的rs返回给client的5xx状态码统计
	"UpstreamRt":              "upstream_rt",                //端口维度的rs发给proxy的平均请求延迟
	"InstanceQps":             "instance_qps",               //实例维度的QPS
	"InstanceRt":              "instance_rt",                //实例维度的请求平均延时
	"InstanceStatusCode2xx":   "instance_status_code_2xx",   //实例维度的slb返回给client的2xx状态码统计
	"InstanceStatusCode3xx":   "instance_status_code_3xx",   //实例维度的slb返回给client的3xx状态码统计
	"InstanceStatusCode4xx":   "instance_status_code_4xx",   //实例维度的slb返回给client的4xx状态码统计
	"InstanceStatusCode5xx":   "instance_status_code_5xx",   //实例维度的slb返回给client的5xx状态码统计
	"InstanceStatusCodeOther": "instance_status_code_other", //实例维度的slb返回给client的Other状态码统计
	"InstanceUpstreamCode4xx": "instance_upstream_code_4xx", //实例维度的rs返回给slb的4xx状态码统计
	"InstanceUpstreamCode5xx": "instance_upstream_code_5xx", //实例维度的rs返回给slb的5xx状态码统计
	"InstanceUpstreamRt":      "instance_upstream_rt",       //实例维度的rs发给proxy的平均请求延迟
}

type SlbExporter struct {
	client     *cms.Client
	DataPoints []struct {
		InstanceId string  `json:"instanceId"`
		Port       string  `json:"port,omitempty"`
		Protocol   string  `json:"protocol,omitempty"`
		Vip        string  `json:"vip,omitempty"`
		Average    float64 `json:"Average"`
	}
	metrics   map[string]*prometheus.GaugeVec
	instances map[string]string
}
