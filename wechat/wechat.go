package wechat

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func init() {
	appid := ""
	secret := ""
	freshTokenTicker := time.NewTicker(7000 * time.Second)
	//requestToken()

	go func() {
		select {
		case <-freshTokenTicker.C:
			{
				accessToken, err := requestToken(appid, secret)
				if err != nil {
					//TODO 错误处理
				}
				log.Printf("token refresh :%s", accessToken)
			}
		}
	}()

}
func requestToken(appid, secret string) (string, error) {
	u, err := url.Parse("https://api.weixin.qq.com/cgi-bin/token")
	if err != nil {
		log.Fatal(err)
	}
	paras := &url.Values{}
	//设置请求参数
	paras.Set("appid", appid)
	paras.Set("secret", secret)
	paras.Set("grant_type", "client_credential")
	u.RawQuery = paras.Encode()
	resp, err := http.Get(u.String())
	//关闭资源
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", errors.New("request token err :" + err.Error())
	}

	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return "", errors.New("request token response json parse err :" + err.Error())
	}
	if jMap["errcode"] == nil || jMap["errcode"] == 0 {
		accessToken, _ := jMap["access_token"].(string)
		return accessToken, nil
	} else {
		//返回错误信息
		errcode := jMap["errcode"].(string)
		errmsg := jMap["errmsg"].(string)
		err = errors.New(errcode + ":" + errmsg)
		return "", err
	}

}

func WechatLogin(js_code, appid, secret string) (map[string]interface{}, error) {

	Code2SessURL := "https://api.weixin.qq.com/sns/jscode2session?appid={appid}&secret={secret}&js_code={code}&grant_type=authorization_code"
	Code2SessURL = strings.Replace(Code2SessURL, "{appid}", appid, -1)
	Code2SessURL = strings.Replace(Code2SessURL, "{secret}", secret, -1)
	Code2SessURL = strings.Replace(Code2SessURL, "{code}", js_code, -1)
	resp, err := http.Get(Code2SessURL)
	//关闭资源
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.New("WechatLogin request err :" + err.Error())
	}

	var jMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jMap)

	if err != nil {
		return nil, errors.New("request token response json parse err :" + err.Error())

	}
	if jMap["errcode"] == nil || jMap["errcode"] == 0 {

		return jMap, nil
	} else {
		//返回错误信息
		errcode := jMap["errcode"].(string)
		errmsg := jMap["errmsg"].(string)
		err = errors.New(errcode + ":" + errmsg)
		return nil, err
	}
}
