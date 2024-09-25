package api

import (
	"encoding/json"
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"wxgateway-service/service"
)

func InitMiniProgramRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/mini-program/{app_id}/phonenumber-code-login", phoneNumberCodeLoginHandler)

}

// @Summary 微信小程序登录
// @Description 根据code 获取到手机号，然后判断手机号是否存在，不存在则返回错误，存在则返回access-token
// @Tags Mini-Program
// @Param code query string true "code"
// @Produce  json
// @Success 200 {object} model.TokenInfo "objects array"
// @Failure 500 {object} common.Response ""
// @Router /mini-program/{app_id}/phonenumber-code-login [get]
func phoneNumberCodeLoginHandler(w http.ResponseWriter, r *http.Request) {
	appId := chi.URLParam(r, "app_id")
	code := r.URL.Query().Get("code")
	tokenInfo, err := service.PhoneCodeLogin(r.Context(), appId, code)
	if err != nil {
		common.HttpError(w, common.ErrService.AppendMsg(err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(tokenInfo)
	w.Write(data)
}
