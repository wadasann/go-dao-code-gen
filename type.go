package main

import (
	"strconv"
	"strings"
)

// convertDatabaseTypeToGoType converts a database column type to its corresponding Go type.
// It handles nullable types by returning appropriate sql.Null* types when nullable is true.
func convertDatabaseTypeToGoType(dataType string, nullable bool) string {
	precision := 0

	var unsigned bool
	// Check if the type is unsigned
	if strings.HasSuffix(dataType, " unsigned") {
		unsigned = true
		dataType = dataType[:len(dataType)-len(" unsigned")]
	}

	// Extract precision and scale from the data type
	dataType, precision, _ = extractPrecisionAndScale(dataType)

	var typ string

switchDataType:
	switch dataType {
	case "bit":
		// bit(1) is typically used as boolean
		if precision == 1 {
			typ = "bool"
			if nullable {
				typ = "sql.NullBool"
			}
			break switchDataType
		} else if precision <= 8 {
			typ = "uint8"
		} else if precision <= 16 {
			typ = "uint16"
		} else if precision <= 32 {
			typ = "uint32"
		} else {
			typ = "uint64"
		}
		if nullable {
			typ = "sql.NullInt64"
		}

	case "bool", "boolean":
		typ = "bool"
		if nullable {
			typ = "sql.NullBool"
		}

	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext", "json":
		// All string types map to string
		typ = "string"
		if nullable {
			typ = "sql.NullString"
		}

	case "tinyint":
		// tinyint(1) is commonly used as boolean
		if precision == 1 {
			typ = "bool"
			if nullable {
				typ = "sql.NullBool"
			}
			break
		}
		typ = "int8"
		if unsigned {
			typ = "int16"
		}
		if nullable {
			typ = "sql.NullInt8"
			if unsigned {
				typ = "sql.NullInt16"
			}
		}

	case "smallint":
		typ = "int16"
		if unsigned {
			typ = "int"
		}
		if nullable {
			typ = "sql.NullInt16"
			if unsigned {
				typ = "sql.NullInt"
			}
		}
	case "mediumint", "int", "integer":
		typ = "int"
		if unsigned {
			typ = "int64"
		}
		if nullable {
			typ = "sql.NullInt"
			if unsigned {
				typ = "sql.NullInt64"
			}
		}
	case "bigint":
		typ = "int64"
		if nullable {
			typ = "sql.NullInt64"
		}
	case "float":
		typ = "float32"
		if nullable {
			typ = "sql.NullFloat64"
		}

	case "decimal", "double":
		typ = "float64"
		if nullable {
			typ = "sql.NullFloat64"
		}

	case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		// All binary types map to byte slice
		typ = "[]byte"

	case "timestamp", "datetime", "date":
		typ = "time.Time"
		if nullable {
			typ = "sql.NullTime"
		}

	case "enum", "set", "time":
		// MySQL time type is not directly supported by the driver
		// Users can parse the string to time.Time in their code
		typ = "string"
		if nullable {
			typ = "sql.NullString"
		}

	default:
		// Handle custom types or types with schema prefix
		if strings.HasPrefix(dataType, "mysql"+".") {
			// Remove schema prefix for types in the same schema
			typ = initialisms.SnakeToCamelIdentifier(dataType[len("mysql")+1:])
		} else {
			typ = initialisms.SnakeToCamelIdentifier(dataType)
		}
	}
	return typ
}

// extractPrecisionAndScale extracts precision and scale from a data type string.
// Returns the cleaned data type, precision, and scale values.
// Examples: "varchar(255)" -> ("varchar", 255, -1)
//
//	"decimal(10,2)" -> ("decimal", 10, 2)
func extractPrecisionAndScale(dataType string) (string, int, int) {
	precision := -1
	scale := -1

	// Find the opening parenthesis
	openParen := strings.Index(dataType, "(")
	if openParen == -1 {
		// No precision/scale specified
		return dataType, precision, scale
	}

	// Find the closing parenthesis
	closeParen := strings.Index(dataType[openParen:], ")")
	if closeParen == -1 {
		// Malformed, no closing parenthesis
		return dataType, precision, scale
	}
	closeParen += openParen

	// Extract the content between parentheses
	content := dataType[openParen+1 : closeParen]

	// Split by comma to separate precision and scale
	parts := strings.Split(content, ",")

	// Parse precision (first part)
	if len(parts) > 0 && strings.TrimSpace(parts[0]) != "" {
		if p, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
			precision = p
		}
	}

	// Parse scale (second part, if present)
	if len(parts) > 1 && strings.TrimSpace(parts[1]) != "" {
		if s, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
			scale = s
		}
	}

	// Remove the precision/scale part from the data type
	dataType = dataType[:openParen] + dataType[closeParen+1:]

	// Normalize enum and set types by removing their value lists
	if strings.HasPrefix(dataType, "enum") {
		dataType = "enum"
	}
	if strings.HasPrefix(dataType, "set") {
		dataType = "set"
	}

	return dataType, precision, scale
}
