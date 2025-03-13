package entity

// ColumnDefaultInterface
// 实体实现，返回DB查询的时候默认的列
type ColumnDefaultInterface interface {
	DefaultColumns() []string
}

// ColumnsMustInterface
// 在进行接口查询的时候，通过fill填充返回值的时候，需要用一些必须字段，比如树实体的parent_id
// 但是如果接口参数中没有给这些参数，就可能会导致错误，实现这个接口就避免这些错误
type ColumnsMustInterface interface {
	MustColumns() []string
}
