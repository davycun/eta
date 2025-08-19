package storage

import (
	"github.com/davycun/eta/pkg/common/global"
	"path"
)

func InitModule() {
	handler := Controller{}
	group := global.GetGin().Group("/storage")
	group.GET("/public/pre_sign/download", handler.PublicPreSignDownload)             // 公共目录下载预签名，无需登录
	group.POST("/public/pre_sign/download", handler.PublicPreSignDownload)            // 公共目录下载预签名，无需登录
	group.GET("/pre_sign/download", handler.PreSignDownload)                          // 下载预签名
	group.POST("/pre_sign/download", handler.PreSignDownload)                         // 下载预签名
	group.POST("/pre_sign/upload", handler.PreSignUpload)                             // 上传预签名
	group.GET(path.Join("/", downloadRoutePath, "/*filepath"), handler.DownloadFile)  // 下载文件
	group.HEAD(path.Join("/", downloadRoutePath, "/*filepath"), handler.DownloadFile) // 下载文件
	group.PUT(path.Join("/", uploadRoutePath, "/*filepath"), handler.UploadFile)      // 上传文件
	group.POST(path.Join("/", uploadRoutePath, "/*filepath"), handler.UploadFile)     // 上传文件
	group.GET("/list_objects", handler.ListObjects)                                   // 对象列表
	group.POST("/delete_objects", handler.DeleteObjects)                              // 删除对象
}
