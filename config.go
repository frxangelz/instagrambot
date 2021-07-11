/* config.go
Instagram Bot - autofollow, like, comment dan unfollow not followback (c) Free Angel - frxangelz@gmail.com
please subscribe to my channel :
https://www.youtube.com/channel/UC15iFd0nlfG_tEBrt6Qz1NQ
*/

package main

import (
	//"encoding/json"
	//"os"
	//	"fmt"
	"strconv"
	"strings"

	"github.com/alyu/configparser"
	"github.com/jimlawless/whereami"
)

var (
	config Configuration
	conf   *configparser.Configuration
)

func trims(str string) string {

	s := strings.Replace(str, "\r", "", -1)
	return strings.Replace(s, "\n", "", -1)
}

func load_config(fname string) error {
	var err error = nil

	conf, err = configparser.Read(fname)
	if err != nil {
		return err
	}

	section, err := conf.Section("main")
	if err != nil {
		return err
	}

	stmp := section.ValueOf("debug_to_console")
	if stmp == "1" {
		config.debug_to_console = true
	} else {
		config.debug_to_console = false
	}

	stmp = section.ValueOf("debug_to_file")
	if stmp == "1" {
		config.debug_to_file = true
	} else {
		config.debug_to_file = false
	}

	stmp = section.ValueOf("max_log_lines")
	config.max_log_lines, _ = strconv.Atoi(stmp)
	if config.max_log_lines < 10 {
		config.max_log_lines = 1000
	}

	config.synch_interval = strToInt64(section.ValueOf("synch_interval"))
	if (config.synch_interval != 0) && (config.synch_interval < 15) {
		mylog(whereami.WhereAmI(), "synch_interval too fast, skip it")
		config.synch_interval = 0
	}

	return nil
}

func save_config(fname string) error {

	return configparser.Save(conf, fname)
}

type Configuration struct {
	debug_to_file    bool
	debug_to_console bool
	max_log_lines    int
	synch_interval   int64
}
