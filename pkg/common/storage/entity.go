package storage

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
)

var (
	Download                = "download"
	Upload                  = "upload"
	StoreRootFolder         = "/data/eta_storage/default"
	DefaultStorage  Storage = Storage{
		Platform: "local",
		//Endpoint:    endpoint,
		Https:       true,
		RootFolder:  StoreRootFolder,
		ForwardPath: "/eta",
	}
)

type Storage struct {
	dorm.JsonType
	Platform      string `json:"platform,omitempty" binding:"oneof=minio aliyun_oss local ''"` //平台：minio，aliyun_oss
	AccessKey     string `json:"accessKey,omitempty"`                                          //accessKey
	SecretKey     string `json:"secretKey,omitempty"`                                          //secretKey
	Endpoint      string `json:"endpoint,omitempty"`                                           //endpoint,阿里云的为https://oss-{region}.aliyuncs.com
	Region        string `json:"region,omitempty"`                                             // region
	Bucket        string `json:"bucket,omitempty"`                                             // bucket
	Https         bool   `json:"https,omitempty"`                                              // 是否启用https
	RootFolder    string `json:"root_folder,omitempty"`                                        // 上传的文件在存储的根目录下, 前面不要加"/"
	AllowFiletype string `json:"allow_filetype,omitempty"`                                     // 允许上传的文件类型，多个用逗号隔开
	ForwardPath   string `json:"forward_path,omitempty"`                                       // 存储路径转发
}

type PreSignedParam struct {
	FileName          string            `form:"file_name" json:"file_name,omitempty" binding:"required"`
	Headers           map[string]string `form:"headers" json:"headers,omitempty"`
	Expires           int64             `form:"expires" json:"expires,omitempty" default:"900"`
	DisableAppendTime bool              `form:"disable_append_time" json:"disable_append_time,omitempty"`
	Redirect          bool              `form:"redirect" json:"redirect,omitempty"` // 当“下载”&&只有一个预授权时，是否重定向
}

type PreSignedResult struct {
	Platform     string `json:"platform,omitempty"`                          // 存储平台
	PreSignedUrl string `json:"pre_signed_url,omitempty" binding:"required"` // 预签名url
	FileName     string `json:"file_name,omitempty" binding:"required"`      // 请求参数里的file_name回填
	FilePath     string `json:"file_path,omitempty"`                         // 文件路径，表名后的路径
}

type ListObjectsResult struct {
	Name         string           `json:"name"`                     // 名称
	IsDir        bool             `json:"is_dir"`                   // 是否目录
	Size         *ctype.Integer   `json:"size,omitempty"`           // 大小
	LastModified *ctype.LocalTime `json:"last_modified,omitempty" ` // 最后修改时间
}

type DeleteObjectsParam struct {
	Keys []string `json:"keys,omitempty"` // 对象key
}

type DeleteObjectsResult struct {
	Key     string `json:"key"`     // 对象key
	Deleted bool   `json:"deleted"` // 是否删除
}
