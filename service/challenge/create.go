package challenge

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_awd/conf"
	"go_awd/dao"
	"go_awd/model"
	"go_awd/pkg/e"
	"go_awd/pkg/util"
	"go_awd/pkg/wdocker"
	"go_awd/serializer"
	"go_awd/service"
	"strconv"
	"time"
)

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
			return serializer.RespCode(e.InvalidWIthUploadFile, c) // 上传文件失败
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

func (s *EmptyService) StartTestChallenge(c *gin.Context, id int64) serializer.Response {
	chalDao := dao.NewChallengeDao(c)
	chal, err := chalDao.GetByID(id)
	if err != nil {
		return serializer.RespCode(e.InvalidWIthNotExistChallenge, c)
	}
	var challengePort string
	if chal.Type == "pwn" {
		challengePort = strconv.Itoa(util.GetPwnPortNotInUse())
	} else {
		challengePort = strconv.Itoa(util.GetWebPortNotInUse())
	}
	flagEnv := util.GenFlagEnv()
	go func() {
		cli := wdocker.NewDockerClient()
		if _, err := cli.CreateContainer( // 开容器
			chal.ImageName, chal.Type+"-"+chal.Title, flagEnv, chal.InnerServerPort, challengePort); err != nil {
			service.Errorln("StartTestChallenge CreateContainer error,", err.Error())
		}
	}()
	return serializer.RespSuccess(
		e.SuccessWithStartTestChallenge,
		gin.H{
			"ip":   conf.DockerServerIP,
			"port": challengePort,
			"env":  flagEnv,
		}, c)
}
