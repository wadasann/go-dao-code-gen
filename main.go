package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/kenshaw/snaker"
)

const (
	shadowPrefix = "fct_"
	shadowSuffix = "_shadow"

	Version = "2025-07-09.v5"
)

var (
	err         error
	tablesMap   map[string]struct{}
	mysqlDB     *sql.DB
	initialisms *snaker.Initialisms
)

func init() {
	tablesMap = make(map[string]struct{})
	initialisms, err = snaker.NewDefaultInitialisms()
	if err != nil {
		slog.Error("error: create initialisms failed", "error", err)
		os.Exit(1)
	}
	for i := rune('a'); i <= rune('z'); i++ {
		initialisms.Add(strings.ToUpper(string(i) + "id"))
	}
}

type AttrEntity struct {
	Name           string
	NameCamel      string
	NameCamelIdent string
	Type           string
	Tag            string
	Comment        string
	IsPk           bool
	HasIndex       bool
}

func main() {
	parseFlags()
	slog.Info("open database connection")
	ctx := context.Background()
	mysqlDB, err = sql.Open("mysql", dsn)
	if err != nil {
		slog.Error(fmt.Sprintf("error: open db failed, %s", err.Error()))
		os.Exit(1)
	}
	defer mysqlDB.Close()
	slog.Info("get tables")
	tableRows, err := mysqlDB.Query("show tables;")
	if err != nil {
		slog.Error("error: exec `show tables` failed", "error", err)
		os.Exit(1)
	}
	defer tableRows.Close()
	tableIndex := make(map[string]struct{})
	shadowTables := make(map[string]string)
	for tableRows.Next() {
		var table string
		tableRows.Scan(&table)
		tableIndex[table] = struct{}{}
	}
	var tables []string
	for index := range tableIndex {
		var table string
		if strings.HasPrefix(index, shadowPrefix) {
			table = strings.TrimPrefix(index, shadowPrefix)
		} else if strings.HasSuffix(index, shadowSuffix) {
			table = strings.TrimSuffix(index, shadowSuffix)
		} else {
			tables = append(tables, index)
			continue
		}
		if _, ok := tableIndex[table]; !ok {
			tables = append(tables, index)
			continue
		}
		shadowTables[index] = table
	}
	if len(tables) == 0 {
		slog.Error("error: no tables be found in database")
		os.Exit(1)
	}
	slog.Info("gen dao.go")
	pkg := filepath.Base(outputDir)
	err = genInitDao(ctx, pkg, shadowTables)
	if err != nil {
		println(err.Error())
	}
	slog.Info("gen tables")
	for _, table := range tables {
		if len(tablesMap) != 0 {
			if _, ok := tablesMap[table]; !ok {
				continue
			}
		}
		columns, err := getTableColumns(ctx, table)
		if err != nil {
			slog.Error("Get table columns failed", "error", err)
			continue
		}
		indexes, err := getTableIndexes(ctx, table)
		if err != nil {
			slog.Error("Get table indexes failed", "error", err)
			continue
		}
		rData, imports, err := getRenderData(ctx, pkg, table, columns, indexes)
		if err != nil {
			slog.Error("Get render data failed", "error", err)
			continue
		}
		if rData == nil {
			continue
		}
		slog.Info(fmt.Sprintf("gen table %s \n", table))
		err = genTable(ctx, table, rData, imports)
		if err != nil {
			println(err.Error())
			continue
		}
		slog.Info(fmt.Sprintf("gen table conds %s \n", table))
		err = genTableConds(ctx, table, rData, imports)
		if err != nil {
			println(err.Error())
			continue
		}
	}
}

const (
	MySQLV5IndexNum = 13
	MySQLV8IndexNum = 15
)

type IndexEntityV5 struct {
	Table        string
	NonUnique    bool
	KeyName      string
	SeqInIndex   int
	ColumnName   string
	Collation    string
	Cardinality  int
	SubPart      sql.NullString
	Packed       sql.NullString
	Null         string
	IndexType    string
	Comment      string
	IndexComment string
}

type IndexEntityV8 struct {
	IndexEntityV5
	Visible    string
	Expression sql.NullString
}

func getTableIndexes(ctx context.Context, table string) (indexes []*IndexEntityV5, err error) {
	indexRows, err := mysqlDB.Query(fmt.Sprintf("show index from `%s`", table))
	if err != nil {
		return indexes, fmt.Errorf("error:  failed to get indexes from table %s, %v", table, err)
	}
	defer indexRows.Close()
	cols, _ := indexRows.Columns()
	if len(cols) == MySQLV8IndexNum {
		for indexRows.Next() {
			indexEntity := IndexEntityV8{}
			err = indexRows.Scan(&indexEntity.Table, &indexEntity.NonUnique, &indexEntity.KeyName,
				&indexEntity.SeqInIndex, &indexEntity.ColumnName, &indexEntity.Collation,
				&indexEntity.Cardinality, &indexEntity.SubPart, &indexEntity.Packed,
				&indexEntity.Null, &indexEntity.IndexType, &indexEntity.Comment,
				&indexEntity.IndexComment, &indexEntity.Visible, &indexEntity.Expression)
			if err != nil {
				return
			}
			indexes = append(indexes, &indexEntity.IndexEntityV5)
		}
	} else {
		for indexRows.Next() {
			indexEntity := IndexEntityV5{}
			indexStruct := sqlbuilder.NewStruct(indexEntity)
			err = indexRows.Scan(indexStruct.Addr(&indexEntity)...)
			if err != nil {
				return
			}
			indexes = append(indexes, &indexEntity)
		}
	}

	return
}

type ColumnEntity struct {
	Field      string
	Type       string
	Collation  sql.NullString
	Null       string
	Key        string
	Default    sql.NullString
	Extra      string
	Privileges string
	Comment    string
}

func getTableColumns(ctx context.Context, table string) (columns []*ColumnEntity, err error) {
	columnRows, err := mysqlDB.Query(fmt.Sprintf("show full columns from `%s`", table))
	if err != nil {
		return columns, fmt.Errorf("error:  failed to get full columns from table %s, %v", table, err)
	}
	defer columnRows.Close()
	for columnRows.Next() {
		columnEntity := &ColumnEntity{}
		columnStruct := sqlbuilder.NewStruct(columnEntity)
		columnRows.Scan(columnStruct.Addr(columnEntity)...)
		columns = append(columns, columnEntity)
	}
	return
}

func getRenderData(ctx context.Context, pkg, table string,
	columns []*ColumnEntity, indexes []*IndexEntityV5) (rData *RenderData, imports []string, err error) {
	if len(columns) == 0 {
		return
	}
	attrs := make([]*AttrEntity, 0, len(columns))
	var primary string
	var timeFields TimeFields
	importsMap := make(map[string]struct{})
	for _, column := range columns {
		nullable := false
		if column.Null == "YES" {
			nullable = true
		}
		dt := convertDatabaseTypeToGoType(column.Type, nullable)
		if strings.HasPrefix(dt, "sql.") {
			if _, ok := importsMap["database/sql"]; !ok {
				importsMap["database/sql"] = struct{}{}
				imports = append(imports, "database/sql")
			}
		}
		if strings.HasPrefix(dt, "time.") {
			if _, ok := importsMap["time"]; !ok {
				importsMap["time"] = struct{}{}
				imports = append(imports, "time")
			}
		}
		if strings.HasPrefix(dt, "decimal.") {
			if _, ok := importsMap["github.com/shopspring/decimal"]; ok {
				importsMap["github.com/shopspring/decimal"] = struct{}{}
				imports = append(imports, "github.com/shopspring/decimal")
			}
		}
		var hasIndex, isPk bool
		if column.Key != "" {
			hasIndex = true
			if column.Key == "PRI" {
				isPk = true
				primary = column.Field
			}
		}
		attr := &AttrEntity{
			Name:           initialisms.SnakeToCamelIdentifier(column.Field),
			NameCamel:      replaceReserved(initialisms.ForceLowerCamelIdentifier(column.Field)),
			NameCamelIdent: initialisms.ForceCamelIdentifier(column.Field),
			Type:           dt,
			Tag:            column.Field,
			Comment:        column.Comment,
			IsPk:           isPk,
			HasIndex:       hasIndex,
		}
		attrs = append(attrs, attr)
		if _, ok := createTimeMap[column.Field]; ok {
			var timeType TimeType
			columnType := strings.ToLower(column.Type)
			columnType = strings.SplitN(columnType, " ", 2)[0]
			columnType = strings.SplitN(columnType, "(", 2)[0]
			switch columnType {
			case "datetime", "timestamp":
				timeType = TimeTypeDatetime
			case "int", "bigint":
				timeType = TimeTypeInt
			default:
				continue
			}
			timeFields.CreateTime = column.Field
			timeFields.CreateType = timeType
			continue
		}
		if _, ok := updateTimeMap[column.Field]; ok {
			var timeType TimeType
			columnType := strings.ToLower(column.Type)
			columnType = strings.SplitN(columnType, " ", 2)[0]
			columnType = strings.SplitN(columnType, "(", 2)[0]
			switch columnType {
			case "datetime", "timestamp":
				timeType = TimeTypeDatetime
			case "int", "bigint":
				timeType = TimeTypeInt
			default:
				continue
			}
			timeFields.UpdateTime = column.Field
			timeFields.UpdateType = timeType
		}
	}
	idxs := make(Indexes)
	for _, index := range indexes {
		if index.NonUnique {
			continue
		}
		idxs[index.KeyName] = append(idxs[index.KeyName], index.ColumnName)
	}
	rData = &RenderData{
		Pkg:                  pkg,
		Table:                table,
		TableLowerCamelIdent: initialisms.ForceLowerCamelIdentifier(table),
		TableUpperCamelIdent: initialisms.ForceCamelIdentifier(table),
		Primary:              primary,
		Attrs:                attrs,
		UniqueIndexes:        idxs,
		TimeFields:           timeFields,
	}
	return
}

func genTable(ctx context.Context, table string, rData *RenderData, imports []string) error {
	if !slices.Contains(imports, "database/sql") {
		imports = append(imports, "database/sql")
	}
	if !slices.Contains(imports, "time") && (rData.TimeFields.CreateTime != "" || rData.TimeFields.UpdateTime != "") {
		imports = append(imports, "time")
	}
	rData.Imports = imports
	content, err := renderTable(table, rData)
	if err != nil {
		return fmt.Errorf("error: render table %s tpl failed, %v", table, err)
	}
	f, err := getTableFile(table)
	if err != nil {
		return fmt.Errorf("error: generate table %s file failed, %v", table, err)
	}

	f.Write(content)
	f.Close()
	return nil
}

func genTableConds(ctx context.Context, table string, rData *RenderData, imports []string) error {
	rData.Imports = imports
	content, err := renderTableConds(table, rData)
	if err != nil {
		return fmt.Errorf("error: render table %s conds tpl failed, %v", table, err)
	}
	f, err := getTableCondsFile(table)
	if err != nil {
		return fmt.Errorf("error: generate table %s conds file failed, %v", table, err)
	}

	f.Write(content)
	f.Close()
	return nil
}

func genInitDao(ctx context.Context, pkg string, shadowTables map[string]string) error {
	renderData := &RenderData{
		Pkg:          pkg,
		ShadowTables: shadowTables,
	}
	content, err := renderInitDao(renderData)
	if err != nil {
		return fmt.Errorf("error: render dao tpl failed, %v", err)
	}
	f, err := getInitDaoFile()
	if err != nil {
		// if dao.go file already exist, return directly
		if err == ErrFileAlreadyExists {
			return err
		}
		return fmt.Errorf("error: generate dao file failed, %v", err)
	}

	f.Write(content)
	f.Close()
	return nil
}
