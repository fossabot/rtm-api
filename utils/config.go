package utils

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

const (
	// AppName is Application name
	AppName = "rtm-api"
	// APIVersion is API version
	APIVersion = "v0"
	// BuildVersion is API build version
	BuildVersion = "v0.3.0"

	MAX_MESSAGE_SIZE = 8192
)

var (
	cfg           *config = NewConfig()
	IsShowVersion bool
)

type config struct {
	Version           string
	HttpPort          string `yaml:"httpPort"`
	Profiling         bool
	ErrorLogging      bool `yaml:"errorLogging"`
	Logging           *Logging
	Realtime          *RealtimeSetting
	MessagingProvider string `yaml:"messagingProvider"`
	NSQ               *NSQ
	Kafka             *Kafka
}

type Logging struct {
	Level string
}

type RealtimeSetting struct {
	IsDisplayConnectionInfo bool
}

type NSQ struct {
	Port           string
	NsqlookupdHost string
	NsqlookupdPort string
	NsqdHost       string
	NsqdPort       string
	Topic          string
	Channel        string
}

type Kafka struct {
	Port           string
	NsqlookupdHost string
	NsqlookupdPort string
	NsqdHost       string
	NsqdPort       string
	Topic          string
	Channel        string
}

func NewConfig() *config {
	log.SetFlags(log.Llongfile)

	logging := &Logging{
		Level: "development",
	}

	realtimeSetting := &RealtimeSetting{}
	nsq := &NSQ{}
	kafka := &Kafka{}

	c := &config{
		Version:           "0",
		HttpPort:          "8102",
		Profiling:         false,
		ErrorLogging:      false,
		Logging:           logging,
		Realtime:          realtimeSetting,
		MessagingProvider: "",
		NSQ:               nsq,
		Kafka:             kafka,
	}

	c.LoadYaml()
	c.LoadEnvironment()
	c.ParseFlag()

	return c
}

func GetConfig() *config {
	return cfg
}

func (c *config) LoadYaml() {
	buf, _ := ioutil.ReadFile("config/app.yaml")
	yaml.Unmarshal(buf, c)
}

func (c *config) LoadEnvironment() {
	var v string

	if v = os.Getenv("HTTP_PORT"); v != "" {
		c.HttpPort = v
	}
	if v = os.Getenv("SC_PORT"); v != "" {
		c.HttpPort = v
	}
	if v = os.Getenv("SC_PROFILING"); v != "" {
		if v == "true" {
			c.Profiling = true
		} else if v == "false" {
			c.Profiling = false
		}
	}
	if v = os.Getenv("SC_ERROR_LOGGING"); v != "" {
		if v == "true" {
			c.ErrorLogging = true
		} else if v == "false" {
			c.ErrorLogging = false
		}
	}

	// Logging
	if v = os.Getenv("SC_LOGGING_LEVEL"); v != "" {
		c.Logging.Level = v
	}
}

func (c *config) ParseFlag() {
	flag.BoolVar(&IsShowVersion, "v", false, "")
	flag.BoolVar(&IsShowVersion, "version", false, "show version")

	flag.StringVar(&c.HttpPort, "httpPort", c.HttpPort, "")

	var profiling string
	flag.StringVar(&profiling, "profiling", "", "")

	var demoPage string
	flag.StringVar(&demoPage, "demoPage", "", "false")

	var errorLogging string
	flag.StringVar(&errorLogging, "errorLogging", "", "false")

	// Logging
	flag.StringVar(&c.Logging.Level, "logging.level", c.Logging.Level, "")

	var isDisplayConnectionInfo string
	flag.StringVar(&isDisplayConnectionInfo, "isDisplayConnectionInfo", "", "Display connection info.")
	if profiling == "true" {
		c.Realtime.IsDisplayConnectionInfo = true
	} else if profiling == "false" {
		c.Realtime.IsDisplayConnectionInfo = false
	}

	flag.StringVar(&c.NSQ.NsqlookupdHost, "nsqlookupdHost", c.NSQ.NsqlookupdHost, "Host name of nsqlookupd")
	flag.StringVar(&c.NSQ.NsqlookupdPort, "nsqlookupdPort", c.NSQ.NsqlookupdPort, "Port no of nsqlookupd")
	flag.StringVar(&c.NSQ.NsqdHost, "nsqdHost", c.NSQ.NsqdHost, "Host name of nsqd")
	flag.StringVar(&c.NSQ.NsqdPort, "nsqdPort", c.NSQ.NsqdPort, "Port no of nsqd")
	flag.StringVar(&c.NSQ.Topic, "topic", c.NSQ.Topic, "Topic name")
	flag.StringVar(&c.NSQ.Channel, "channel", c.NSQ.Channel, "Channel name. If it's not set, channel is hostname.")
	flag.Parse()
}
