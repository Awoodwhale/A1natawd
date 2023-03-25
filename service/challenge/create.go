package challenge

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_awd/cache"
	"go_awd/conf"
	"go_awd/dao"
	"go_awd/model"
	"go_awd/pkg/e"
	"go_awd/pkg/util"
	"go_awd/pkg/wdocker"
	wjwt "go_awd/pkg/wjwt"
	"go_awd/serializer"
	"go_awd/service"
	"strconv"
	"strings"
	"time"
)

// CreateOrUpdateChallenge
// @Description: 创建或者更新题目，包含镜像
// @receiver s *CreateOrUpdateChallengeImageService
// @param c *gin.Context
// @return serializer.Response
func (s *CreateOrUpdateChallengeImageService) CreateOrUpdateChallenge(c *gin.Context) serializer.Response {
	if s.BaseScore == 0 {
		s.BaseScore = 10 // 默认10分
	}
	cli := wdocker.NewDockerClient()
	chalDao := dao.NewChallengeDao(c)
	// 自定义上传tar包build镜像
	if file, err := c.FormFile("file"); err == nil { // 存在tar文件
		dockerTarPath := fmt.Sprintf("./docker/%v.tar", s.Title)
		imageName := fmt.Sprintf("%v-%v", s.Title, time.Now().Unix())
		if err := c.SaveUploadedFile(file, dockerTarPath); err != nil {
			return serializer.RespCode(e.InvalidWithUploadFile, c) // 上传文件失败
		}
		if err := chalDao.CreateOrUpdateChallenge(&model.Challenge{
			Title:           s.Title,
			Info:            s.Info,
			BaseScore:       s.BaseScore,
			InnerServerPort: s.InnerServerPort,
			ImageName:       imageName,
			Type:            s.Type,
			State:           "building",
		}); err != nil {
			service.Errorln("CreateOrUpdateChallenge dao error,", err.Error())
			return serializer.RespCode(e.InvalidWithCreateChallenge, c) // 创建题目失败
		}
		go func() { // 开个协程去build image
			chal := &model.Challenge{
				Title:           s.Title,
				Info:            s.Info,
				BaseScore:       s.BaseScore,
				InnerServerPort: s.InnerServerPort,
				ImageName:       imageName,
				Type:            s.Type,
				State:           "success",
			}
			if err := cli.BuildImage(dockerTarPath, imageName); err != nil {
				service.Errorln("CreateOrUpdateChallenge BuildImage error,", err.Error())
				chal.State = "error" //  build失败了
			}
			if !cli.CheckImageExist(imageName) {
				chal.State = "error" //  镜像不存在说明build失败了
			}
			// 更新state
			if err := dao.NewChallengeDaoByDB(chalDao.DB).CreateOrUpdateChallenge(chal); err != nil {
				service.Errorln("CreateOrUpdateChallenge BuildImage UpdateByTitle error,", err.Error())
				return
			}
		}()
		return serializer.RespSuccess(e.SuccessWithUploadChallenge, nil, c)
	}
	// 使用dockerhub链接的方式pull镜像
	if s.ImageName != "" {
		if err := chalDao.CreateOrUpdateChallenge(&model.Challenge{
			Title:           s.Title,
			Info:            s.Info,
			BaseScore:       s.BaseScore,
			InnerServerPort: s.InnerServerPort,
			ImageName:       s.ImageName,
			Type:            s.Type,
			State:           "building",
		}); err != nil {
			service.Errorln("CreateOrUpdateChallenge dao error,", err.Error())
			return serializer.RespCode(e.InvalidWithCreateChallenge, c) // 创建题目失败
		}
		go func() { // 开个协程去pull image
			chal := &model.Challenge{
				Title:           s.Title,
				Info:            s.Info,
				BaseScore:       s.BaseScore,
				InnerServerPort: s.InnerServerPort,
				ImageName:       s.Title,
				Type:            s.Type,
				State:           "success",
			}
			if err := cli.PullImage(s.ImageName); err != nil {
				service.Errorln("CreateOrUpdateChallenge PullImage error,", err.Error())
				chal.State = "error"
			}
			if !cli.CheckImageExist(s.ImageName) {
				chal.State = "error" //  镜像不存在说明pull失败了
			}
			// 更新state
			if err := dao.NewChallengeDaoByDB(chalDao.DB).CreateOrUpdateChallenge(chal); err != nil {
				service.Errorln("CreateOrUpdateChallenge PullImage UpdateByTitle error,", err.Error())
				return
			}
		}()
		return serializer.RespSuccess(e.SuccessWithUploadChallenge, nil, c)
	}
	return serializer.RespCode(e.InvalidWithCreateChallenge, c)
}

// StartTestChallenge
// @Description: 开启题目测试容器
// @receiver s *EmptyService
// @param c *gin.Context
// @param id int64
// @return serializer.Response
func (s *EmptyService) StartTestChallenge(c *gin.Context, id int64) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)

	// 如果这个题目已经当前这个管理员开启了一次了，就不能再开了，返回已开启的信息
	if exist, err := cache.RedisClient.SIsMember(cache.AdminStartTestChallengeKey(claims.ID), id).Result(); err != nil || exist {
		containerInfo, err := cache.RedisClient.HGetAll(cache.AdminContainerKey(claims.ID, id)).Result()
		if err != nil {
			return serializer.RespCode(e.InvalidWithContainerInfoLost, c) // 容器信息丢失
		}
		res := gin.H{}
		for k, v := range containerInfo {
			res[k] = v
		}
		env, _ := containerInfo["env"]
		res["env"] = strings.Split(env, "\n")
		return serializer.RespSuccess(e.SuccessWithFindStartedTestChallenge, res, c)
	}

	chalDao := dao.NewChallengeDao(c)
	chal, err := chalDao.GetByID(id)
	if err != nil {
		return serializer.RespCode(e.InvalidWithNotExistChallenge, c)
	}
	if chal.State != "success" {
		return serializer.RespCode(e.InvalidWithNotSuccessChallenge, c)
	}
	var challengePort string
	if chal.Type == "pwn" {
		challengePort = strconv.Itoa(util.GetPwnPortNotInUse())
	} else {
		challengePort = strconv.Itoa(util.GetWebPortNotInUse())
	}
	sshPort := strconv.Itoa(util.GetSSHPortNotInUse()) // 暴露出去容器的ssh端口
	sshUname := conf.SSHDefaultUsername                // ssh的username
	sshPwd := strconv.FormatInt(model.GenID(), 10)     // ssh的password
	containerEnv := util.GenEnv(                       // 生成env
		util.WithSSHUsername(sshUname),
		util.WithSSHPassword(sshPwd),
	)
	go func() {
		cli := wdocker.NewDockerClient()
		containerName := fmt.Sprintf("%v-%v-%v", chal.Type, chal.Title, claims.ID)
		containerID, err := cli.CreateContainerWithSSH( // 开容器
			chal.ImageName, containerName, containerEnv, chal.InnerServerPort, challengePort, sshPort)
		if err != nil {
			service.Errorln("StartTestChallenge CreateContainerWithSSH error,", err.Error())
			return
		}
		// 把admin开启的这个题目ID放到redis中
		if err := cache.RedisClient.SAdd(cache.AdminStartTestChallengeKey(claims.ID), chal.ID).Err(); err != nil {
			service.Errorln("redis set AdminStartTestChallengeKey error,", err.Error())
			return
		}
		// 容器info
		containerInfo := map[string]any{
			"challenge_id": chal.ID,
			"container_id": containerID,
			"type":         chal.Type,
			"ip":           conf.DockerServerIP,
			"port":         challengePort,
			"ssh_port":     sshPort,
			"ssh_username": sshUname,
			"ssh_password": sshPwd,
			"env":          strings.Join(containerEnv, "\n"),
		}
		// 如果开启容器成功，将当前这位管理员开启的容器info存入redis，同一道题一个管理员只能启动一个容器
		if err := cache.RedisClient.HMSet(cache.AdminContainerKey(claims.ID, chal.ID), containerInfo).Err(); err != nil {
			service.Errorln("redis set AdminContainerKey error,", err.Error())
			// 出错了要删除容器，同时要去把键给删了
			cache.RedisClient.SRem(cache.AdminStartTestChallengeKey(claims.ID), chal.ID)
			if err := cli.RemoveContainer(containerID); err != nil {
				service.Errorln("StartTestChallenge RemoveContainer error,", err.Error())
				return
			}
		}
		time.AfterFunc(conf.ContainerExistTime, func() { // 超时之后自动删除容器相关信息
			service.Debugln("Auto remove container...")
			cache.RedisClient.SRem(cache.AdminStartTestChallengeKey(claims.ID), chal.ID) // 删除管理员开启的题目记录
			cache.RedisClient.Del(cache.AdminContainerKey(claims.ID, chal.ID))           // 删除管理员开启的容器记录
			_ = cli.RemoveContainer(containerID)                                         // 删除container
		})
	}()

	return serializer.RespSuccess(
		e.SuccessWithStartTestChallenge,
		gin.H{
			"challenge_id": chal.ID,
			"type":         chal.Type,
			"ip":           conf.DockerServerIP,
			"port":         challengePort,
			"ssh_port":     sshPort,
			"ssh_username": sshUname,
			"ssh_password": sshPwd,
			"env":          containerEnv,
		}, c)
}
