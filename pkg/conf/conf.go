// SPDX-License-Identifier: BSD-2
// Copyright 2021 VMware, Inc.

package conf

import (
	"github.com/photon-os-container-builder/pkg/log"
	"github.com/photon-os-container-builder/pkg/parser"
	"github.com/spf13/viper"
)

const (
	DefaultLogLevel       = "info"
	DefaultReleaseVersion = "4.0"
	DefaultStorageDir     = "/var/lib/machines"
	DefaultGPGDir         = "/etc/pki/rpm-gpg"
	DefaultPackages       = "systemd,dbus,iproute2,tdnf,photon-release,photon-repos,curl,shadow,ncurses-terminfo"

	DefaultParentLink  = "eth0"
	DefaultNetworkKind = "macvlan"
	DefaultAddressPool = "172.16.85.50/24 "
	DefaultPoolOffSet  = 64

	Version  = "0.1"
	ConfPath = "/etc/photon-os-container/"
	ConfFile = "photon-os-container"
)

// Config file key value
type Network struct {
	Kind        string `mapstructure:"Kind"`
	ParentLink  string `mapstructure:"ParentLink"`
	AddressPool string `mapstructure:"AddressPool"`
	PoolOffset  int    `mapstructure:"PoolOffset "`
}

type System struct {
	Packages string `mapstructure:"Packages"`

	Release  string `mapstructure:"Release"`
	LogLevel string `mapstructure:"LogLevel"`
}
type Config struct {
	Network Network `mapstructure:"Network"`
	System  System  `mapstructure:"System"`
}

func Parse() (*Config, error) {
	viper.SetConfigName(ConfFile)
	viper.AddConfigPath(ConfPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("%+v", err)
	}

	viper.SetDefault("Network.AddressPool", DefaultAddressPool)
	viper.SetDefault("Network.PoolOffset", DefaultAddressPool)
	viper.SetDefault("Network.ParentLink", DefaultParentLink)
	viper.SetDefault("Network.Kind", DefaultNetworkKind)

	viper.SetDefault("System.LogLevel", DefaultLogLevel)
	viper.SetDefault("System.Release", DefaultReleaseVersion)
	viper.SetDefault("System.Packages", DefaultPackages)

	c := Config{}
	if err := viper.Unmarshal(&c); err != nil {
		log.Warnf("Failed to parse config file: '/etc/photon-os-container/photon-os-container.toml'")
	}

	log.SetLevel(c.System.LogLevel)

	if _, err := parser.ParseIP(c.Network.AddressPool); err != nil {
		log.Debugf("Failed to parse address pool. Default will be used", err)
		c.Network.AddressPool = DefaultAddressPool
	}

	return &c, nil
}