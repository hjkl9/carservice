package miniprogram

type Code2SessionResponse struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}

type Code2Session interface {
	Code2session(code string) (*Code2SessionResponse, error)
}
