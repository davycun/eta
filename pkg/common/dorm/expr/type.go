package expr

const (
	VarTypeColumn           = "column"
	VarTypeValue            = "value"
	Parenthesis   QuoteType = "()"
	Bracket       QuoteType = "[]"
	Brace         QuoteType = "{}"
	SingleQuote   QuoteType = "'"
	DoubleQuote   QuoteType = `"`
)

type QuoteType string

type Expression struct {
	Expr string   `json:"expr,omitempty"` //表达式
	Vars []ExpVar `json:"vars,omitempty"`
}

type ExpVar struct {
	Type  string `json:"type,omitempty" binding:"oneof=column value ''"` //value
	Value any    `json:"value,omitempty"`
}

type ExpColumn struct {
	Expression
	Type  string `json:"type,omitempty"`  //取出字段后的类型
	Alias string `json:"alias,omitempty"` //取出字段后的别名
}
