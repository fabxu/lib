package sqldb

import (
	"regexp"
	"strings"

	"gorm.io/gorm"
)

func SQLEscape(value string) string {
	pattern := regexp.MustCompile(`(?m)(['"%_\n\r\t\\\x00\x08\x1a])`)
	return pattern.ReplaceAllString(value, "\\$0")
}

func QuoteTo(tx *gorm.DB, names ...string) string {
	var builder strings.Builder

	quotedNames := make([]string, len(names))

	for i, name := range names {
		tx.Dialector.QuoteTo(&builder, name)
		quotedNames[i] = builder.String()
		builder.Reset()
	}

	return strings.Join(quotedNames, ".")
}

func CheckDBType(db *gorm.DB) DBType {
	if db.Config.Dialector.Name() == string(DBTypePostgres) {
		return DBTypePostgres
	}

	return DBTypeMySQL
}
