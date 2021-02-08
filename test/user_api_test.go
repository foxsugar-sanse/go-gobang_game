package test

import (
	"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

const(
	apiV1 string = "/v1"
	apiV2 string = "/v2"
)

func TestUserLogin(t *testing.T) {
	var con config.ConFig = &config.Config{}
	conf := con.InitConfig()
	apiPath := "http://" + conf.ConfData.Run.Ipaddr + conf.ConfData.Run.Port + apiV2
	mux := http.NewServeMux()
	mux.HandleFunc(apiPath + "/user/login",func(w http.ResponseWriter, r *http.Request){})
	req := httptest.NewRequest("Post",apiPath + "/user/login",nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w,req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}
