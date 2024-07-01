package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/models"
	"github.com/gin-gonic/gin"
	"io"
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

/**
{
  "type": "PUSH_ARTIFACT",
  "occur_at": 1719318828,
  "operator": "admin",
  "event_data": {
    "resources": [
      {
        "digest": "sha256:746da633881a7c6c6f9a4d77225c6aa1728394b01fc2fb41bc0591b5239bbd98",
        "tag": "1.0.0",
        "resource_url": "192.168.78.129:8787/library/followme:1.0.0"
      }
    ],
    "repository": {
      "date_created": 1719316269,
      "name": "followme",
      "namespace": "library",
      "repo_full_name": "library/followme",
      "repo_type": "public"
    }
  }
}
*/

func (s *Server) webHookGithub(ctx *gin.Context) {
	all, err := io.ReadAll(ctx.Request.Body)
	if err == nil {
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

		}
	}

	var push models.GitHubPush
	if err := json.Unmarshal(all, &push); err == nil {

	}

	ctx.JSON(200, gin.H{})
}

func (s *Server) webHookGitee(ctx *gin.Context) {
	all, err := io.ReadAll(ctx.Request.Body)
	if err == nil {
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

		}
	}

	var push models.GitHubPush
	if err := json.Unmarshal(all, &push); err == nil {

	}

	ctx.JSON(200, gin.H{})
}

func (s *Server) webHookGitlib(ctx *gin.Context) {
	all, err := io.ReadAll(ctx.Request.Body)
	if err == nil {
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

		}
	}

	var push models.GitlabPush
	if err := json.Unmarshal(all, &push); err == nil {
		if tag.ObjectKind == "push" {

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
