package server

import (
	"github.com/dollarkillerx/harbor_easy_cicd/internal/models"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/request"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/resp"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"strconv"
)

func (s *Server) login(ctx *gin.Context) {
	var inputPayload = struct {
		Token string `json:"token"`
	}{}

	if err := ctx.ShouldBindJSON(&inputPayload); err != nil {
		resp.Resp(ctx, false, err.Error(), nil)
		return
	}

	if s.conf.AdminAuth.Token != inputPayload.Token {
		resp.Resp(ctx, false, "参数错误", nil)
		return
	}

	resp.Resp(ctx, true, "success", nil)
}

func (s *Server) tasks(ctx *gin.Context) {
	var tasks []models.Task
	if err := s.db.Model(&models.Task{}).Order("created_at desc").Find(&tasks).Error; err != nil {
		log.Info().Msgf("tasks error: %s", err)
		resp.Resp(ctx, false, "数据库异常", nil)
		return
	}

	resp.Resp(ctx, true, "success", tasks)
}

func (s *Server) task(ctx *gin.Context) {
	var input request.TaskPayload
	if err := ctx.ShouldBindJSON(&input); err != nil {
		resp.Resp(ctx, false, err.Error(), nil)
		return
	}

	if err := input.Validate(); err != nil {
		resp.Resp(ctx, false, err.Error(), nil)
		return
	}

	var err error

	switch input.Type {
	case request.TaskAdd:
		err = s.db.Model(&models.Task{}).Create(&models.Task{
			HarborKey:   input.HarborKey,
			TaskName:    input.TaskName,
			Path:        input.Path,
			Cmd:         input.Cmd,
			Heartbeat:   input.Heartbeat,
			Run:         true,
			LastRunTime: 0,
		}).Error
	case request.TaskDel:
		err = s.db.Model(&models.Task{}).Where("id = ?", input.TaskId).Delete(&models.Task{}).Error
	case request.TaskEdit:
		err = s.db.Model(&models.Task{}).Where("id = ?", input.TaskId).Updates(&models.Task{
			HarborKey: input.HarborKey,
			TaskName:  input.TaskName,
			Path:      input.Path,
			Cmd:       input.Cmd,
			Heartbeat: input.Heartbeat,
		}).Error
	case request.TaskStart:
		err = s.db.Model(&models.Task{}).Where("id = ?", input.TaskId).
			Update("run", true).Error
	case request.TaskStop:
		err = s.db.Model(&models.Task{}).Where("id = ?", input.TaskId).
			Update("run", false).Error
	}

	if err != nil {
		log.Error().Msgf("up err: %v", err)
		resp.Resp(ctx, false, err.Error(), nil)
		return
	}

	resp.Resp(ctx, true, "success", nil)
}

func (s *Server) logs(ctx *gin.Context) {
	id := ctx.Query("id")
	var logs []models.TaskLogs
	if id == "" {
		s.db.Model(&models.TaskLogs{}).Order("created_at desc").Find(&logs)
	} else {
		atoi, err := strconv.Atoi(id)
		if err == nil {
			s.db.Model(&models.TaskLogs{}).Where("task_id = ?", atoi).Order("created_at desc").Limit(20).Find(&logs)
		}
	}

	resp.Resp(ctx, true, "success", logs)
}
