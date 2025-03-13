package storage

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/storage"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	downloadRoutePath = "download"
	uploadRoutePath   = "upload"
	publicContextKey  = "eta_storage_public_visit"
	anonymityUserId   = "10000"
	anonymityAppId    = "10001"
	anonymityUser     = user.User{BaseEntity: entity.BaseEntity{ID: anonymityUserId}}
	anonymityApp      = app.App{BaseEntity: entity.BaseEntity{ID: anonymityAppId}}
)

type Controller struct {
	controller.DefaultController
}

func (handler *Controller) PublicPreSignDownload(c *gin.Context) {
	var (
		err   error
		ct    = ctx.GetContext(c)
		bc, _ = setting.GetStorageConfig(ct.GetAppGorm())
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			if bc.PublicFolder == "" {
				return errs.NewClientError("未配置公共目录，请检查 setting 配置")
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			//a, err = app.LoadDefaultApp()
			//if err != nil {
			//	return errs.NewClientError("default app not found or no storage config")
			//}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			ct, err = user.NewContext(anonymityUser, anonymityApp, "")
			if err != nil {
				return errs.NewClientError("初始化 context 失败")
			}
			ctx.BindGinContext(ct, c)
			ct.Set(publicContextKey, true)
			return nil
		}).Err
	if err != nil {
		handler.ProcessResult(c, nil, err)
		return
	}

	handler.preSign(c, ct, storage.Download)
}

func (handler *Controller) PreSignDownload(c *gin.Context) {
	handler.preSign(c, nil, storage.Download)
}

func (handler *Controller) PreSignUpload(c *gin.Context) {
	handler.preSign(c, nil, storage.Upload)
}

func (handler *Controller) DownloadFile(c *gin.Context) {
	var (
		filePath = c.Param("filepath")
		err      error
		aId, uId string
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			aId, uId, err = storage.VerifyLocalPreSignedParam(c, filePath)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			return initAppUserContext(c, aId, uId)
		}).Err

	if err != nil {
		handler.ProcessResult(c, nil, err)
		return
	}

	if filePath == "" || filePath == storage.SEP {
		controller.Fail(c, 404, "文件路径不能为空", nil)
		return
	}
	bc, _ := setting.GetStorageConfig(ctx.GetContext(c).GetAppGorm())
	if bc.Platform != "local" || bc.RootFolder == "" {
		handler.ProcessResult(c, nil, errors.New("app not found or storage local config error"))
		return
	}

	filePath = strings.TrimPrefix(filePath, storage.SEP)
	localRoot := strings.TrimSuffix(bc.Storage.RootFolder, storage.SEP)

	osFileName := strings.Join([]string{localRoot, filePath}, storage.SEP)
	fileInfo, err := os.Stat(osFileName)
	if os.IsNotExist(err) {
		controller.Fail(c, 404, "文件不存在", nil)
		return
	}
	if err != nil { // 未知错误
		controller.Fail(c, 404, "获取文件信息出错", nil)
		return
	}
	if !fileInfo.Mode().IsRegular() { //必须是常规文件
		controller.Fail(c, 404, "只能下载文件", nil)
		return
	}

	routePath := strings.Join([]string{"/storage", downloadRoutePath}, storage.SEP)
	fileServer := http.StripPrefix(routePath, http.FileServer(gin.Dir(localRoot, false)))
	fileServer.ServeHTTP(c.Writer, c.Request)
}

func (handler *Controller) UploadFile(c *gin.Context) {
	var (
		filePath = c.Param("filepath")
		err      error
		aId, uId string
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			aId, uId, err = storage.VerifyLocalPreSignedParam(c, filePath)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			return initAppUserContext(c, aId, uId)
		}).Err

	if err != nil {
		handler.ProcessResult(c, nil, err)
		return
	}

	//app, exists := user.GetContextApp(ctx.GetContext(c))
	//app.Storage.Platform, app.Storage.RootFolder = "local", "/tmp/delta_test" // todo delete
	bc, _ := setting.GetStorageConfig(ctx.GetContext(c).GetAppGorm())
	if bc.Platform != "local" || bc.RootFolder == "" {
		handler.ProcessResult(c, nil, errors.New("app not found or storage local config error"))
		return
	}

	rawData, err := handler.getRawData(c)
	if err != nil {
		handler.ProcessResult(c, nil, err)
		return
	}
	err = handler.validateUploadFile(c, filePath, &rawData)
	if err != nil {
		controller.Fail(c, 400, err.Error(), nil)
		return
	}

	svc := storage.NewStorageSvc(ctx.GetContext(c), bc.Storage)
	err = svc.PutObject(filePath, rawData)
	if err != nil {
		handler.ProcessResult(c, nil, err)
		return
	}

	controller.Success(c, nil)
}

func (handler *Controller) ListObjects(c *gin.Context) {
	var (
		ct     = ctx.GetContext(c)
		cfg, _ = setting.GetStorageConfig(ct.GetAppGorm())
	)
	//app.Storage.Platform, app.Storage.RootFolder = "local", "/tmp/delta_test" // todo delete
	if cfg.Platform == "" {
		handler.ProcessResult(c, nil, errors.New("app not found or storage local config error"))
		return
	}

	prefix := c.Query("prefix")
	svc := storage.NewStorageSvc(ctx.GetContext(c), cfg.Storage)
	res, err := svc.ListObjects(prefix)
	handler.ProcessResult(c, res, err)
}
func (handler *Controller) DeleteObjects(c *gin.Context) {
	var (
		param storage.DeleteObjectsParam
		ct    = ctx.GetContext(c)
		appDb = ct.GetAppGorm()
		//bc          = setting.GetBackendConfig(ct.GetAppGorm())
		bc, _ = setting.GetStorageConfig(appDb)
	)
	if bc.Platform == "" {
		handler.ProcessResult(c, nil, errors.New("app not found or storage config err"))
		return
	}
	err := controller.BindBody(c, &param)
	if err != nil {
		handler.ProcessResult(c, nil, err)
		return
	}

	svc := storage.NewStorageSvc(ctx.GetContext(c), bc.Storage)
	res, err := svc.DeleteObjects(&param)
	handler.ProcessResult(c, res, err)
}

func (handler *Controller) preSign(c *gin.Context, ct *ctx.Context, preSignedType string) {
	if ct == nil {
		ct = ctx.GetContext(c)
	}
	var (
		err    error
		param  []storage.PreSignedParam
		result []storage.PreSignedResult
		//ap     *app.LocatedApp
		svc   *storage.StoreSvc
		bc, _ = setting.GetStorageConfig(ct.GetAppGorm())
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			if bc.Platform == "" {
				return errs.NewClientError("ap not found or no storage config")
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			svc = storage.NewStorageSvc(ct, bc.Storage)
			// 下载单个文件，支持GET
			var p storage.PreSignedParam
			if ct.GetBool(publicContextKey) && preSignedType == storage.Download {
				er := controller.BindQuery(c, &p)
				if er == nil && p.FileName != "" {
					param = append(param, p)
				}
			}
			if len(param) <= 0 {
				return controller.BindBody(c, &param)
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			if ct.GetBool(publicContextKey) {
				pf := strutil.Trim(bc.PublicFolder, storage.SEP)
				for i, p := range param {
					param[i].FileName = strings.Join([]string{pf, strutil.Trim(p.FileName, storage.SEP)}, storage.SEP)
				}
			}
			return handler.validatePreSign(c, &param)
		}).
		Call(func(cl *caller.Caller) error {
			result = make([]storage.PreSignedResult, len(param))
			return svc.BatchPreSignedUrl(preSignedType, &param, &result)
		}).Err

	if preSignedType == storage.Download && len(param) == 1 && param[0].Redirect {
		c.Redirect(http.StatusFound, result[0].PreSignedUrl)
	} else {
		handler.ProcessResult(c, &result, err)
	}
}

func (handler *Controller) validatePreSign(c *gin.Context, param *[]storage.PreSignedParam) error {
	//app, exists := user.GetContextApp(ctx.GetContext(c))
	var (
		bc, _ = setting.GetStorageConfig(ctx.GetContext(c).GetAppGorm())
	)
	if bc.Platform == "" {
		return errors.New("no storage config")
	}
	maxBatchSize := 100
	if len(*param) > maxBatchSize {
		return errors.New(fmt.Sprintf(`批量大小不能超过 %d`, maxBatchSize))
	}

	var fileTypes []string
	fileTypeMatch := true
	if bc.AllowFiletype != "" {
		fileTypes = strings.Split(bc.AllowFiletype, ",")
		fileTypeMatch = false
	}
	for i := range *param {
		p := &(*param)[i]
		if p.FileName == "" {
			return errors.New("文件名不能为空")
		}
		if strings.HasPrefix(p.FileName, storage.SEP) {
			p.FileName = p.FileName[1:]
		}
		if p.Expires == 0 {
			p.Expires = 900
		}
		if len(strings.Split(p.FileName, storage.SEP)) < 2 {
			return errors.New("文件名需要包含表名")
		}
		// 判断 p.FileName 是否以 fileTypes 里的元素为后缀
		for _, fileType := range fileTypes {
			if strings.HasSuffix(p.FileName, fileType) {
				fileTypeMatch = true
				break
			}
		}
		if !fileTypeMatch {
			return errors.New("文件类型不合法")
		}
	}

	return nil
}

func (handler *Controller) validateUploadFile(c *gin.Context, filePath string, rawData *[]byte) (err error) {

	// 上传文件路径
	if filePath == "" {
		return errors.New("文件路径不能为空")
	}
	filePath = strings.TrimPrefix(filePath, storage.SEP)
	filePath = strings.TrimSuffix(filePath, storage.SEP)

	bc, _ := setting.GetStorageConfig(ctx.GetContext(c).GetAppGorm())

	// 文件类型
	var fileTypes []string
	fileTypeMatch := true
	if bc.AllowFiletype != "" {
		fileTypes = strings.Split(bc.AllowFiletype, ",")
		fileTypeMatch = false
	}
	for _, fileType := range fileTypes {
		if strings.HasSuffix(filePath, fileType) {
			fileTypeMatch = true
			break
		}
	}
	if !fileTypeMatch {
		return errors.New("文件类型不合法")
	}
	// headers.Size 获取文件大小
	if len(*rawData) > 1024*1024*200 {
		return errors.New("文件太大了")
	}
	return nil
}

func (handler *Controller) getRawData(c *gin.Context) ([]byte, error) {
	body := c.Request.Body
	return io.ReadAll(body)
}

func initAppUserContext(c *gin.Context, appId, userId string) (err error) {
	var (
		us user.User
		ap app.App
		ct = ctx.GetContext(c)
	)
	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			if userId == anonymityUserId {
				us = anonymityUser
			} else {
				us, err = user.LoadUserById(global.GetLocalGorm(), userId)
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			if appId == anonymityAppId {
				ap = anonymityApp
			} else {
				ap, err = app.LoadAppById(global.GetLocalGorm(), appId)
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			if us.ID != "" && ap.ID != "" {
				ct, err = user.NewContext(us, ap, "")
			} else {
				return errs.NewServerError("用户或应用不存在")
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			ctx.BindGinContext(ct, c)
			ct.Set(publicContextKey, true)
			return nil
		}).Err

	return err
}
