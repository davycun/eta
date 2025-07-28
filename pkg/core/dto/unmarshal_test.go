package dto

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshall(t *testing.T) {
	st := `{
    "filters": [
        {
            "logical_operator": "and",
            "column": "category",
            "operator": "=",
            "value": "setting_address_category"
        },
        {
            "logical_operator": "and",
            "column": "name",
            "operator": "=",
            "value": "setting_address_name"
        }
    ]
}`
	pm := &Param{}
	err := json.Unmarshal([]byte(st), pm)
	assert.Nil(t, err)
}
