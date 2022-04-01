package configuration

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configurator struct {
	l logrus.FieldLogger
}

func NewConfigurator(l logrus.FieldLogger) *Configurator {
	return &Configurator{l}
}

type Configuration struct {
	Transports []TransportConfiguration `yaml:"transports"`
}

type TransportConfiguration struct {
	Enabled            bool     `yaml:"enabled"`
	Source             uint32   `yaml:"source"`
	Departure          uint32   `yaml:"departure"`
	Transport          []uint32 `yaml:"transport"`
	Arrival            uint32   `yaml:"arrival"`
	Destination        uint32   `yaml:"destination"`
	OpenGateDuration   uint32   `yaml:"open_gate_duration"`
	ClosedGateDuration uint32   `yaml:"closed_gate_duration"`
	RideDuration       uint32   `yaml:"ride_duration"`
}

func (c *Configurator) GetConfiguration() (*Configuration, error) {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		c.l.Printf("yamlFile.Get err   #%v ", err)
		return nil, err
	}

	con := &Configuration{}
	err = yaml.Unmarshal(yamlFile, con)
	if err != nil {
		c.l.Fatalf("Unmarshal: %v", err)
		return nil, err
	}

	return con, nil
}

func (c Configuration) GetTransportConfiguration(index byte) (*TransportConfiguration, error) {
	if len(c.Transports) > 0 && index < byte(len(c.Transports)) {
		w := &TransportConfiguration{}
		w = &c.Transports[index]
		return w, nil
	}
	return nil, errors.New(fmt.Sprintf("Index out of bounds: %d", index))
}
