package streaming

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Streamin(ctx *gin.Context) {
	videoId := ctx.Param("video_id")
	PlayList := ctx.Param("playlist")

	PlayListDataChan := make(chan []byte)
	ErrorChan := make(chan error)

	go func() {
		byteVideo, err := readPlayList(videoId, PlayList)
		if err != nil {
			ErrorChan <- err
			return
		}
		PlayListDataChan <- byteVideo
	}()

	select {
	case playListByte := <-PlayListDataChan:
		ctx.Header("Content-Type", "application/vnd.apple.mpegurl")
		ctx.Header("Content-Disposition", "inline")

		ctx.Writer.Write(playListByte)
	case err := <-ErrorChan: // This case checks if there is an error available on the errChan channel
		// Handle the error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to read file from server",
			"error":   err.Error(),
		})
	case <-time.After(5 * time.Second):
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "request time out",
		})
	}

}

func readPlayList(videoId string, playList string) ([]byte, error) {
	playListPath := fmt.Sprintf("storage/%s/%s", videoId, playList)
	PlayListData, err := ioutil.ReadFile(playListPath)
	if err != nil {
		return nil, err
	}
	return PlayListData, err

}
