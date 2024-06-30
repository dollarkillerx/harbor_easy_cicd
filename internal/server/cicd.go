package server

import (
	"fmt"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/models"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/utils"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex

func (s *Server) cicd(hk harborHook) {
	mu.Lock()
	defer mu.Unlock()

	var tasks []models.Task
	s.db.Model(&models.Task{}).Find(&tasks)

	for _, i := range tasks {
		if i.HarborKey == hk.EventData.Repository.Name {
			if i.Run {
				s.noticeLog(i.HarborKey, i.TaskName, "获取到任务")
				s.cicdLogic(i, hk)
			}
		}
	}
}

func (s *Server) initLog(task models.Task) uint {
	var log = models.TaskLogs{
		TaskId:   task.ID,
		TaskName: task.TaskName,
		Message:  "初始化任务",
	}
	s.db.Model(&models.TaskLogs{}).Create(&log)
	return log.ID
}

func (s *Server) log(id uint, success bool, message string) {
	s.db.Model(&models.TaskLogs{}).Where("id = ?", id).Updates(&models.TaskLogs{
		Success: success,
		Message: message,
	})
}

func (s *Server) cicdLogic(task models.Task, hk harborHook) {
	logId := s.initLog(task)

	dockerImg := fmt.Sprintf("%s/%s", s.conf.HarborAddress, strings.Split(hk.EventData.Resources[0].ResourceUrl, "/")[2])
	composeFile := fmt.Sprintf("%s/%s", task.Path, "docker-compose.yaml")
	file, err := os.ReadFile(composeFile)
	if err != nil {
		log.Error().Msgf("Cicd Error: 获取目录不存在 %s", composeFile)
		s.noticeLog(task.HarborKey, task.TaskName, fmt.Sprintf("cicd Error: 获取目录不存在 %s", composeFile))
		s.log(logId, false, fmt.Sprintf("Cicd Error: 获取目录不存在 %s", composeFile))
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
		s.log(logId, false, fmt.Sprintf("Cicd Error: 获取目录不存在 %s", composeFile))
		return
	}

	// ls コマンドを実行
	resp, err := utils.Exec(fmt.Sprintf("docker pull %s", dockerImg))
	if err != nil {
		log.Error().Msgf("Cicd Error: 执行错误 %s %s", err, resp)
		s.noticeLog(task.HarborKey, task.TaskName, fmt.Sprintf("cicd Error: 执行错误 %s %s", err, resp))
		s.log(logId, false, fmt.Sprintf("Cicd Error: 执行错误 %s %s", err, resp))
		return
	}

	resp, err = utils.Exec(task.Cmd)
	if err != nil {
		log.Error().Msgf("Cicd Error: 执行错误 %s %s", err, resp)
		s.noticeLog(task.HarborKey, task.TaskName, fmt.Sprintf("cicd Error: 执行错误 %s %s", err, resp))
		s.log(logId, false, fmt.Sprintf("Cicd Error: 执行错误 %s %s", err, resp))
		return
	}

	s.db.Model(&models.Task{}).Where("id = ?", task.ID).Update("last_run_time", time.Now().Unix())

	task.Heartbeat = strings.TrimSpace(task.Heartbeat)
	if task.Heartbeat != "" {
		resp, err := http.Get(task.Heartbeat)
		if err != nil {
			log.Error().Msgf("Cicd Error: 执行错误 Heartbeat 验证失败 %s", task.Heartbeat)
			s.noticeLog(task.HarborKey, task.TaskName, fmt.Sprintf("cicd Error: Heartbeat验证失败 %s", task.Heartbeat))
			s.log(logId, false, fmt.Sprintf("Cicd Error: 执行错误 Heartbeat 验证失败 %s", task.Heartbeat))
			return
		}
		if resp.StatusCode != 200 {
			log.Error().Msgf("Cicd Error: 执行错误 Heartbeat 验证失败 %s", task.Heartbeat)
			s.noticeLog(task.HarborKey, task.TaskName, fmt.Sprintf("cicd Error: Heartbeat验证失败 %s", task.Heartbeat))
			s.log(logId, false, fmt.Sprintf("Cicd Error: 执行错误 Heartbeat 验证失败 %s", task.Heartbeat))
			return
		}

		s.noticeLog(task.HarborKey, task.TaskName, fmt.Sprintf("success %s", task.Heartbeat))
		s.log(logId, true, fmt.Sprintf("success %s", task.Heartbeat))
		return
	}

	s.noticeLog(task.HarborKey, task.TaskName, "success")
	s.log(logId, true, fmt.Sprintf("success"))
}
