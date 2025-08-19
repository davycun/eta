package server

import (
	"github.com/davycun/eta/pkg/common/config"
	"github.com/gin-gonic/gin/binding"
)

var (
	confFile  = ""
	argConfig = config.Configuration{}
	destCfg   = &config.Configuration{}
)

func init() {
	config.BindArgConfig(StartCommand, &confFile, &argConfig)
}

func readConfig() error {
	destCfg = config.LoadConfig(confFile, &argConfig)
	//如果json串中有一个字段是数值，但是反序列化的时候针对这个字段没有指定具体的是float或者int
	//那么默认json会反序列化为float64类型，这也就是为什么我用map去接受反序列化的时候，明明序列化之前是int，但是反序列化后map里面是float64的原因
	//如果设置了EnableDecoderUseNumber，那么这种情况下反序列化的目标就会被指定为json.Number对象（其实是个string，type Number string）
	//binding.EnableDecoderUseNumber = true
	binding.EnableDecoderDisallowUnknownFields = true

	return nil
}
