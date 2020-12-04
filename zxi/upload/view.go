package upload

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

func (v *View) ShowUploads(c *gin.Context) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	userInter, _ := c.Get("userInfo")
	userMate := userInter.(models.UserInfo)
	page, err := strconv.Atoi(pageStr)
	size, err := strconv.Atoi(sizeStr)
	v.errCheck(c, err, errState.ErrBadReq)

	uploadList, count := v.handler.GetUploadTable(userMate.Id, page, size)
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"upload_list": uploadList,
		"count":       count,
	})
}

func (v *View) ShowUploadInfo(c *gin.Context) {
	uploadIdStr := c.Query("upload_id")
	uploadId, err := strconv.Atoi(uploadIdStr)
	v.errCheck(c, err, errState.ErrBadReq)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"upload_info": v.handler.GetUploadInfo(uploadId),
	})
}

func (v *View) ShowProgress(c *gin.Context) {
	uploadIdStr := c.Query("upload_id")
	uploadId, err := strconv.Atoi(uploadIdStr)
	v.errCheck(c, err, errState.ErrBadReq)

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"progress": v.handler.GetProgress(uploadId),
	})
}

func (v *View) UploadBuffer(c *gin.Context) {
	uploadIdStr := c.PostForm("upload_id")
	blockIdStr := c.PostForm("block_id")
	offsetStr := c.PostForm("offset")
	bufferStr := c.PostForm("buffer")
	uploadId, err := strconv.Atoi(uploadIdStr)
	blockId, err := strconv.Atoi(blockIdStr)
	v.errCheck(c, err, errState.ErrBadReq)

	fileMate := v.handler.GetFileInfoByUploadId(uploadId)
	err = v.handler.CreateOrIgnorePath(path.Join(FilePath, fileMate.Hash))
	v.errCheck(c, err, errState.ErrCreatePath)
	savePath := path.Join(FilePath, fileMate.Hash, offsetStr)
	err = v.handler.SaveFile([]byte(bufferStr), savePath)
	v.errCheck(c, err, errState.ErrSaveFile)
	v.handler.UpdateBlockComplete(blockId, 1)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (v *View) PauseUpload(c *gin.Context) {
	uploadIdStr := c.PostForm("upload_id")
	uploadId, err := strconv.Atoi(uploadIdStr)
	v.errCheck(c, err, errState.ErrBadReq)
	v.handler.UpdateUploading(uploadId, 0)
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"upload_id": uploadId,
	})
}

func (v *View) StartUpload(c *gin.Context) {
	uploadIdStr := c.PostForm("upload_id")
	uploadId, err := strconv.Atoi(uploadIdStr)
	v.errCheck(c, err, errState.ErrBadReq)

	v.handler.UpdateUploading(uploadId, 1)
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"upload_id": uploadId,
	})
}

func (v *View) CancelUpload(c *gin.Context) {
	uploadIdStr := c.PostForm("upload_id")
	uploadId, err := strconv.Atoi(uploadIdStr)
	v.errCheck(c, err, errState.ErrBadReq)

	v.handler.DeleteUpload(uploadId)
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"upload_id": uploadId,
	})
}
