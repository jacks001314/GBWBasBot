package conf

import (
	hadoop_conf "common/proto/hadoop/hadoop_common/conf"
)

var (
	YARN_DEFAULT hadoop_conf.Resource = hadoop_conf.Resource{"yarn-default.xml", false}
	YARN_SITE    hadoop_conf.Resource = hadoop_conf.Resource{"yarn-site.xml", false}
)

const (
	YARN_PREFIX                      = "yarn."
	RM_PREFIX                        = YARN_PREFIX + "resourcemanager."
	RM_ADDRESS                       = RM_PREFIX + "address"
	DEFAULT_RM_ADDRESS               = "0.0.0.0:8032"
	RM_SCHEDULER_ADDRESS             = RM_PREFIX + "scheduler.address"
	DEFAULT_RM_SCHEDULER_ADDRESS     = "0.0.0.0:8030"
	RM_AM_EXPIRY_INTERVAL_MS         = YARN_PREFIX + "am.liveness-monitor.expiry-interval-ms"
	DEFAULT_RM_AM_EXPIRY_INTERVAL_MS = 600000
)

type yarn_configuration struct {
	conf hadoop_conf.Configuration
}

type YarnConfiguration interface {
	GetRMAddress() (string, error)
	GetRMSchedulerAddress() (string, error)

	SetRMAddress(address string) error
	SetRMSchedulerAddress(address string) error

	Get(key string, defaultValue string) (string, error)
	GetInt(key string, defaultValue int) (int, error)

	Set(key string, value string) error
	SetInt(key string, value int) error
}

func (yarn_conf *yarn_configuration) Get(key string, defaultValue string) (string, error) {
	return yarn_conf.conf.Get(key, defaultValue)
}

func (yarn_conf *yarn_configuration) GetInt(key string, defaultValue int) (int, error) {
	return yarn_conf.conf.GetInt(key, defaultValue)
}

func (yarn_conf *yarn_configuration) GetRMAddress() (string, error) {
	return yarn_conf.conf.Get(RM_ADDRESS, DEFAULT_RM_ADDRESS)
}

func (yarn_conf *yarn_configuration) GetRMSchedulerAddress() (string, error) {
	return yarn_conf.conf.Get(RM_SCHEDULER_ADDRESS, DEFAULT_RM_SCHEDULER_ADDRESS)
}

func (yarn_conf *yarn_configuration) Set(key string, value string) error {
	return yarn_conf.conf.Set(key, value)
}

func (yarn_conf *yarn_configuration) SetInt(key string, value int) error {
	return yarn_conf.conf.SetInt(key, value)
}

func (yarn_conf *yarn_configuration) SetRMAddress(address string) error {
	return yarn_conf.conf.Set(RM_ADDRESS, address)
}

func (yarn_conf *yarn_configuration) SetRMSchedulerAddress(address string) error {
	return yarn_conf.conf.Set(RM_SCHEDULER_ADDRESS, address)
}

func NewYarnConfiguration() (YarnConfiguration, error) {
	c, err := hadoop_conf.NewConfigurationResources([]hadoop_conf.Resource{YARN_DEFAULT, YARN_SITE})
	return &yarn_configuration{conf: c}, err
}
