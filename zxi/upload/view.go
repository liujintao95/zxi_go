package upload

import (
	"fmt"
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
	if err != nil {
		fmt.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": core.ERR_BAD_REQ,
			"err_msg":  err.Error(),
		})
		return
	}

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
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": core.ERR_BAD_REQ,
			"err_msg":  err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"success":     true,
		"upload_info": handlers.GetUploadInfo(uploadId),
	})
}

func (v *View) ShowProgress(g *gin.Context) {
	uploadIdStr := g.Query("id")
	uploadId, err := strconv.Atoi(uploadIdStr)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": core.ERR_BAD_REQ,
			"err_msg":  err.Error(),
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"success":  true,
		"progress": handlers.GetProgress(uploadId),
	})
}

func (v *View) UploadFile(g *gin.Context) {
	uploadIdStr := g.PostForm("id")
	file, err := g.FormFile("file")
	uploadId, err := strconv.Atoi(uploadIdStr)
	if err != nil {
		fmt.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": core.ERR_BAD_REQ,
			"err_msg":  err.Error(),
		})
		return
	}
	fileMate := handlers.GetFileInfoByUploadId(uploadId)
	savePath := path.Join("files", fileMate.Hash+file.Filename)
	err = g.SaveUploadedFile(file, savePath)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": core.ERR_SAVE_FILE,
			"err_msg":  err.Error(),
		})
		return
	}
	handlers.UpdateFileComplete(uploadId, 1)

	g.JSON(http.StatusOK, gin.H{"success": true})
}

func (v *View) UploadBlock(g *gin.Context) {
	uploadIdStr := g.PostForm("upload_id")
	blockIdStr := g.PostForm("block_id")
	block, err := g.FormFile("block")
	uploadId, err := strconv.Atoi(uploadIdStr)
	blockId, err := strconv.Atoi(blockIdStr)
	if err != nil {
		fmt.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": core.ERR_BAD_REQ,
			"err_msg":  err.Error(),
		})
		return
	}
	fileMate := handlers.GetFileInfoByUploadId(uploadId)
	uploadBlockMate := handlers.GetUploadBlockInfo(blockId)

	err = handlers.CreateOrIgnorePath(path.Join("blocks", fileMate.Hash))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": core.ERR_SAVE_FILE,
			"err_msg":  err.Error(),
		})
		return
	}
	savePath := path.Join(
		"blocks", fileMate.Hash, strconv.Itoa(uploadBlockMate.Offset))
	err = g.SaveUploadedFile(block, savePath)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": core.ERR_SAVE_FILE,
			"err_msg":  err.Error(),
		})
		return
	}
	handlers.UpdateBlockComplete(blockId, 1)

	g.JSON(http.StatusOK, gin.H{"success": true})
}

func (v *View) PauseUpload(g *gin.Context)  {
	uploadIdStr := g.PostForm("upload_id")
	uploadId, err := strconv.Atoi(uploadIdStr)
	if err != nil {
		fmt.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": core.ERR_BAD_REQ,
			"err_msg":  err.Error(),
		})
		return
	}
	handlers.UpdateUploading(uploadId, 0)
	g.JSON(http.StatusOK, gin.H{
		"success": true,
		"upload_id": uploadId,
	})
}

func (v *View) StartUpload(g *gin.Context)  {
	uploadIdStr := g.PostForm("upload_id")
	uploadId, err := strconv.Atoi(uploadIdStr)
	if err != nil {
		fmt.Println(err)
		g.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": core.ERR_BAD_REQ,
			"err_msg":  err.Error(),
		})
		return
	}
	handlers.UpdateUploading(uploadId, 1)
	g.JSON(http.StatusOK, gin.H{
		"success": true,
		"upload_id": uploadId,
	})
}