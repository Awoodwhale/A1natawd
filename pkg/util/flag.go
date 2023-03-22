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

func GenFlagEnv() string {
	flagEnv := conf.FlagEnv
	if flagEnv == "" {
		flagEnv = "FLAG"
	}
	return fmt.Sprintf("%v=%v", flagEnv, GenFlag())
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

func getPortNotInUse(maxCount, startPort int) int {
	rand.Seed(time.Now().Unix())
	startScanPort := rand.Intn(maxCount) + startPort // 生成指定范围的随机数[0, MaxTeamCount)
	for ; startScanPort <= maxCount+startPort; startScanPort++ {
		addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(startScanPort))
		conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
		if err != nil {
			return startScanPort
		}
		_ = conn.Close()
	}
	startScanPort = startPort
	for ; startScanPort <= maxCount+startPort; startScanPort++ {
		addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(startScanPort))
		conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
		if err != nil {
			return startScanPort
		}
		_ = conn.Close()
	}
	return -1
}
