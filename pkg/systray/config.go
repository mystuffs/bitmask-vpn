// Copyright (C) 2018 LEAP
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package systray

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"0xacab.org/leap/bitmask-vpn/pkg/config"
	"golang.org/x/text/message"
)

const (
	oneDay   = time.Hour * 24
	oneMonth = oneDay * 30
)

var (
	configPath = path.Join(config.Path, "systray.json")
)

// Config holds the configuration of the systray
type Config struct {
	LastNotification time.Time
	Donated          time.Time
	SelectGateway    bool
	UserStoppedVPN   bool
	Version          string           `json:"-"`
	Printer          *message.Printer `json:"-"`
}

// ParseConfig reads the configuration from the configuration file
func ParseConfig() *Config {
	var conf Config

	f, err := os.Open(configPath)
	if err != nil {
		conf.save()
		return &conf
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&conf)
	return &conf
}

func (c *Config) setUserStoppedVPN(vpnStopped bool) error {
	c.UserStoppedVPN = vpnStopped
	return c.save()
}

func (c *Config) hasDonated() bool {
	return c.Donated.Add(oneMonth).After(time.Now())
}

func (c *Config) needsNotification() bool {
	return !c.hasDonated() && c.LastNotification.Add(oneDay).Before(time.Now())
}

func (c *Config) setNotification() error {
	c.LastNotification = time.Now()
	return c.save()
}

func (c *Config) setDonated() error {
	c.Donated = time.Now()
	return c.save()
}

func (c *Config) save() error {
	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(c)
}
