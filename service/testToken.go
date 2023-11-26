package service

import (
	"bytes"
	"douyincloud-gin-demo/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var token = ""

func getCltToken() string {
	res, err := getClientToken()
	if err != nil {
		fmt.Println("get clientToken err", err)
		return ""
	}
	return res.Data.AccessToken
}

// 定义获取token的接口的请求地址
const getTokenURL = "https://open.douyin.com/oauth/client_token/"

// 定义获取token的接口的请求参数结构体
type GetTokenRequest struct {
	ClientKey    string `json:"client_key"`    // 应用的client_key
	ClientSecret string `json:"client_secret"` // 应用的client_secret
	GrantType    string `json:"grant_type"`    // 固定值为"client_credential"
}

// 定义获取token的接口的返回结果结构体
type GetTokenResponse struct {
	Data    GetTokenData `json:"data"`    // 返回结果的对象
	Message string       `json:"message"` // 整体的返回结果的描述信息
	Code    int          `json:"code"`    // 整体的返回结果的状态码
}

// 定义返回结果的对象结构体
type GetTokenData struct {
	AccessToken string `json:"access_token"` // token的值
	ExpiresIn   int    `json:"expires_in"`   // token的有效期，单位为秒
	Scope       string `json:"scope"`        // token的权限范围，目前固定为"*"
}

// 定义一个函数，用于调用获取token的接口，传入应用的client_key和client_secret，返回token的值
func getClientToken() (GetTokenResponse, error) {
	// 构造获取token的接口的请求参数
	request := GetTokenRequest{
		ClientKey:    config.Cfg.AppId,  // 请替换为您的client_key
		ClientSecret: config.Cfg.Secret, // 请替换为您的client_secret
		GrantType:    "client_credential",
	}
	// 将请求参数转换为JSON格式
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return GetTokenResponse{}, err
	}
	// 创建一个HTTP客户端
	client := &http.Client{}
	// 创建一个HTTP请求
	req, err := http.NewRequest("POST", getTokenURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return GetTokenResponse{}, err
	}
	// 设置请求头的Content-Type为application/json
	req.Header.Set("Content-Type", "application/json")
	// 发送HTTP请求，并获取HTTP响应
	resp, err := client.Do(req)
	if err != nil {
		return GetTokenResponse{}, err
	}
	// 延迟关闭HTTP响应的Body
	defer resp.Body.Close()
	// 读取HTTP响应的Body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetTokenResponse{}, err
	}
	// 定义一个获取token的接口的返回结果变量
	var response GetTokenResponse
	// 将HTTP响应的Body转换为获取token的接口的返回结果
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return GetTokenResponse{}, err
	}
	// 返回获取token的接口的返回结果
	return response, nil
}
