/*
Copyright (year) Bytedance Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"douyincloud-gin-demo/config"
	"douyincloud-gin-demo/service"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitCfg()
	// component.InitComponents()

	r := gin.Default()

	r.GET("/api/hello", service.Hello)
	r.GET("/api/test", service.Test)
	r.GET("/api/openid", service.Openid)
	r.POST("/api/set_name", service.SetName)

	r.GET("/api/pre_sign_url", service.GetPreSignUrlHandler)
	r.POST("/api/censor_img", service.CensorImg)
	r.GET("/api/test_ci", service.TestCI)

	r.POST("/api/censor_img2", service.CensorImg2)
	r.POST("/api/censor_img2ns", service.CensorImg2Nos)
	r.GET("/api/get_ci", service.GetCI2)

	r.Run(":8000")
}
