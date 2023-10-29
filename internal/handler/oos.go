package handler

import (
	"colatiger/api/response"
	v1 "colatiger/api/v1/req"
	"colatiger/api/v1/res"
	"colatiger/config"
	"github.com/gin-gonic/gin"
	"github.com/jassue/gin-wire/app/pkg/request"
	progress "github.com/markity/minio-progress"
	"github.com/minio/minio-go/v7"
	"github.com/sony/sonyflake"
	"go.uber.org/zap"
	"path"
	"strconv"
	"time"
)

type OssHandler struct {
	log    *zap.Logger
	client *minio.Client
	sf     *sonyflake.Sonyflake
	conf   *config.Configuration
}

func NewOssHandler(log *zap.Logger, client *minio.Client, sf *sonyflake.Sonyflake, conf *config.Configuration) *OssHandler {
	return &OssHandler{
		log:    log,
		client: client,
		sf:     sf,
		conf:   conf,
	}
}

func (o *OssHandler) Upload(ctx *gin.Context) {

	var form v1.UploadReq

	if err := ctx.ShouldBind(&form); err != nil {
		response.FailByErr(ctx, request.GetError(form, err))
		return
	}

	var sucList []res.OssSuccess
	var failList []res.OssFail

	for _, obj := range form.Files {
		file, err := obj.Open()
		defer file.Close()
		if err != nil {
			o.log.Error("文件读取失败")
			failList = append(failList, res.OssFail{
				OriFileName: obj.Filename,
				Status:      "fail",
				Message:     err.Error(),
			})
			continue
		}
		id, _ := o.sf.NextID()
		name := strconv.FormatUint(id, 10) + path.Ext(obj.Filename)
		objectName := time.Now().Format("20060102") + "/" + name

		// 创建进度条对象, 需要在参数中输入文件的大小
		progressBar := progress.NewUploadProgress(obj.Size)
		result, err := o.client.PutObject(ctx, o.conf.Oss.BucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: "application/octet-stream", Progress: progressBar})

		if err != nil {
			failList = append(failList, res.OssFail{
				OriFileName: obj.Filename,
				Status:      "fail",
				Message:     err.Error(),
			})
		} else {
			sucList = append(sucList, res.OssSuccess{
				Bucket:      result.Bucket,
				Key:         strconv.FormatUint(id, 10),
				PreviewUrl:  result.Location,
				OriFileName: obj.Filename,
				FileSize:    result.Size,
			})
		}

	}

	response.Success(ctx, res.OssReq{
		Total:      len(form.Files),
		FailNum:    len(failList),
		SuccessNum: len(sucList),
		OssFail:    failList,
		OssSuccess: sucList,
	})
}
