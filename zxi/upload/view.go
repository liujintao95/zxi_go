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

func (v *View)ShowUploads(g *gin.Context) {
    userInter, _ := g.Get("userInfo")
    userMate := userInter.(models.UserInfo)
    uploadList := handlers.GetUploadList(userMate.Id)
    g.JSON(http.StatusOK, gin.H{
        "success":  true,
        "upload_list": uploadList,
    })
}

func (v *View)ShowUploadInfo(g *gin.Context) {
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
        "upload_info": handlers.GetUploadInfo(uploadId),
    })
}

func (v *View)UploadFile(g *gin.Context) {
    uploadIdStr := g.PostForm("id")
    f, err := g.FormFile("file_old")
    uploadId, err := strconv.Atoi(uploadIdStr)
    if err != nil {
        g.JSON(http.StatusBadRequest, gin.H{
            "success":  false,
            "err_code": core.ERR_BAD_REQ,
            "err_msg":  err.Error(),
        })
        return
    }
    err = g.SaveUploadedFile(f, path.Join("files", f.Filename))
    if err != nil {
        g.JSON(http.StatusBadRequest, gin.H{
            "success":  false,
            "err_code": core.ERR_SAVE_FILE,
            "err_msg":  err.Error(),
        })
        return
    }
    handlers.SetIsComplete(uploadId, 1)

    g.JSON(http.StatusOK, gin.H{"success":  true})
}