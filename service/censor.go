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
const imgAntidirtURL = "https://open.douyin.com/api/apps/v1/censor/image"
const imgAntidirtURLCld = "http://open.douyin.com/api/apps/v1/censor/image"

//	type CensorImgReq struct {
//		appId string `json:"target"`
//		Name   string `json:"name"`
//	}

func CensorImg(ctx *gin.Context) {
	var req PictureDetectRequest
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
	ns := ctx.Query("ns")
	url := imgAntidirtURL
	if ns != "" {
		url = imgAntidirtURLCld
	}
	resp, err := pictureDetect(req, url, token)
	fmt.Printf("\nresp:%+v,\nerr:%v\n,ns:%v", resp, err, ns)
	if err != nil {
		fmt.Printf("1 detect err:%v", err)
		getCltToken(true)
		resp, err = pictureDetect(req, url, token)
		fmt.Printf("\nresp:%+v,\nerr:%v\n,ns:%v", resp, err, ns)
	}
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
func TestCI(ctx *gin.Context) {
	appId := ctx.Query("app_id")
	if appId == "" {
		Failure(ctx, fmt.Errorf("param invalid"))
		return
	}
	image := ctx.Query("image")
	if appId == "" {
		Failure(ctx, fmt.Errorf("param invalid"))
		return
	}
	token, err := getCltToken(false)
	fmt.Println("token:", token, err)
	if token == "" || err != nil {
		getCltToken(true)
		Failure(ctx, fmt.Errorf("get token error:%s,tk:%s", err, token))
		return
	}
	req := PictureDetectRequest{AppId: appId, Image: image}
	fmt.Printf("\n%+v\n", req)
	resp, err := pictureDetect(req, imgAntidirtURL, token)
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

type PictureDetectRequest struct {
	AppId     string `json:"app_id"`     // 应用的client_key
	Image     string `json:"image"`      // 应用的client_secret
	ImageData string `json:"image_data"` // 图片的URL或文件流
	// Scenes       []string `json:"scenes"`        // 图片检测场景
}
type PictureDetectResponse struct {
	ErrNo    int        `json:"err_no"`   // 整体的返回结果的状态码
	ErrMsg   string     `json:"err_msg"`  // 整体的返回结果的描述信息
	Predicts []Predicts `json:"predicts"` // 图片检测结果的数组
	ErrTips  string     `json:"err_tips"` // 失败 ：err_no 不为 0 时返回。
}

// 定义图片检测结果的结构体
type Predicts struct {
	ModelName string `json:"model_name"` // 图片检测场景的名称
	Hit       bool   `json:"hit"`        // 图片检测结果的标签
}

func checkPictureDetectPredicts(predicts []Predicts) (bool, []string) {
	// 定义一个变量，存储图片是否含有任何违法违规内容，初始值为false
	hasIllegalContent := false
	// 定义一个变量，存储图片含有违法违规内容的场景的名称，初始值为空切片
	illegalScenes := []string{}
	for _, predict := range predicts {
		if predict.Hit {
			hasIllegalContent = true
			illegalScenes = append(illegalScenes, predict.ModelName)
		}
	}
	return hasIllegalContent, illegalScenes
}

// func pictureDetect(app_id, imgURL string) (PictureDetectResponse, error) {
func pictureDetect(request PictureDetectRequest, apiURL, token string) (PictureDetectResponse, error) {
	// request := PictureDetectRequest{
	// 	AppId: app_id, Image: imgURL,
	// }
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return PictureDetectResponse{}, err
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
		return PictureDetectResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		fmt.Println("!!!!! token not empty", token)
		req.Header.Set("access-token", token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return PictureDetectResponse{}, err
	}
	defer resp.Body.Close()
	// 读取HTTP响应的Body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return PictureDetectResponse{}, err
	}

	var response PictureDetectResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return PictureDetectResponse{}, err
	}
	return response, nil
}
