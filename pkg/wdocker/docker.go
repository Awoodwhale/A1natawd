package wdocker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"go_awd/pkg/wlog"
	"io"
	"os"
)

type Docker struct {
	*client.Client
}

func NewDockerClient() *Docker {
	// TODO: 使用 DockerServerIP 进行docker服务
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		wlog.DockerLogger.Errorln("Docker Client create error,", err.Error())
		return nil
	}
	return &Docker{cli}
}

func (d *Docker) CheckImageExist(imageName string) bool {
	cli, ctx := d.Client, context.Background()
	_, _, err := cli.ImageInspectWithRaw(ctx, imageName)
	if err == nil {
		return true
	}
	return false
}

func (d *Docker) PullImage(imageName string) error {
	cli, ctx := d.Client, context.Background()
	if d.CheckImageExist(imageName) {
		wlog.DockerLogger.Infoln("Docker image has exist, image name:", imageName)
		return nil
	}
	resp, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})

	if err != nil {
		// 处理错误
		wlog.DockerLogger.Errorln("Docker pull image error,", err.Error())
		return err
	}

	if info, err := io.ReadAll(resp); err == nil {
		wlog.DockerLogger.Infof("Docker pull image info:\n%v", string(info))
	}

	if err := resp.Close(); err != nil {
		wlog.DockerLogger.Errorln("Docker pull image resp close error,", err.Error())
	}

	return nil
}

func (d *Docker) BuildImage(path, imageName string) error {
	cli, ctx := d.Client, context.Background()
	dockerBuildContext, err := os.Open(path)
	if err != nil {
		wlog.DockerLogger.Errorln("Docker build image error,", err.Error())
		return nil
	}
	defer func(dockerBuildContext *os.File) {
		err := dockerBuildContext.Close()
		if err != nil {
			wlog.DockerLogger.Errorln("dockerBuildContext.Close error,", err.Error())
		}
	}(dockerBuildContext)
	opts := types.ImageBuildOptions{
		Tags:           []string{imageName},
		Memory:         10 * 1024 * 1024,  // 10MB
		MemorySwap:     256 * 1024 * 1024, // 256MB
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Dockerfile:     "Dockerfile",
	}
	resp, err := cli.ImageBuild(ctx, dockerBuildContext, opts)
	if err != nil {
		wlog.DockerLogger.Errorln("Docker build image error,", err.Error())
		return err
	}
	if info, err := io.ReadAll(resp.Body); err == nil {
		wlog.DockerLogger.Infof("Docker build image info:\n%v", string(info))
	}
	return nil
}

func (d *Docker) CreateContainerWithSSH(imageName, containerName string, containerEnv []string, innerPort, exposePort, exposeSSHPort string) (string, error) {
	cli, ctx := d.Client, context.Background()
	// 配置容器内部port与flag
	containerConfig := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			nat.Port(innerPort): {},
			nat.Port("22"):      {},
		},
		Env: containerEnv,
	}
	// 配置宿主机的host config
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(innerPort): []nat.PortBinding{ // bind challenge server port
				{
					HostIP:   "0.0.0.0",
					HostPort: exposePort,
				},
			},
			nat.Port("22"): []nat.PortBinding{ // bind ssh server port
				{
					HostIP:   "0.0.0.0",
					HostPort: exposeSSHPort,
				},
			},
		},
		// Container resource limits
		Resources: container.Resources{
			NanoCPUs: 1,
			Memory:   134217728,
		},
	}
	// create container
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		wlog.DockerLogger.Errorln("Docker create container error,", err.Error())
		return "", err
	}
	// open new container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		// if error happen, remove the container
		if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			wlog.DockerLogger.Errorln("Docker remove container error,", err.Error())
			return "", err
		}
		wlog.DockerLogger.Errorln("Docker run container error,", err.Error())
		return "", err
	}
	wlog.DockerLogger.Infoln("Docker run success container, id:", resp.ID)
	return resp.ID, nil
}

func (d *Docker) RemoveContainer(containerID string) error {
	cli, ctx := d.Client, context.Background()
	if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		wlog.DockerLogger.Errorln("Docker remove container error,", err.Error())
		return err
	}
	wlog.DockerLogger.Infoln("Docker remove container success, id:", containerID)
	return nil
}

func (d *Docker) RemoveImage(imageName string) error {
	cli, ctx := d.Client, context.Background()
	if exist := d.CheckImageExist(imageName); !exist {
		return nil // 不存在这个镜像
	}
	if _, err := cli.ImageRemove(ctx, imageName, types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	}); err != nil {
		wlog.DockerLogger.Errorln("Docker remove image error,", err.Error())
		return err
	}
	wlog.DockerLogger.Infoln("Docker remove image success, imageName:", imageName)
	return nil
}
