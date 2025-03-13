package auth

var (
	defaultAuth2Role = make([]Auth2Role, 0, 10)
)

func RegistryDefaultAuth2role(a2rList ...Auth2Role) {
	for _, v := range a2rList {
		defaultAuth2Role = append(defaultAuth2Role, v)
	}
}
