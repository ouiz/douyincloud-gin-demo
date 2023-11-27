package service

import (
	"bytes"
	"douyincloud-gin-demo/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// const textAntidirtURL = "http://developer.toutiao.com/api/v2/tags/text/antidirt"
const imgCensor2s = "https://developer.toutiao.com/api/apps/censor/image"
const imgCensor2Cld = "http://developer.toutiao.com/api/apps/censor/image"

//	type CensorImgReq struct {
//		appId string `json:"target"`
//		Name   string `json:"name"`
//	}

func CensorImg2(ctx *gin.Context) {
	var req PictureDetectRequest2
	err := ctx.Bind(&req)
	if err != nil {
		Failure(ctx, err)
		return
	}
	token, err := getCltToken(false)
	fmt.Println("token:", token, err)
	if token == "" || err != nil {
		Failure(ctx, fmt.Errorf("get token error:%s,tk:%s", err, token))
		return
	}
	resp, err := pictureDetect2(req, imgCensor2s, token)
	fmt.Printf("\nresp:%+v,\nerr:%v\n", resp, err)
	if err != nil {
		getCltToken(true)
		Failure(ctx, err)
		return
	}
	if resp.ErrNo != 0 {
		Failure(ctx, fmt.Errorf("%s_%s", resp.ErrMsg, resp.ErrTips))
		return
	}
	// checkPictureDetectPredicts()
	SuccessData(ctx, resp.Predicts)
}
func CensorImg2Nos(ctx *gin.Context) {
	var req PictureDetectRequest2
	err := ctx.Bind(&req)
	if err != nil {
		Failure(ctx, err)
		return
	}
	resp, err := pictureDetect2(req, imgCensor2Cld, "")
	fmt.Printf("\n2nos resp:%+v,\nerr:%v\n", resp, err)
	if err != nil {
		Failure(ctx, err)
		return
	}
	if resp.ErrNo != 0 {
		Failure(ctx, fmt.Errorf("%s_%s", resp.ErrMsg, resp.ErrTips))
		return
	}
	// checkPictureDetectPredicts()
	SuccessData(ctx, resp.Predicts)
}

func GetCI2(ctx *gin.Context) {
	appId := ctx.Query("app_id")
	if appId == "" {
		Failure(ctx, fmt.Errorf("param appid invalid"))
		return
	}
	image := ctx.Query("image")
	if image == "" {
		Failure(ctx, fmt.Errorf("param image invalid"))
		return
	}
	ft := ctx.Query("fresh_t")
	freshT := false
	if ft != "" {
		freshT = true
	}
	token, err := getCltToken(freshT)

	fmt.Println("token:", token, err, freshT)
	if token == "" || err != nil {
		Failure(ctx, fmt.Errorf("get token error:%s,tk:%s", err, token))
		return
	}
	req := PictureDetectRequest2{AppId: appId, AccessToken: token, Image: image}
	fmt.Printf("\n%+v\n", req)

	nos := ctx.Query("nos")
	url := imgCensor2s
	if nos != "" {
		url = imgCensor2Cld
	}
	resp, err := pictureDetect2(req, url, token)
	fmt.Printf("\nresp:%+v,\nerr:%v\n", resp, err)
	if err != nil {
		Failure(ctx, err)
		return
	}
	if resp.ErrNo != 0 {
		Failure(ctx, fmt.Errorf("%s_%s", resp.ErrMsg, resp.ErrTips))
		return
	}
	// checkPictureDetectPredicts()
	SuccessData(ctx, resp.Predicts)
}

type PictureDetectRequest2 struct {
	AppId       string `json:"app_id"` // 应用的client_key
	AccessToken string `json:"access_token"`
	Image       string `json:"image"`      // 应用的client_secret
	ImageData   string `json:"image_data"` // 图片的URL或文件流
	// Scenes       []string `json:"scenes"`        // 图片检测场景
}
type PictureDetectResponse2 struct {
	ErrNo     int        `json:"error"`    // 整体的返回结果的状态码
	ErrMsg    string     `json:"message"`  // 整体的返回结果的描述信息
	Predicts  []Predicts `json:"predicts"` // 图片检测结果的数组
	ErrTips   string     `json:"err_tips"` // 失败 ：err_no 不为 0 时返回。
	ModelName string     `json:"model_name"`
	Hit       bool       `json:"hit"`
}

// "error": 0,
//   "message": "image censor success",
// // 定义图片检测结果的结构体
// type Predicts struct {
// 	ModelName string `json:"model_name"` // 图片检测场景的名称
// 	Hit       bool   `json:"hit"`        // 图片检测结果的标签
// }

// func pictureDetect(app_id, imgURL string) (PictureDetectResponse, error) {
func pictureDetect2(request PictureDetectRequest2, apiURL, token string) (PictureDetectResponse2, error) {
	// request := PictureDetectRequest{
	// 	AppId: app_id, Image: imgURL,
	// }
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return PictureDetectResponse2{}, err
	}
	// 创建一个HTTP客户端
	client := &http.Client{}
	// apiURL := imgAntidirtURLCld
	// if config.Cfg.IsLocal {
	// 	apiURL = imgAntidirtURL
	// }
	fmt.Println("url,isLocal", apiURL, config.Cfg.IsLocal)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return PictureDetectResponse2{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		fmt.Println("!!!!! token not empty", token)
		req.Header.Set("access-token", token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return PictureDetectResponse2{}, err
	}
	defer resp.Body.Close()
	// 读取HTTP响应的Body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return PictureDetectResponse2{}, err
	}

	var response PictureDetectResponse2
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return PictureDetectResponse2{}, err
	}
	return response, nil
}
