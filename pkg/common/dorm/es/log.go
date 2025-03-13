package es

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"net/http"
	"strings"
	"time"
)

const (
	optSearch = "search"
	optUpsert = "upsert"
	optDelete = "delete"
)

func LatencyLog(start time.Time, idx string, opt string, body []byte, statusCode int) {

	var (
		latency = time.Now().Sub(start)
	)

	logger.Infof("| elasticsearch | %s | %13v |%s| %s\n%s",
		getCodeColor(statusCode), latency, utils.FmtTextBlue(idx), opt, utils.BytesToString(body),
	)
}

func getCodeColor(code int) string {

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return utils.FmtTextGreen(fmt.Sprintf("%3d", code))
	default:
		return utils.FmtTextRed(fmt.Sprintf("%3d", code))
	}
}

func getSearchResultCode(err error) int {
	if err != nil {
		if x, ok := err.(*types.ElasticsearchError); ok {
			return x.Status
		} else {
			return 500
		}
	}
	return 200
}
func getBulkResultCode(resp *bulk.Response) int {
	if resp != nil && resp.Errors {
		for _, v := range resp.Items {
			for _, y := range v {
				return y.Status
			}
		}
	}
	return 200
}
func getBulkErrorMsg(resp *bulk.Response, opt string) string {
	bd := strings.Builder{}
	if resp != nil && resp.Errors {
		for _, v := range resp.Items {
			for _, y := range v {
				if y.Status != http.StatusOK {
					bd.WriteString(fmt.Sprintf("the document which id is %s and index is %s %s fail because %s \n", y.Id_, y.Index_, opt, getErrorMsg(y.Error)))
				}
			}
		}
	}
	return bd.String()
}
