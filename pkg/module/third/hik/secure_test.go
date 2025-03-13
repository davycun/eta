package hik_test

import (
	"github.com/davycun/eta/pkg/module/third/hik"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/proxy"
	"testing"
)

func M_TestAk(t *testing.T) {

	hk := hik.NewHikClient("192.168.2.211", 446, "25173628", "LMmQbBHQWKmoneaQSYQg")
	hk.SocksProxy("121.43.63.60:23114", &proxy.Auth{User: "mdt", Password: "Mdt123"})

	reqBody := map[string]any{"pageNo": 1, "pageSize": 20}

	rs := make(map[string]interface{})
	err := hk.HttpPost("/artemis/api/resource/v1/cameras", reqBody, &rs, map[string]string{})
	assert.Nil(t, err)
	assert.NotNil(t, rs["data"])

}

func M_TestVideo(t *testing.T) {
	reqBody := map[string]any{
		"cameraIndexCode": "e3c2c2e5baa04084a2f891bd95279f2d",
		"streamType":      0,
		"protocol":        "hls",
		"transmode":       1,
		"expand":          "transcode=1",
		"streamform":      "ps",
	}

	hk := hik.NewHikClient("192.168.2.211", 446, "25173628", "LMmQbBHQWKmoneaQSYQg")
	hk.SocksProxy("121.43.63.60:23114", &proxy.Auth{User: "mdt", Password: "Mdt123"})

	rs := make(map[string]interface{})

	err := hk.HttpPost("/artemis/api/video/v2/cameras/previewURLs", reqBody, &rs, map[string]string{})

	assert.Nil(t, err)
	assert.NotNil(t, rs["data"])

}
