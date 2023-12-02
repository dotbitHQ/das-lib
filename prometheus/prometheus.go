package prometheus

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"net"
	"sync"
	"time"
)

var (
	log          = logger.NewLogger("prometheus", logger.LevelInfo)
	promRegister = prometheus.NewRegistry()
)

type Option func(*Prometheus)

type Prometheus struct {
	pusher      *push.Pusher
	Metrics     Metric
	env         common.DasNetType
	pushGateway string
	serverName  string
}

type Metric struct {
	l         sync.Mutex
	api       *prometheus.SummaryVec
	errNotify *prometheus.CounterVec
}

func New(opts ...Option) *Prometheus {
	p := &Prometheus{}
	for _, opt := range opts {
		opt(p)
	}
	if p.pushGateway == "" || p.serverName == "" {
		panic("push gateway or server name is empty")
	}
	if p.env != common.DasNetTypeMainNet &&
		p.env != common.DasNetTypeTestnet2 &&
		p.env != common.DasNetTypeTestnet3 {
		panic("env only can include (1|2|3)")
	}
	p.pusher = push.New(p.pushGateway, p.serverName)
	return p
}

func WithEnv(env common.DasNetType) Option {
	return func(p *Prometheus) {
		p.env = env
	}
}

func WithPushGateway(pushGateway string) Option {
	return func(p *Prometheus) {
		p.pushGateway = pushGateway
	}
}

func WithServerName(serverName string) Option {
	return func(p *Prometheus) {
		p.serverName = serverName
	}
}

func (m *Metric) Api() *prometheus.SummaryVec {
	if m.api == nil {
		m.l.Lock()
		defer m.l.Unlock()
		m.api = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: "api",
		}, []string{"method", "http_status", "err_no", "err_msg"})
		promRegister.MustRegister(m.api)
	}
	return m.api
}

func (m *Metric) ErrNotify() *prometheus.CounterVec {
	if m.errNotify == nil {
		m.l.Lock()
		defer m.l.Unlock()
		m.errNotify = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "notify",
		}, []string{"title", "text"})
		promRegister.MustRegister(m.errNotify)
	}
	return m.errNotify
}

func (t *Prometheus) Run() {
	t.pusher.Gatherer(promRegister)
	t.pusher.Grouping("env", fmt.Sprint(t.env))
	t.pusher.Grouping("instance", GetLocalIp("eth0"))

	go func() {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()

		for range ticker.C {
			_ = t.pusher.Push()
		}
	}()
}

func GetLocalIp(interfaceName string) string {
	ief, err := net.InterfaceByName(interfaceName)
	if err != nil {
		log.Error("GetLocalIp: ", err)
		return ""
	}
	addrs, err := ief.Addrs()
	if err != nil {
		log.Error("GetLocalIp: ", err)
		return ""
	}

	var ipv4Addr net.IP
	for _, addr := range addrs {
		if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			break
		}
	}
	if ipv4Addr == nil {
		log.Errorf("GetLocalIp interface %s don't have an ipv4 address", interfaceName)
		return ""
	}
	return ipv4Addr.String()
}
