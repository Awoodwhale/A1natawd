package wdocker

import (
	"go_awd/pkg/util"
	"go_awd/pkg/wlog"
	"testing"
)

var cli *Docker

func init() {
	wlog.InitLogger("./", 5, 7)
	cli = NewDockerClient()
}

func TestDocker_PullImage(t *testing.T) {
	if err := cli.PullImage("hello-world:latest"); err != nil {
		t.Error(err)
	}
}

func TestDocker_BuildImage(t *testing.T) {
	if err := cli.BuildImage("./test.tar", "test:v1"); err != nil {
		t.Error(err)
	}
}

func TestDocker_CreateContainer(t *testing.T) {
	if _, err := cli.CreateContainerWithSSH("test:v1",
		"test_flask",
		util.GenEnv(),
		"8080",
		"40001",
		"50001"); err != nil {
		t.Error(err)
	}
}

func TestDocker_RemoveContainer(t *testing.T) {
	if err := cli.RemoveContainer("e57004cbf8131692c46774e15ff1a453a882bc51889507650a3e5cd0be7b6978"); err != nil {
		t.Error(err)
	}
}
