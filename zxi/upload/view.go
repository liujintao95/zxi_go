package upload

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strconv"
	"zxi_go/core"
	"zxi_go/zxi/models"
)

type View struct {
}

var view = View{}

func (v *View) ShowUploads(g *gin.Context) {
	userInter, _ := g.Get("userInfo")
	userMate := userInter.(models.UserInfo)
	pageStr := g.Query("page")
	sizeStr := g.Query("size")
	page, err := strconv.Atoi(pageStr)
	size, err := strconv.Atoi(sizeStr)
	core.CustomError(g, err, core.ErrBadReq)

	uploadList, count := handlers.GetUploadList(userMate.Id, page, size)
	g.JSON(http.StatusOK, gin.H{
		"success":     true,
		"upload_list": uploadList,
		"count":       count,
	})
}

func (v *View) ShowUploadInfo(g *gin.Context) {
	uploadIdStr := g.Query("id")
	uploadId, err := strconv.Atoi(uploadIdStr)
	core.CustomError(g, err, core.ErrBadReq)

	g.JSON(http.StatusOK, gin.H{
		"success":     true,
		"upload_info": handlers.GetUploadInfo(uploadId),
	})
}

func (v *View) ShowProgress(g *gin.Context) {
	uploadIdStr := g.Query("id")
	uploadId, err := strconv.Atoi(uploadIdStr)
	core.CustomError(g, err, core.ErrBadReq)

	g.JSON(http.StatusOK, gin.H{
		"success":  true,
		"progress": handlers.GetProgress(uploadId),
	})
}

func (v *View) UploadFile(g *gin.Context) {
	uploadIdStr := g.PostForm("upload_id")
	fileStr := g.PostForm("file")
	uploadId, err := strconv.Atoi(uploadIdStr)
	core.CustomError(g, err, core.ErrBadReq)

	fileMate := handlers.GetFileInfoByUploadId(uploadId)
	savePath := path.Join("files", fileMate.Hash)
	err = handlers.SaveFile([]byte(fileStr), savePath)
	core.CustomError(g, err, core.ErrSaveFile)
	handlers.UpdateFileComplete(uploadId, 1)

	g.JSON(http.StatusOK, gin.H{"success": true})
}

func (v *View) UploadBlock(g *gin.Context) {
	uploadIdStr := g.PostForm("upload_id")
	blockIdStr := g.PostForm("block_id")
	blockStr := g.PostForm("block")
	uploadId, err := strconv.Atoi(uploadIdStr)
	blockId, err := strconv.Atoi(blockIdStr)
	core.CustomError(g, err, core.ErrBadReq)

	fileMate := handlers.GetFileInfoByUploadId(uploadId)
	uploadBlockMate := handlers.GetUploadBlockInfo(blockId)
	err = handlers.CreateOrIgnorePath(path.Join("blocks", fileMate.Hash))
	core.CustomError(g, err, core.ErrCreatePath)
	savePath := path.Join(
		"blocks", fileMate.Hash, strconv.Itoa(uploadBlockMate.Offset))
	err = handlers.SaveFile([]byte(blockStr), savePath)
	core.CustomError(g, err, core.ErrSaveFile)
	handlers.UpdateBlockComplete(blockId, 1)

	g.JSON(http.StatusOK, gin.H{"success": true})
}

func (v *View) PauseUpload(g *gin.Context)  {
	uploadIdStr := g.PostForm("upload_id")
	uploadId, err := strconv.Atoi(uploadIdStr)
	core.CustomError(g, err, core.ErrBadReq)
	handlers.UpdateUploading(uploadId, 0)
	g.JSON(http.StatusOK, gin.H{
		"success": true,
		"upload_id": uploadId,
	})
}

func (v *View) StartUpload(g *gin.Context)  {
	uploadIdStr := g.PostForm("upload_id")
	uploadId, err := strconv.Atoi(uploadIdStr)
	core.CustomError(g, err, core.ErrBadReq)

	handlers.UpdateUploading(uploadId, 1)
	g.JSON(http.StatusOK, gin.H{
		"success": true,
		"upload_id": uploadId,
	})
}

func (v *View) CancelUpload(g *gin.Context)  {
	uploadIdStr := g.PostForm("upload_id")
	uploadId, err := strconv.Atoi(uploadIdStr)
	core.CustomError(g, err, core.ErrBadReq)

	handlers.DeleteUpload(uploadId)
	g.JSON(http.StatusOK, gin.H{
		"success": true,
		"upload_id": uploadId,
	})
}