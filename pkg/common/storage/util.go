package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/dongle"
	"github.com/speps/go-hashids/v2"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	downloadRoutePath = "download"
	uploadRoutePath   = "upload"
	hashIdSalt        = "eta local storage presigned"
	timeFmt           = "20060102T150405Z"
	SEP               = "/"
)

/*
NewObjectKey 文件名转换

	ex: tbl_xxx/a/b/c/d.jpg -> tbl_xxx/1704704407/a/b/c/d.jpg
*/
func NewObjectKey(fileName, rootFolder string, appendTime bool) (objKey, filePath string) {
	st := strings.Split(fileName, SEP)
	if appendTime {
		st[1] = strings.Join([]string{strconv.FormatInt(time.Now().Unix(), 10), SEP, st[1]}, "")
	}
	if rootFolder != "" {
		st[0] = strings.Join([]string{rootFolder, st[0]}, SEP)
	}
	objKey = strings.Join(st, SEP)
	filePath = strings.Join(st[1:], SEP) // 表名后（不含表名）的路径
	return
}

func S3Client(platform string, accessKey, secretKey, endpoint string, region string, https bool) (*s3.Client, error) {
	var (
		endpointResolver = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: endpoint}, nil
		})
		pathStyle = true
	)
	switch platform {
	case "aliyun_oss":
		pathStyle = false
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(endpointResolver),
		config.WithClientLogMode(aws.LogRetries|aws.LogRequestWithBody|aws.LogResponseWithBody|aws.LogRequestEventMessage|aws.LogResponseEventMessage|aws.LogSigning),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = pathStyle
		o.EndpointOptions.DisableHTTPS = !https
	})

	return s3Client, nil
}

func buildLocalPreSignedParam(host, objKey string, lifetimeSecs int64, appId string, userId string) url.Values {
	hd := hashids.NewData()
	hd.Salt = hashIdSalt
	hd.MinLength = 10
	h, _ := hashids.NewWithData(hd)

	a, _ := convertor.ToInt(appId)
	u, _ := convertor.ToInt(userId)
	now := time.Now().UTC()

	uId, _ := h.EncodeInt64([]int64{u})
	algoKey, _ := h.EncodeInt64([]int64{now.Unix(), a, u})
	//fmt.Println(algoKey)
	//d, _ := h.DecodeWithError(e)
	//fmt.Println(d)

	date := now.Format(timeFmt)
	expired := strconv.FormatInt(lifetimeSecs, 10)
	cre := fmt.Sprintf("%s/%s", appId, uId)

	toSignStr := fmt.Sprintf("%s%s%s%s%s", host, objKey, date, expired, cre)
	logger.Debugf("toSignStr: %s, key: %s", toSignStr, algoKey)

	params := url.Values{} //拼接query参数
	params.Add(ParamKeyAlgorithm, DefaultAlgorithm)
	params.Add(ParamKeyCredential, cre)
	params.Add(ParamKeyDate, date)
	params.Add(ParamKeyExpires, expired)
	params.Add(ParamKeySignature, dongle.Encrypt.FromString(toSignStr).ByHmacSha256(algoKey).ToHexString())
	return params
}

func VerifyLocalPreSignedParam(c *gin.Context, filepath string) (aId, uId string, err error) {
	var (
		algoStr    = c.Query(ParamKeyAlgorithm)
		creStr     = c.Query(ParamKeyCredential)
		dateStr    = c.Query(ParamKeyDate)
		expiresStr = c.Query(ParamKeyExpires)
		signStr    = c.Query(ParamKeySignature)
	)

	date, err := time.Parse(timeFmt, dateStr)
	if err != nil {
		logger.Errorf("%s 解析失败. %v", ParamKeyDate, err)
		return "", "", errors.New("X-Delta-Date 解析失败")
	}
	expires, err := convertor.ToInt(expiresStr)
	if err != nil {
		logger.Errorf("%s 解析失败. %v", ParamKeyExpires, err)
		return "", "", errors.New("X-Delta-Expires 解析失败")
	}
	if date.Add(time.Duration(expires) * time.Second).Before(time.Now()) {
		logger.Warnf("签名已过期")
		return "", "", errors.New("签名已过期")
	}
	cres := strings.Split(creStr, SEP)
	if len(cres) != 2 {
		logger.Warnf("%s 格式错误", ParamKeyCredential)
		return "", "", errors.New(fmt.Sprintf("%s 格式错误", ParamKeyCredential))
	}
	appIdStr := cres[0]
	appId, _ := convertor.ToInt(appIdStr)
	userIdStr := cres[1]

	hd := hashids.NewData()
	hd.Salt = hashIdSalt
	hd.MinLength = 10
	h, _ := hashids.NewWithData(hd)

	userIds, _ := h.DecodeInt64WithError(userIdStr)
	if len(userIds) != 1 {
		logger.Warnf("%s 格式错误", ParamKeyCredential)
		return "", "", errors.New(fmt.Sprintf("%s 格式错误", ParamKeyCredential))
	}
	userId := userIds[0]
	logger.Debugf("storage 本地存储签名. appId:%s, userId:%s", appIdStr, convertor.ToString(userId))

	host := c.Request.Host
	objKey := strings.TrimPrefix(filepath, SEP)
	toSignStr := fmt.Sprintf("%s%s%s%s%s", host, objKey, dateStr, expiresStr, creStr)
	algoKey, _ := h.EncodeInt64([]int64{date.Unix(), appId, userId})

	switch algoStr {
	case DefaultAlgorithm:
		logger.Debugf("toSignStr: %s, key: %s", toSignStr, algoKey)
		newSignStr := dongle.Encrypt.FromString(toSignStr).ByHmacSha256(algoKey).ToHexString()
		if newSignStr != signStr {
			logger.Warnf("签名不匹配")
			return "", "", errors.New("签名不匹配")
		}
		return appIdStr, convertor.ToString(userId), nil
	default:
		logger.Errorf("不支持的签名算法.%s", algoStr)
		return "", "", errors.New("不支持的签名算法")
	}
}

func AppStorageFolder(appId string) string {
	lastSlashIndex := strings.LastIndex(StoreRootFolder, SEP)
	var rootFolder string
	if lastSlashIndex != -1 {
		rootFolder = StoreRootFolder[:lastSlashIndex]
	} else {
		rootFolder = "/data/delta_storage"
	}
	return fmt.Sprintf(`%s/%s`, rootFolder, appId)
}
