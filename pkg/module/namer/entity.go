package namer

const (
	TypeUser = "user"
	TypeDept = "dept"
)

type IdName struct {
	ID        string `json:"id,omitempty" gorm:"column:id;comment:实体ID"`
	Name      string `json:"name,omitempty" gorm:"column:name;comment:实体名称"`
	Account   string `json:"account,omitempty" gorm:"column:account;comment:实体账号"`
	Namespace string `json:"namespace,omitempty" gorm:"column:namespace;comment:命名空间针对,同一个app下不同定制化" binding:"required"`
	Category  string `json:"category,omitempty" gorm:"column:category;comment:实体分类"`
	ParentId  string `json:"parent_id,omitempty" gorm:"column:parent_id;comment:父ID"`
	Tp        string `json:"tp,omitempty" `
}
