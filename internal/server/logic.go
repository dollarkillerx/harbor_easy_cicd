package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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

/*
*
tag commit
Gitee:
Github:
*/
func (s *Server) webHookGit(ctx *gin.Context) {
	all, err := io.ReadAll(ctx.Request.Body)
	if err == nil {
		//log.Info().Msgf("%s", all)
	}

	log.Info().Msgf("ok ? %v", validateSignature(all, ctx.GetHeader("X-Hub-Signature-256")))

	marshal, err := json.Marshal(ctx.Request.Header)
	if err == nil {
		fmt.Println(string(marshal))
	}
}

// 验证 GitHub 发送的请求
func validateSignature(payload []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte("WUSmhf2tKuVhm"))
	mac.Write(payload)
	expectedMAC := mac.Sum(nil)
	expectedSignature := "sha256=" + hex.EncodeToString(expectedMAC)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
