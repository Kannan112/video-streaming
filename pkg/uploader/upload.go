package uploader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	storageLocation = "storage"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch video file from request",
			"error":   err.Error(),
		})
		return
	}

	FileUuid := uuid.New()
	FileName := FileUuid.String()
	FolderPath := storageLocation + "/" + FileName
	FilePath := storageLocation + "/" + FileName + "/" + "video.mp4"

	// creating a directory
	allPaths := []string{storageLocation, FolderPath}
	for _, dirNames := range allPaths {
		if err := os.MkdirAll(dirNames, 0755); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "failed to fetch video file from request",
				"error":   err.Error(),
			})
			return
		}
	}

	newFile, err := os.Create(FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to crate file to copy video file",
			"error":   err.Error(),
		})
		return
	}
	defer newFile.Close()
	fmt.Println("File is created successfully.")

	src, err := file.Open()

	//copy uploaded file to new file
	_, err = io.Copy(newFile, src)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Video Uploaded Successfully",
		"video_id": FileUuid,
	})
	go func() {
		err = CreatePlaylistAndSegments(FilePath, FolderPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to create segments and playlist",
				"error":   err.Error(),
			})
			return
		}
		fmt.Println("exited without error")
	}()
}

func CreatePlaylistAndSegments(filePath string, folderPath string) error {
	//defer wg.Done()
	//TODO : calculate segment duration depending on video length
	segmentDuration := 3
	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-i", filePath,
		"-profile:v", "baseline", // baseline profile is compatible with most devices
		"-level", "3.0",
		"-start_number", "0", // start number segments from 0
		"-hls_time", strconv.Itoa(segmentDuration), //duration of each segment in second
		"-hls_list_size", "0", // keep all segments in the playlist
		"-f", "hls",
		fmt.Sprintf("%s/playlist.m3u8", folderPath),
	)
	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create HLS: %v \nOutput: %s ", err, string(output))
	}
	return nil
}
