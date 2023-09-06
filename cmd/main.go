package main

import (
	"github.com/gin-gonic/gin"

	"github.com/kannan112/video-streaming/pkg/streaming"
	"github.com/kannan112/video-streaming/pkg/uploader"
)

func main() {
	r := gin.Default()

	// route for uploading video
	r.POST("/upload", uploader.Upload)

	// route for streaming video using hls
	r.GET("/play/:video_id/:playlist", streaming.Streamin)

	r.Run(":7001")
}
