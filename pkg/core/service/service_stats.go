package service

import (
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
)

func (s *DefaultService) Count(args *dto.Param, result *dto.Result) error {
	args.OnlyCount = true
	return s.Retrieve(args, result, iface.MethodQuery)
}

func (s *DefaultService) Aggregate(args *dto.Param, result *dto.Result) error {
	return s.Retrieve(args, result, iface.MethodAggregate)
}

// Partition 实现的是postgresql的窗口函数查询。假设有一张表，记录了大区、年份、月份、营收，四个字段
// 我们需要查询 每个大区在2020~2023年中 每个月的营收 和 当月所有大区的月总营收和年总营收，那么可以如下查询
// select 大区,年份,月份,营收,
//
//		sum(营收) over (partition by 年份,月份) as 月总计,
//	 sum(营收) over (partition by 年份) as 年总计
//
// from 营收表
// where 年份 in ('2020','2021','2023')
// 注意如果传入了distinct ，并且如果需要order by，那么order by中必须出现distinct的字段，并且排在order by语句的最左边
func (s *DefaultService) Partition(args *dto.Param, result *dto.Result) error {
	return s.Retrieve(args, result, iface.MethodPartition)
}
