package dict

var (
	defaultDictionary = make([]Dictionary, 0, 100)
)

func Registry(dictList ...Dictionary) {
	defaultDictionary = append(defaultDictionary, dictList...)
}
