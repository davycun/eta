package weixin

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/httpclient"
	"net/http"
)

type Menu struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	Key        string `json:"key"`
	Url        string `json:"url"`
	Value      string `json:"value"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Digest     string `json:"digest"`
	ShowCover  string `json:"show_cover"`
	CoverUrl   string `json:"cover_url"`
	ContentUrl string `json:"content_url"`
	SourceUrl  string `json:"source_url"`
}
type MenuInfo struct {
	WxError
	IsMenuOpen   int `json:"is_menu_open"`
	SelfMenuInfo struct {
		Button []struct {
			Menu
			SubButton struct {
				List []struct {
					Menu
					NewsInfo struct {
						List []Menu
					}
				}
			} `json:"sub_button"`
		}
	} `json:"selfmenu_info"`
}

func (w *WeiXin) QueryMenu() (MenuInfo, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token=%s", w.GetAccessToken())
	var mi MenuInfo
	hc := httpclient.DefaultHttpClient.
		Url(url).
		Method(http.MethodGet).
		Do(&mi)
	return mi, hc.Error
}
