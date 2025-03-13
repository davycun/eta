package utils_test

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFormatStrToTime(t *testing.T) {
	strs := []string{"1999-07-16", "2024-07-26T12:07:55.654321+08:00", "2024-07-26T12:07:55+08:00"}

	for _, str := range strs {
		toTime, err := utils.FormatStrToTime(str)
		assert.NoError(t, err)
		logger.Infof("utils.FormatStrToTime: %v", toTime)
	}

	now := time.Now()
	zone, offset := now.Zone()
	logger.Infof("now: %v, name: %v, offset: %v", now, zone, offset)

	parse, err := time.Parse("2006-01-02", "1999-07-16")
	assert.NoError(t, err)
	logger.Infof("time.Parse: %v", parse)

}
