package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dollarkillerx/harbor_easy_cicd/internal/models"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *Server) webHook(ctx *gin.Context) {
	// PUSH_ARTIFACT

	var hook harborHook
	if err := ctx.ShouldBindJSON(&hook); err != nil {
		ctx.JSON(400, gin.H{
			"error": err,
		})
		return
	}

	go s.cicd(hook)

	ctx.JSON(200, gin.H{
		"message": "success",
	})
}

type harborHook struct {
	Type      string `json:"type"`
	OccurAt   int    `json:"occur_at"`
	Operator  string `json:"operator"`
	EventData struct {
		Resources []struct {
			Digest      string `json:"digest"`
			Tag         string `json:"tag"`
			ResourceUrl string `json:"resource_url"`
		} `json:"resources"`
		Repository struct {
			DateCreated  int    `json:"date_created"`
			Name         string `json:"name"`
			Namespace    string `json:"namespace"`
			RepoFullName string `json:"repo_full_name"`
			RepoType     string `json:"repo_type"`
		} `json:"repository"`
	} `json:"event_data"`
}

func (s *Server) webHookGithub(ctx *gin.Context) {
	all, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "400",
		})
		return
	}

	if !validateSignature(all, ctx.GetHeader("X-Hub-Signature-256")) {
		ctx.JSON(401, gin.H{
			"error": "auth error",
		})
		return
	}
	// 判断 tag push?
	// tag 判断
	var tag models.GithubTag
	if err := json.Unmarshal(all, &tag); err == nil {
		if tag.RefType == "tag" {
			var tagTasks []models.GitTask
			s.db.Model(&models.GitTask{}).Where("matching = ?", models.MathTag).Where("git_type = ?", models.Github).Find(&tagTasks)
			for i, v := range tagTasks {
				idx := i
				// 项目判断
				if strings.Contains(tag.Repository.FullName, v.Repository) {
					if strings.Contains(tag.Ref, v.Tag) {
						// develop
						go s.gitCicd(tagTasks[idx])
					}
				}
			}
		}
	}

	var push models.GitHubPush
	if err := json.Unmarshal(all, &push); err == nil {
		var tagTasks []models.GitTask
		s.db.Model(&models.GitTask{}).Where("matching = ?", models.MathTag).Where("git_type = ?", models.Github).Find(&tagTasks)
		for i, v := range tagTasks {
			idx := i
			// 项目判断
			if strings.Contains(push.Repository.FullName, v.Repository) {
				// 分支判断
				if strings.Contains(push.Ref, v.Branch) {
					// comment判断
					if strings.Contains(push.HeadCommit.Message, v.Comment) {
						// develop
						go s.gitCicd(tagTasks[idx])
					}
				}
			}
		}
	}

	ctx.JSON(200, gin.H{})
}

func (s *Server) webHookGitee(ctx *gin.Context) {
	all, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "400",
		})
		return
	}

	if !validateSignature(all, ctx.GetHeader("X-Hub-Signature-256")) {
		ctx.JSON(401, gin.H{
			"error": "auth error",
		})
		return
	}
	// 判断 tag push?
	// tag 判断
	var tag models.GiteaTag
	if err := json.Unmarshal(all, &tag); err == nil {
		if tag.Action == "published" {
			var tagTasks []models.GitTask
			s.db.Model(&models.GitTask{}).Where("matching = ?", models.MathTag).Where("git_type = ?", models.Github).Find(&tagTasks)
			for i, v := range tagTasks {
				// 项目判断
				if strings.Contains(tag.Repository.FullName, v.Repository) {
					if strings.Contains(tag.Release.TagName, v.Tag) {
						// develop
						idx := i
						go s.gitCicd(tagTasks[idx])
					}
				}
			}
		}
	}

	var push models.GitHubPush
	if err := json.Unmarshal(all, &push); err == nil {
		var tagTasks []models.GitTask
		s.db.Model(&models.GitTask{}).Where("matching = ?", models.MathTag).Where("git_type = ?", models.Github).Find(&tagTasks)
		for i, v := range tagTasks {
			idx := i
			// 项目判断
			if strings.Contains(push.Repository.FullName, v.Repository) {
				// 分支判断
				if strings.Contains(push.Ref, v.Branch) {
					// comment判断
					if strings.Contains(push.HeadCommit.Message, v.Comment) {
						// develop
						go s.gitCicd(tagTasks[idx])
					}
				}
			}
		}
	}

	ctx.JSON(200, gin.H{})
}

func (s *Server) webHookGitlib(ctx *gin.Context) {
	all, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "400",
		})
		return
	}

	// 判断 tag push?
	// tag 判断
	var tag models.GitlabTag
	if err := json.Unmarshal(all, &tag); err == nil {
		if tag.ObjectKind == "tag_push" {
			var tagTasks []models.GitTask
			s.db.Model(&models.GitTask{}).Where("matching = ?", models.MathTag).Where("git_type = ?", models.Github).Find(&tagTasks)
			for i, v := range tagTasks {
				idx := i
				// 项目判断
				if strings.Contains(tag.Repository.Url, v.Repository) {
					// 分支判断
					if strings.Contains(tag.Ref, v.Tag) {
						// comment判断
						go s.gitCicd(tagTasks[idx])
					}
				}
			}
		}
	}

	var push models.GitlabPush
	if err := json.Unmarshal(all, &push); err == nil {
		if tag.ObjectKind == "push" {
			var tagTasks []models.GitTask
			s.db.Model(&models.GitTask{}).Where("matching = ?", models.MathTag).Where("git_type = ?", models.Github).Find(&tagTasks)
			for i, v := range tagTasks {
				idx := i
				// 项目判断
				if strings.Contains(push.Repository.Url, v.Repository) {
					// 分支判断
					if strings.Contains(push.Ref, v.Branch) {
						// comment判断
						if len(push.Commits) > 0 {
							if strings.Contains(push.Commits[0].Message, v.Comment) {
								// develop
								go s.gitCicd(tagTasks[idx])
							}
						}
					}
				}
			}
		}
	}

	ctx.JSON(200, gin.H{})
}

// 验证 GitHub 发送的请求
func validateSignature(payload []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte("WUSmhf2tKuVhm"))
	mac.Write(payload)
	expectedMAC := mac.Sum(nil)
	expectedSignature := "sha256=" + hex.EncodeToString(expectedMAC)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func (s *Server) gitCicd(task models.GitTask) {
	logId := s.initGithubLog(task)

	s.noticeLog(fmt.Sprintf("%s %s", task.GitType, task.Repository), task.TaskName, "获取到任务")

	err := os.Chdir(task.Path)
	if err != nil {
		log.Error().Msgf("Git Cicd Error: 获取目录不存在 %s", task.Path)
		s.noticeLog(fmt.Sprintf("%s %s", task.GitType, task.Repository), task.TaskName, fmt.Sprintf("Git Cicd Error: 获取目录不存在 %s", task.Path))
		s.log(logId, false, fmt.Sprintf("Cicd Error: 获取目录不存在 %s", task.Path))
		return
	}

	// ls コマンドを実行
	resp, err := utils.Exec(fmt.Sprintf("git pull"))
	if err != nil {
		log.Error().Msgf("Git Cicd Error: 执行错误 %s %s", err, resp)
		s.noticeLog(fmt.Sprintf("%s %s", task.GitType, task.Repository), task.TaskName, fmt.Sprintf("Git Cicd Error: 执行错误 %s %s", err, resp))
		s.log(logId, false, fmt.Sprintf("Git Cicd Error: 执行错误 %s %s", err, resp))
		return
	}

	resp, err = utils.Exec(task.Cmd)
	if err != nil {
		log.Error().Msgf("Git Cicd Error: 执行错误 %s %s", err, resp)
		s.noticeLog(fmt.Sprintf("%s %s", task.GitType, task.Repository), task.TaskName, fmt.Sprintf("Git Cicd Error: 执行错误 %s %s", err, resp))
		s.log(logId, false, fmt.Sprintf("Git Cicd Error: 执行错误 %s %s", err, resp))
		return
	}

	task.Heartbeat = strings.TrimSpace(task.Heartbeat)
	if task.Heartbeat != "" {
		resp, err := http.Get(task.Heartbeat)
		if err != nil {
			log.Error().Msgf("Cicd Error: 执行错误 Heartbeat 验证失败 %s", task.Heartbeat)
			s.noticeLog(fmt.Sprintf("%s %s", task.GitType, task.Repository), task.TaskName, fmt.Sprintf("Cicd Error: 执行错误 Heartbeat 验证失败 %s", task.Heartbeat))
			s.log(logId, false, fmt.Sprintf("Cicd Error: 执行错误 Heartbeat 验证失败 %s", task.Heartbeat))
			return
		}
		if resp.StatusCode != 200 {
			log.Error().Msgf("Cicd Error: 执行错误 Heartbeat 验证失败 %s", task.Heartbeat)
			s.noticeLog(fmt.Sprintf("%s %s", task.GitType, task.Repository), task.TaskName, fmt.Sprintf("Cicd Error: 执行错误 Heartbeat 验证失败 %s", task.Heartbeat))
			s.log(logId, false, fmt.Sprintf("Cicd Error: 执行错误 Heartbeat 验证失败 %s", task.Heartbeat))
			return
		}

		s.noticeLog(fmt.Sprintf("%s %s", task.GitType, task.Repository), task.TaskName, fmt.Sprintf("success %s", task.Heartbeat))
		s.log(logId, true, fmt.Sprintf("success %s", task.Heartbeat))
		return
	}

	s.noticeLog(fmt.Sprintf("%s %s", task.GitType, task.Repository), task.TaskName, "success")
	s.log(logId, true, "success")
}

func (s *Server) initGithubLog(task models.GitTask) uint {
	var log = models.TaskLogs{
		TaskId:   task.ID,
		TaskName: task.TaskName,
		Message:  LogJson{"获取到任务"}.ToJson(),
	}
	s.db.Model(&models.TaskLogs{}).Create(&log)
	return log.ID
}
