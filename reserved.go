package main

const (
	reservedSuffix = "Reserved"
)

var reservedKeywords = map[string]struct{}{
	"break":       {},
	"default":     {},
	"func":        {},
	"interface":   {},
	"select":      {},
	"case":        {},
	"defer":       {},
	"go":          {},
	"map":         {},
	"struct":      {},
	"chan":        {},
	"else":        {},
	"goto":        {},
	"package":     {},
	"switch":      {},
	"const":       {},
	"fallthrough": {},
	"if":          {},
	"range":       {},
	"type":        {},
	"continue":    {},
	"for":         {},
	"import":      {},
	"return":      {},
	"var":         {},
}

func replaceReserved(fieldName string) string {
	if _, ok := reservedKeywords[fieldName]; !ok {
		return fieldName
	}
	return fieldName + reservedSuffix
}
