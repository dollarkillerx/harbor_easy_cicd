package server

import (
	"github.com/gin-gonic/gin"
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
