package config

const (
	RMQADDR      = "amqp://guest:guest@172.17.84.205:5672/"
	EXCHANGENAME = "syslog_direct"
	CONSUMERCNT  = 4
)

var (
	RoutingKeys [4]string = [4]string{"info", "debug", "warn", "error"}
)
