package config

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"time"
)

const (
	defaultMaxHeaderBytes = 1 << 20

	ParamNamePort              = "PORT"
	ParamNameHTTPTimeout       = "HTTP_TIMEOUT"
	ParamNameMaxHeaderBytes    = "MAX_HEADER_BYTES"
	ParamWaitBeforeProcessTask = "WAIT_BEFORE_PROCESS_TASK"
	ParamTaskTimeout           = "TASK_TIMEOUT"
	ParamLogPath               = "LOG_PATH"
)

var (
	ErrParamNotFound = errors.New("param not found")
	ErrBadConfig     = errors.New("bad config")
)

type Config struct {
	params map[string]string
}

func NewConfig() *Config {
	return &Config{
		params: map[string]string{
			ParamNamePort:              "8080",
			ParamNameHTTPTimeout:       "10s",
			ParamNameMaxHeaderBytes:    strconv.Itoa(defaultMaxHeaderBytes),
			ParamWaitBeforeProcessTask: "5",
			ParamTaskTimeout:           "20",
			ParamLogPath:               "/tmp/log.txt",
		},
	}
}

func (c *Config) InitFromEnv() error {
	for _, k := range c.getParamNames() {
		v := os.Getenv(k)
		if v != "" {
			err := c.setParam(k, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Config) setParam(key, value string) error {
	if _, ok := c.params[key]; !ok {
		return fmt.Errorf("%w: %s", ErrParamNotFound, key)
	}
	c.params[key] = value
	return nil
}

func (c *Config) getParamNames() []string {
	return slices.Collect(maps.Keys(c.params))
}

func (c *Config) GetPort() string {
	return c.params[ParamNamePort]
}

func (c *Config) GetTimeout() string {
	return c.params[ParamNameHTTPTimeout]
}

func (c *Config) GetMaxHeaderBytes() string {
	return c.params[ParamNameMaxHeaderBytes]
}

func (c *Config) GetWaitBeforeProcessTask() time.Duration {
	s, _ := strconv.Atoi(c.params[ParamWaitBeforeProcessTask])
	return time.Duration(s) * time.Second
}

func (c *Config) GetTaskTimeout() time.Duration {
	s, _ := strconv.Atoi(c.params[ParamTaskTimeout])
	return time.Duration(s) * time.Second
}

func (c *Config) GetLogPath() string {
	return c.params[ParamLogPath]
}
