package download

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
	"strconv"
	"zxi_go/core/errState"
	"zxi_go/core/logger"
	"zxi_go/core/models"
	"zxi_go/utils/customError"
)

type View struct {
	logger   *logrus.Logger
	handler  *Handler
	errCheck func(*gin.Context, error, int)
}

func NewView() *View {
	return &View{
		logger:   logger.LogInit(),
		handler:  NewHandler(),
		errCheck: customError.ErrorCheck,
	}
}

func (v *View) CreateDownload(c *gin.Context) {
	downloadPath := c.PostForm("download_path")
	fileIdStr := c.PostForm("file_id")
	userInter, _ := c.Get("userInfo")
	userMate := userInter.(models.UserInfo)
	fileId, err := strconv.Atoi(fileIdStr)
	v.errCheck(c, err, errState.ErrBadReq)

	lastId := v.handler.CreateOrIgnoreDownload(fileId, userMate.Id, downloadPath)
	v.handler.CreateOrIgnoreDownloadBlock(lastId)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (v *View) ShowDownload(c *gin.Context) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	userInter, _ := c.Get("userInfo")
	userMate := userInter.(models.UserInfo)
	page, err := strconv.Atoi(pageStr)
	size, err := strconv.Atoi(sizeStr)
	v.errCheck(c, err, errState.ErrBadReq)

	downloadList, count := v.handler.GetDownloadTable(userMate.Id, page, size)
	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"download_list": downloadList,
		"count":         count,
	})
}

func (v *View) ShowDownloadInfo(c *gin.Context) {
	downloadIdStr := c.Query("download_id")
	downloadId, err := strconv.Atoi(downloadIdStr)
	v.errCheck(c, err, errState.ErrBadReq)

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"download_info": v.handler.GetDownloadInfo(downloadId),
	})
}

func (v *View) ShowProgress(c *gin.Context) {
	downloadIdStr := c.Query("download_id")
	downloadId, err := strconv.Atoi(downloadIdStr)
	v.errCheck(c, err, errState.ErrBadReq)

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"progress": v.handler.GetProgress(downloadId),
	})
}

func (v *View) DownloadBuffer(c *gin.Context) {
	downloadIdStr := c.PostForm("download_id")
	blockIdStr := c.PostForm("block_id")
	offsetStr := c.PostForm("offset")
	downloadId, err := strconv.Atoi(downloadIdStr)
	blockId, err := strconv.Atoi(blockIdStr)
	v.errCheck(c, err, errState.ErrBadReq)

	fileMate := v.handler.GetFileInfoByDownloadId(downloadId)
	filePath := path.Join(FilePath, fileMate.Hash, offsetStr)
	bufferStr, err := v.handler.readFile(filePath)
	v.errCheck(c, err, errState.ErrSaveFile)
	v.handler.UpdateBlockComplete(blockId, 1)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"buffer":  bufferStr,
	})
}

func (v *View) PauseDownload(c *gin.Context) {
	downloadIdStr := c.PostForm("download_id")
	downloadId, err := strconv.Atoi(downloadIdStr)
	v.errCheck(c, err, errState.ErrBadReq)
	v.handler.UpdateDownloading(downloadId, 0)
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"download_id": downloadId,
	})
}

func (v *View) StartDownload(c *gin.Context) {
	downloadIdStr := c.PostForm("download_id")
	downloadId, err := strconv.Atoi(downloadIdStr)
	v.errCheck(c, err, errState.ErrBadReq)

	v.handler.UpdateDownloading(downloadId, 1)
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"download_id": downloadId,
	})
}

func (v *View) CancelDownload(c *gin.Context) {
	downloadIdStr := c.PostForm("download_id")
	downloadId, err := strconv.Atoi(downloadIdStr)
	v.errCheck(c, err, errState.ErrBadReq)

	v.handler.DeleteDownload(downloadId)
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"upload_id": downloadId,
	})
}
