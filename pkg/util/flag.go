package util

import (
	"fmt"
	"github.com/docker/distribution/uuid"
	"go_awd/conf"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type Env struct {
	Flag        string
	SSHUsername string
	SSHPassword string
}

func (e *Env) toString() []string {
	return []string{e.Flag, e.SSHUsername, e.SSHPassword}
}

type Option func(d *Env)

func WithFlag(flag string) Option {
	return func(e *Env) {
		e.Flag = fmt.Sprintf("%v=%v", conf.FlagEnv, flag)
	}
}

func WithSSHUsername(username string) Option {
	return func(e *Env) {
		e.SSHUsername = fmt.Sprintf("%v=%v", conf.SSHUsernameEnv, username)
	}
}

func WithSSHPassword(pwd string) Option {
	return func(e *Env) {
		e.SSHPassword = fmt.Sprintf("%v=%v", conf.SSHPasswordEnv, pwd)
	}
}

func GenEnv(option ...Option) []string {
	env := &Env{}
	for _, o := range option {
		o(env)
	}
	if env.Flag == "" {
		WithFlag(GenFlag())(env)
	}
	if env.SSHUsername == "" {
		WithSSHUsername(conf.SSHDefaultUsername)(env)
	}
	if env.SSHPassword == "" {
		WithSSHPassword(conf.SSHDefaultPassword)(env)
	}
	return env.toString()
}

func GenFlag() string {
	flagPrefix := conf.FlagPrefix
	if flagPrefix == "" {
		flagPrefix = "flag"
	}
	return fmt.Sprintf("%v{%v}", flagPrefix, uuid.Generate().String())
}

func GetPwnPortNotInUse() int {
	return getPortNotInUse(conf.MaxTeamCount, conf.PwnBoxStartPort)
}

func GetWebPortNotInUse() int {
	return getPortNotInUse(conf.MaxTeamCount, conf.WebBoxStartPort)
}

func GetSSHPortNotInUse() int {
	return getPortNotInUse(conf.MaxTeamCount, conf.SSHStartPort)
}

func getPortNotInUse(maxCount, startPort int) int {
	rand.Seed(time.Now().Unix())
	startScanPort := rand.Intn(maxCount) + startPort // 生成指定范围的随机数[0, MaxTeamCount)
	for ; startScanPort <= maxCount+startPort; startScanPort++ {
		addr := net.JoinHostPort(conf.DockerServerIP, strconv.Itoa(startScanPort))
		conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
		if err != nil {
			return startScanPort
		}
		_ = conn.Close()
	}
	startScanPort = startPort
	for ; startScanPort <= maxCount+startPort; startScanPort++ {
		addr := net.JoinHostPort(conf.DockerServerIP, strconv.Itoa(startScanPort))
		conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
		if err != nil {
			return startScanPort
		}
		_ = conn.Close()
	}
	return -1
}
