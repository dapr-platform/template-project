package service

import (
	"context"
	"encoding/json"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
	"github.com/dapr-platform/common"
	"github.com/dapr/go-sdk/client"
	"github.com/pkg/errors"
	"wxgateway-service/config"
	"wxgateway-service/model"
)

var cacheApps map[string]*miniProgram.MiniProgram

func init() {
	cacheApps = make(map[string]*miniProgram.MiniProgram)
}
func PhoneCodeLogin(ctx context.Context, appName string, code string) (tokenInfo *model.TokenInfo, err error) {
	app, err := getMiniProgramApp(appName)
	if err != nil {
		return nil, errors.Wrap(err, "get mini program app error")
	}
	res, err := app.PhoneNumber.GetUserPhoneNumber(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, "get user phone number error")
	}
	if res.PhoneInfo == nil {
		return nil, errors.New("phone info is empty")
	}
	phone := res.PhoneInfo.PhoneNumber

	req := make(map[string]string)
	req["field"] = "mobile"
	req["value"] = phone
	req["client_id"] = config.CLIENT_ID
	req["client_secret"] = config.CLIENT_SECRET
	reqData, _ := json.Marshal(req)
	content := &client.DataContent{
		Data:        reqData,
		ContentType: "application/x-www-form-urlencoded",
	}

	data, err := common.GetDaprClient().InvokeMethodWithContent(ctx, "authz-service", "/oauth/token-by-field", "POST", content)
	if err != nil {
		return nil, errors.Wrap(err, "invoke method error")
	}
	tokenInfo = &model.TokenInfo{}
	err = json.Unmarshal(data, tokenInfo)
	return
}

func getMiniProgramApp(appName string) (app *miniProgram.MiniProgram, err error) {
	if app, ok := cacheApps[appName]; ok {
		return app, nil
	}
	if appName == "" {
		err = errors.New("appName is empty")
		return
	}
	config, err := config.GetWeChatConfig()
	if err != nil {
		err = errors.Wrap(err, "get wechat config error")
		return
	}
	appConfig := config.MiniPrograms[appName]
	if appConfig.AppID == "" || appConfig.AppSecret == "" {
		err = errors.New("app config is empty")
		return
	}
	app, err = miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:     appConfig.AppID,     // 小程序appid
		Secret:    appConfig.AppSecret, // 小程序app secret
		HttpDebug: true,
		Log: miniProgram.Log{
			Level:  "debug",
			File:   "./wechat.log",
			Stdout: true, //  是否打印在终端
		},
	})
	cacheApps[appName] = app
	return
}
