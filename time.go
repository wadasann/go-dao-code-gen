package main

type TimeType string

const (
	TimeTypeDatetime TimeType = "datetime"
	TimeTypeInt      TimeType = "int"
)

var (
	createTimeMap = map[string]struct{}{
		"create_time":   {},
		"create_at":     {},
		"create_date":   {},
		"create_on":     {},
		"creation_time": {},
		"creation_at":   {},
		"creation_date": {},
		"creation_on":   {},
		"created_time":  {},
		"created_at":    {},
		"created_date":  {},
		"created_on":    {},
		"add_time":      {},
		"add_at":        {},
		"add_date":      {},
		"insert_time":   {},
		"insert_at":     {},
		"insert_date":   {},
		"insert_on":     {},
		"inserted_at":   {},
		"inserted_on":   {},
		"ctime":         {},
		"c_time":        {},
	}
	updateTimeMap = map[string]struct{}{
		"update_time":   {},
		"update_at":     {},
		"update_date":   {},
		"update_on":     {},
		"updated_time":  {},
		"updated_at":    {},
		"updated_date":  {},
		"updated_on":    {},
		"modify_time":   {},
		"modify_at":     {},
		"modify_date":   {},
		"modify_on":     {},
		"modified_time": {},
		"modified_at":   {},
		"modified_date": {},
		"modified_on":   {},
		"edit_time":     {},
		"edit_at":       {},
		"edit_date":     {},
		"edit_on":       {},
		"edited_time":   {},
		"edited_at":     {},
		"edited_date":   {},
		"edited_on":     {},
		"utime":         {},
		"u_time":        {},
	}
)

type TimeFields struct {
	CreateTime string
	CreateType TimeType
	UpdateTime string
	UpdateType TimeType
}
