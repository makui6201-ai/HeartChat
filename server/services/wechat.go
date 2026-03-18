package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const wxCode2SessionURL = "https://api.weixin.qq.com/sns/jscode2session"

type wxSessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// WeChatLogin exchanges a WeChat login code for the user's stable openid.
// It requires WECHAT_APPID and WECHAT_SECRET environment variables.
// httpClient is the shared *http.Client defined in deepseek.go (same package).
func WeChatLogin(code string) (openID string, err error) {
	appID := os.Getenv("WECHAT_APPID")
	secret := os.Getenv("WECHAT_SECRET")
	if appID == "" || secret == "" {
		return "", fmt.Errorf("WECHAT_APPID and WECHAT_SECRET must be set")
	}

	url := fmt.Sprintf(
		"%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		wxCode2SessionURL, appID, secret, code,
	)

	resp, err := httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("wx code2session request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("wx code2session read body: %w", err)
	}

	var wxResp wxSessionResponse
	if err := json.Unmarshal(body, &wxResp); err != nil {
		return "", fmt.Errorf("wx code2session parse: %w", err)
	}
	if wxResp.ErrCode != 0 {
		return "", fmt.Errorf("wx code2session error %d: %s", wxResp.ErrCode, wxResp.ErrMsg)
	}
	return wxResp.OpenID, nil
}
