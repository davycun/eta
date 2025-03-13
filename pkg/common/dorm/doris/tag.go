package doris

import "github.com/davycun/eta/pkg/common/tag"

const (
	TagName = "doris"
)

type Tag struct {
	tg tag.Tag
}

func NewTag(text string) Tag {
	tg := Tag{
		tg: tag.New(text),
	}

	return tg
}

func (t Tag) AggType() string {
	return t.tg.Get("agg_type")
}
