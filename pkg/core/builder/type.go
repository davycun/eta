package builder

type Builder interface {
	Build() (listSql, countSql string, err error)
}
