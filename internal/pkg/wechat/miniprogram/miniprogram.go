package miniprogram

import (
	"bytes"
	"carservice/internal/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type MiniProgram interface {
	GetAccessToken() (string, error)
	GetUserPhoneNumber(string, string) (string, error)
	Code2session(code string) (*Code2SessionResponse, error)
}

type MiniProgramProvider struct {
	config config.WechatConf
}

// access_token	string	获取到的凭证
// expires_in	number	凭证有效时间，单位：秒。目前是7200秒之内的值。
// errcode		number	错误码
// errmsg		string	错误信息
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint   `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

//	{
//	    "errcode":0,
//	    "errmsg":"ok",
//	    "phone_info": {
//	        "phoneNumber":"xxxxxx",
//	        "purePhoneNumber": "xxxxxx",
//	        "countryCode": 86,
//	        "watermark": {
//	            "timestamp": 1637744274,
//	            "appid": "xxxx"
//	        }
//	    }
//	}
type PhoneNumberResponse struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	PhoneInfo struct {
		PhoneNumber     string
		PurePhoneNumber string
		CountryCode     int
		// Watermark       struct {
		// 	Timestamp uint
		// 	Appid     string
		// } `json:"watermark"`
	} `json:"phone_info"`
}

func NewWechatProvider(config config.WechatConf) *MiniProgramProvider {
	return &MiniProgramProvider{
		config,
	}
}

// GetAccessToken
// ? Should private?
// return real access token and ?errcode.
func (m *MiniProgramProvider) GetAccessToken() (string, error) {
	apiurl := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		m.config.MiniProgram.AppId,
		m.config.MiniProgram.Secret,
	)
	resp, err := http.Get(apiurl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var response AccessTokenResponse
	json.Unmarshal(body, &response)
	switch response.Errcode {
	case -1:
		return "", errors.New("系统繁忙，请稍后重试")
	case 0: // OK
		return response.AccessToken, nil
	case 40001:
		return "", errors.New("AppSecret 错误或者 AppSecret 不属于这个小程序，请开发者确认 AppSecret 的正确性")
	case 40002:
		return "", errors.New("请确保 grant_type 字段值为 client_credential")
	case 40013:
		return "", errors.New("不合法的 AppID，请开发者检查 AppID 的正确性，避免异常字符，注意大小写")
	default:
		return "", errors.New("其他错误")
	}
}

// GetUserPhoneNumber
// 1. depends on GetAccessToken(...args)
// 2. depends on Code of Frontend.
// return real phone number and ?errcode.
func (m *MiniProgramProvider) GetUserPhoneNumber(accessToken, code string) (string, error) {
	apiurl := "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=" + accessToken
	data := map[string]string{
		"code": code,
	}
	jsonData, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, apiurl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var response PhoneNumberResponse
	json.Unmarshal(body, &response)
	switch response.Errcode {
	case 0:
		return response.PhoneInfo.PhoneNumber, nil
	case -1:
		return "", errors.New("系统繁忙，请稍后重试")
	case 40029:
		return "", errors.New("js_code 无效")
	case 45011:
		return "", errors.New("API 调用太频繁，请稍候再试")
	case 40013:
		return "", errors.New("请求appid身份与获取code的小程序appid不匹配")
	default:
		return "", errors.New("其他错误")
	}
}

func (m *MiniProgramProvider) Code2session(code string) (*Code2SessionResponse, error) {
	apiurl := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		m.config.MiniProgram.AppId,
		m.config.MiniProgram.Secret,
		code,
	)
	resp, err := http.Get(apiurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response Code2SessionResponse
	json.Unmarshal(body, &response)

	return &response, nil
}
