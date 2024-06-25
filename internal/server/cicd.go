package server

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/dollarkillerx/harbor_easy_cicd/internal/conf"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/utils"
	"github.com/rs/zerolog/log"
)

var mu sync.Mutex

func (s *Server) cicd(hk harborHook) {
	mu.Lock()
	defer mu.Unlock()

	for _, i := range s.conf.Tasks {
		if i.HarborKey == hk.EventData.Repository.Name {
			s.noticeLog(i.HarborKey, i.TaskName, "获取到任务")
			s.cicdLogic(i, hk)
		}
	}
}

func (s *Server) cicdLogic(task conf.Task, hk harborHook) {
	dockerImg := fmt.Sprintf("%s/%s", s.conf.HarborAddress, strings.Split(hk.EventData.Resources[0].ResourceUrl, "/")[2])
	composeFile := fmt.Sprintf("%s/%s", task.Path, "docker-compose.yaml")
	file, err := os.ReadFile(composeFile)
	if err != nil {
		log.Error().Msgf("Cicd Error: 获取目录不存在 %s", composeFile)
		s.noticeLog(task.HarborKey, task.TaskName, fmt.Sprintf("cicd Error: 获取目录不存在 %s", composeFile))
		return
	}

	// 如果不存在 则 更新新的 tag
	if !strings.Contains(string(file), dockerImg) {
		utils.ReplaceImage(composeFile, dockerImg)
	}

	err = os.Chdir(task.Path)
	if err != nil {
		log.Error().Msgf("Cicd Error: 获取目录不存在 %s", task.Path)
		s.noticeLog(task.HarborKey, task.TaskName, fmt.Sprintf("cicd Error: 获取目录不存在 %s", task.Path))
		return
	}

	// ls コマンドを実行
	cmd := exec.Command(task.Cmd)

	err = cmd.Run()
	if err != nil {
		log.Error().Msgf("Cicd Error: 执行错误 %s", err)
		s.noticeLog(task.HarborKey, task.TaskName, fmt.Sprintf("cicd Error: 执行错误 %s", err))
		return
	}

	task.Heartbeat = strings.TrimSpace(task.Heartbeat)
	if task.Heartbeat != "" {
		resp, err := http.Get(task.Heartbeat)
		if err != nil {
			log.Error().Msgf("Cicd Error: 执行错误 Heartbeat 验证失败 %s", task.Heartbeat)
			s.noticeLog(task.HarborKey, task.TaskName, fmt.Sprintf("cicd Error: Heartbeat验证失败 %s", task.Heartbeat))
			return
		}
		if resp.StatusCode != 200 {
			log.Error().Msgf("Cicd Error: 执行错误 Heartbeat 验证失败 %s", task.Heartbeat)
			s.noticeLog(task.HarborKey, task.TaskName, fmt.Sprintf("cicd Error: Heartbeat验证失败 %s", task.Heartbeat))
			return
		}

		s.noticeLog(task.HarborKey, task.TaskName, fmt.Sprintf("success %s", task.Heartbeat))

		return
	}

	s.noticeLog(task.HarborKey, task.TaskName, "success")
}
