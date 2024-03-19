package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		files := form.File["file"]

		for _, file := range files {
			// 保存上传的文件到本地
			err := c.SaveUploadedFile(file, file.Filename)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			// 上传到七牛云的配置信息
			accessKey := ""
			secretKey := ""
			bucket := ""

			// 上传文件
			filePath := fmt.Sprintf("%s:%s", bucket, "stone/"+file.Filename)
			key := "stone/" + file.Filename // 文件上传到七牛云的路径
			// 生成上传凭证
			mac := qbox.NewMac(accessKey, secretKey)
			putPolicy := storage.PutPolicy{
				Scope: filePath, // 指定文件夹路径
			}
			upToken := putPolicy.UploadToken(mac)
			// 构建上传配置
			cfg := storage.Config{}
			// 空间对应的机房
			// 使用指定的区域
			cfg.Zone = &storage.Zone{
				SrcUpHosts: []string{"up-z1.qiniup.com"},
				RsHost:     "rs-z1.qiniu.com",
			}
			// 是否使用https域名
			cfg.UseHTTPS = false
			// 构建表单上传的对象
			formUploader := storage.NewFormUploader(&cfg)
			// 可选配置
			putExtra := storage.PutExtra{}

			err = formUploader.PutFile(c, nil, upToken, key, file.Filename, &putExtra)
			if err != nil {
				clean(file)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			clean(file)
			fmt.Println(fmt.Sprintf("外链地址: http://qiniu.yunxue521.top/stone/%s", file.Filename))
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	r.Run(":8080")
}

func clean(file *multipart.FileHeader) {
	// 上传成功后删除本地文件
	err := os.Remove(file.Filename)
	if err != nil {
		fmt.Println("删除本地文件失败:", err)
	} else {
		fmt.Println("删除本地文件成功:", file.Filename)
	}
}
