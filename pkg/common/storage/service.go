package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/slice"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type StoreSvc struct {
	//service.DefaultService
	storage Storage
	c       *ctx.Context
}

func NewStorageSvc(c *ctx.Context, s Storage) *StoreSvc {
	svc := &StoreSvc{}
	svc.c = c
	svc.storage = s
	return svc
}

func (s *StoreSvc) BatchPreSignedUrl(preSignedType string, param *[]PreSignedParam, result *[]PreSignedResult) error {
	for i := 0; i < len(*param); i++ {
		err := s.PreSignedUrl(preSignedType, &(*param)[i], &(*result)[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *StoreSvc) PreSignedUrl(preSignedType string, param *PreSignedParam, result *PreSignedResult) error {
	var (
		fileName    = strings.TrimPrefix(param.FileName, SEP)
		contentType = param.Headers["Content-Type"]
		storage     = s.storage
		platform    = storage.Platform
		accessKey   = storage.AccessKey
		secretKey   = storage.SecretKey
		endpoint    = storage.Endpoint
		region      = storage.Region
		bucket      = storage.Bucket
		https       = storage.Https
		forwardPath = storage.ForwardPath
		rootFolder  = strings.TrimPrefix(storage.RootFolder, SEP)
	)
	result.Platform = platform
	result.FileName = fileName

	if forwardPath != "" {
		// 主要是看minio经过nginx转发时用的是不是https
		endpoint = fmt.Sprintf("%s://%s", utils.HttpScheme(https), utils.RequestHost(s.c.GetGinContext()))
		//endpoint = fmt.Sprintf("%s://%s", endpointSchema, "zzlydj.wgly.hangzhou.gov.cn")
		if platform == "local" {
			endpoint = fmt.Sprintf("%s%s", endpoint, forwardPath)
		}
	}
	var presigner S3Presigner
	if platform != "local" {
		s3client, err := S3Client(platform, accessKey, secretKey, endpoint, region, https)
		if err != nil {
			return err
		}
		presignClient := s3.NewPresignClient(s3client)
		presigner = S3Presigner{PresignClient: presignClient}
	}

	switch preSignedType {
	case Download: // 下载
		if platform == "local" {
			objKey, _ := NewObjectKey(fileName, "", false)
			result.FilePath = objKey
			rawUrl := strings.Join([]string{endpoint, "storage", downloadRoutePath, objKey}, SEP)
			uri, err := url.Parse(rawUrl)
			if err != nil {
				return err
			}
			host := s.c.GetGinContext().Request.Host
			uri.RawQuery = buildLocalPreSignedParam(
				host, objKey, param.Expires,
				s.c.GetContextAppId(),
				s.c.GetContextUserId(),
			).Encode()
			result.PreSignedUrl = uri.String()
			return nil
		}

		objKey, _ := NewObjectKey(fileName, rootFolder, false)
		result.FilePath = objKey
		presignedGetRequest, _err := presigner.GetObject(bucket, objKey, param.Expires)
		if _err != nil {
			return _err
		}
		result.PreSignedUrl = presignedGetRequest.URL

	case Upload: // 上传
		if platform == "local" {
			objKey, _ := NewObjectKey(fileName, "", !param.DisableAppendTime)
			result.FilePath = objKey
			rawUrl := strings.Join([]string{endpoint, "storage", uploadRoutePath, objKey}, SEP)
			uri, err := url.Parse(rawUrl)
			if err != nil {
				return err
			}
			host := s.c.GetGinContext().Request.Host
			uri.RawQuery = buildLocalPreSignedParam(
				host, objKey, param.Expires,
				s.c.GetContextAppId(),
				s.c.GetContextUserId(),
			).Encode()
			result.PreSignedUrl = uri.String()
			return nil
		}
		objKey, _ := NewObjectKey(fileName, rootFolder, !param.DisableAppendTime)
		result.FilePath = objKey
		presignedPutRequest, _err := presigner.PutObject(bucket, objKey, param.Expires, contentType)
		if _err != nil {
			return _err
		}
		result.PreSignedUrl = presignedPutRequest.URL
	}
	if forwardPath != "" {
		u, _err := url.Parse(result.PreSignedUrl)
		if _err != nil {
			return _err
		}
		u.Scheme = utils.Scheme(s.c.GetGinContext()) // 要换成请求进来时的scheme
		u.Path = fmt.Sprintf("%s%s", forwardPath, u.Path)
		result.PreSignedUrl = u.String()
	}
	return nil
}

func (s *StoreSvc) PutObject(objKey string, rawData []byte) (err error) {
	var (
		storage   = s.storage
		platform  = storage.Platform
		accessKey = storage.AccessKey
		secretKey = storage.SecretKey
		endpoint  = storage.Endpoint
		region    = storage.Region
		bucket    = storage.Bucket
		https     = storage.Https
	)

	switch s.storage.Platform {
	case "local":
		localRoot := strings.TrimSuffix(s.storage.RootFolder, SEP)
		dst := filepath.Join(localRoot, objKey)

		if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
			logger.Errorf("文件目录写入失败。error: %v", err)
			return errors.New("文件目录写入失败")
		}
		err = os.WriteFile(dst, rawData, 0750)
		if err != nil {
			logger.Errorf("文件写入失败。error: %v", err)
			return errors.New("文件写入失败")
		}
	case "minio", "aliyun_oss":
		s3client, err := S3Client(platform, accessKey, secretKey, endpoint, region, https)
		if err != nil {
			return err
		}
		output, err := s3client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &bucket,
			Key:    &objKey,
			Body:   bytes.NewReader(rawData),
		})
		logger.Infof("PutObject output: %v", output)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown platform")
	}

	return
}

func (s *StoreSvc) GetObject(objKey string) (rawData []byte, err error) {
	var (
		storage   = s.storage
		platform  = storage.Platform
		accessKey = storage.AccessKey
		secretKey = storage.SecretKey
		endpoint  = storage.Endpoint
		region    = storage.Region
		bucket    = storage.Bucket
		https     = storage.Https
	)

	switch s.storage.Platform {
	case "local":
		localRoot := strings.TrimSuffix(s.storage.RootFolder, SEP)
		dst := filepath.Join(localRoot, objKey)
		return os.ReadFile(dst)
	case "minio", "aliyun_oss":
		s3client, err := S3Client(platform, accessKey, secretKey, endpoint, region, https)
		if err != nil {
			return rawData, err
		}
		output, err := s3client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &objKey,
		})
		if err != nil {
			return rawData, err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logger.Warnf("Close body error: %v", err)
			}
		}(output.Body)
		return io.ReadAll(output.Body)
	default:
		return nil, errors.New("unknown platform")
	}

}

func (s *StoreSvc) ListObjects(prefix string) (*[]ListObjectsResult, error) {
	result := make([]ListObjectsResult, 0)
	if !s.c.GetContextIsManager() {
		return &result, errors.New("只有管理员才能列举文件")
	}
	var (
		storage    = s.storage
		platform   = storage.Platform
		accessKey  = storage.AccessKey
		secretKey  = storage.SecretKey
		endpoint   = storage.Endpoint
		region     = storage.Region
		bucket     = storage.Bucket
		https      = storage.Https
		rootFolder = storage.RootFolder
		delimiter  = SEP
	)
	//var ( // minio
	//	platform   = "minio"
	//	accessKey  = "mdc"                    //storage.AccessKey
	//	secretKey  = "mdc@1234"               //storage.SecretKey
	//	endpoint   = "http://127.0.0.1:29000" //storage.Endpoint
	//	region     = "us-east-1"              //storage.Region
	//	bucket     = "mdc"                    //storage.Bucket
	//	https      = false                    //storage.Https
	//	rootFolder = "eta/dev_storage"
	//	delimiter  = SEP
	//)
	//var ( // local
	//	platform  = "local"
	//	accessKey = "mdc_ak"                //storage.AccessKey
	//	secretKey = "mdc_sk123"             //storage.SecretKey
	//	endpoint  = "http://127.0.0.1:9001" //storage.Endpoint
	//	region    = "us-east-1"             //storage.Region
	//	bucket    = "mdc"                   //storage.Bucket
	//	https     = false                   //storage.Https
	//  rootFolder = "/data/tmp/"
	//	delimiter = SEP
	//)

	switch platform {
	case StorePlatformLocal:
		dst := filepath.Join(rootFolder, prefix)
		files, err := utils.ListFiles(dst)
		if err != nil {
			return nil, err
		}
		result = append(result, slice.Map(files, func(index int, item os.DirEntry) ListObjectsResult {
			if item.IsDir() {
				return ListObjectsResult{Name: item.Name(), IsDir: true}
			} else {
				var (
					size *ctype.Integer
					mt   *ctype.LocalTime
				)
				info, err1 := item.Info()
				if err1 == nil {
					size = ctype.NewIntPrt(info.Size())
					mt = ctype.NewLocalTimePrt(info.ModTime())
				}
				return ListObjectsResult{Name: item.Name(), IsDir: false, Size: size, LastModified: mt}
			}
		})...)

		return &result, nil
	case StorePlatformMinio, StorePlatformAliYunOss:
		rootFolder = strings.TrimPrefix(rootFolder, SEP)
		s3client, err1 := S3Client(platform, accessKey, secretKey, endpoint, region, https)
		if err1 != nil {
			return &result, err1
		}
		prefix = filepath.Join(rootFolder, prefix) + SEP
		input := &s3.ListObjectsV2Input{
			Bucket:            &bucket,
			ContinuationToken: nil,
			Delimiter:         &delimiter,
			Prefix:            &prefix,
			RequestPayer:      "",
			StartAfter:        nil,
		}
		for {
			resp, err2 := s3client.ListObjectsV2(context.TODO(), input)
			if err2 != nil {
				return &result, err2
			}
			result = append(result, slice.Map(resp.CommonPrefixes, func(index int, item types.CommonPrefix) ListObjectsResult {
				n := strings.TrimSuffix(strings.TrimPrefix(*item.Prefix, prefix), SEP)
				return ListObjectsResult{Name: n, IsDir: true}
			})...)

			result = append(result, slice.Map(resp.Contents, func(index int, item types.Object) ListObjectsResult {
				n := strings.TrimPrefix(*item.Key, prefix)
				return ListObjectsResult{Name: n, IsDir: false, Size: ctype.NewIntPrt(*item.Size), LastModified: ctype.NewLocalTimePrt(*item.LastModified)}
			})...)

			if !*resp.IsTruncated {
				break
			}
			input.ContinuationToken = resp.NextContinuationToken
		}

		return &result, nil
	default:
		return &result, errors.New("unknown platform")
	}

}

func (s *StoreSvc) DeleteObjects(param *DeleteObjectsParam) (*[]DeleteObjectsResult, error) {
	result := make([]DeleteObjectsResult, 0)
	if !s.c.GetContextIsManager() {
		return &result, errors.New("只有管理员才能删除文件")
	}
	if param == nil {
		return &result, nil
	}
	var (
		storage    = s.storage
		platform   = storage.Platform
		accessKey  = storage.AccessKey
		secretKey  = storage.SecretKey
		endpoint   = storage.Endpoint
		region     = storage.Region
		bucket     = storage.Bucket
		https      = storage.Https
		rootFolder = strings.TrimPrefix(storage.RootFolder, SEP)
	)

	switch platform {
	case "local":
		var (
			toDeleteFiles = slice.Map(param.Keys, func(index int, item string) string { return filepath.Join(s.storage.RootFolder, item) })
			errs          = make([]error, 0)
		)
		for _, f := range toDeleteFiles {
			key := strings.TrimPrefix(f, s.storage.RootFolder)
			err := fileutil.RemoveFile(f)
			if err != nil {
				result = append(result, DeleteObjectsResult{Key: key, Deleted: false})
				errs = append(errs, err)
			} else {
				result = append(result, DeleteObjectsResult{Key: key, Deleted: true})
			}
		}

		var err error
		if len(errs) > 0 {
			err = errors.New(fmt.Sprintf("%s", strings.Join(slice.Map(errs, func(index int, item error) string {
				return item.Error()
			}), "\n")))
		}
		return &result, err

	case "minio", "aliyun_oss":
		s3client, err1 := S3Client(platform, accessKey, secretKey, endpoint, region, https)
		if err1 != nil {
			return &result, err1
		}
		prefix := rootFolder + SEP
		input := &s3.DeleteObjectsInput{
			Bucket: &bucket,
			Delete: &types.Delete{
				Objects: slice.Map(param.Keys, func(index int, item string) types.ObjectIdentifier {
					objKey, _ := NewObjectKey(item, rootFolder, false)
					return types.ObjectIdentifier{Key: &objKey}
				})},
			BypassGovernanceRetention: nil,
			ChecksumAlgorithm:         "",
			ExpectedBucketOwner:       nil,
			MFA:                       nil,
			RequestPayer:              "",
		}
		resp, err2 := s3client.DeleteObjects(context.TODO(), input)
		if err2 != nil {
			return &result, err2
		}
		for _, obj := range resp.Deleted {
			result = append(result, DeleteObjectsResult{Key: strings.TrimPrefix(*obj.Key, prefix), Deleted: true})
		}
		if len(resp.Errors) > 0 {
			return &result, errors.New(fmt.Sprintf("%s", strings.Join(slice.Map(resp.Errors, func(index int, item types.Error) string {
				return *item.Message
			}), "\n")))
		}

		return &result, nil
	default:
		return &result, errors.New("unknown platform")
	}

}
