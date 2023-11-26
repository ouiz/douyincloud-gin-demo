package service

import (
	"douyincloud-gin-demo/config"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos/enum"
)

// func GetPreSignUrlHandler(c *gin.Context, w http.ResponseWriter, r *http.Request) {
func GetPreSignUrlHandler(c *gin.Context) {
	// 初始化对象存储 client
	client, err := tos.NewClientV2(config.CfgCm.OsEndpoint, tos.WithRegion(config.CfgCm.OsRegion), tos.WithCredentials(tos.NewStaticCredentials(config.Cfg.OsAK, config.Cfg.OsSK)))
	if err != nil {
		// fmt.Fprint(w, "client init error")
		Failure(c, fmt.Errorf("%s", "client init error"))
		return
	}
	// objectName := r.URL.Query().Get("object_name")
	objectName := c.Query("object_name")
	if objectName == "" {
		// 如果为空，返回错误信息
		//   c.String(http.StatusBadRequest, "X-TT-OPENID 为空")
		Failure(c, fmt.Errorf("params objectName 为空"))
		return
	}
	// 调用Tos SDK 生成上传对象预签名
	url, err := client.PreSignedURL(&tos.PreSignedURLInput{
		HTTPMethod: enum.HttpMethodPut,
		Bucket:     config.CfgCm.OsBucketName,
		Key:        objectName,
	})
	if err != nil {
		// fmt.Fprint(w, "get pre sign url error")
		Failure(c, fmt.Errorf("%s", "get pre sign url error"))
		return
	}
	data := make(map[string]string)
	data[objectName] = url.SignedUrl
	// msg, err := json.Marshal(data)
	// if err != nil {
	// 	fmt.Fprint(w, "json marshal error")
	// 	return
	// }

	// w.Header().Set("content-type", "application/json")
	// w.Write(msg)
	SuccessData(c, data)
}
