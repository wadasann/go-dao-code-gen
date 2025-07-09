package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	mysqlDrv "github.com/go-sql-driver/mysql"
)

const (
	cmdQuote = `
Usage 1:
	go-dao-code-gen -h 127.0.0.1 -u root -p 123456 -D dbname -o ./dao 

	or

	go-dao-code-gen -dsn='root:123456@tcp(127.0.0.1:3306)/dbname?charset=utf8' -o ./dao

Usage 2(Specified tables):
	go-dao-code-gen -h 127.0.0.1 -u root -p 123456 -D dbname -o ./dao -tables "tbl1,tbl2"
`
)

var (
	help, version bool

	dsn string
	// host, user, passwd, port, dbname, params
	h, u, p, P, D, params string // When dsn parameters are present, these parameters are ignored
	tables                string

	outputDir string // Output directory
)

func parseFlags() {
	// Help
	flag.BoolVar(&help, "help", false, "Show command usage.")
	flag.BoolVar(&version, "version", false, "Show command version.")

	// Database config
	flag.StringVar(&dsn, "dsn", "", "Mysql dsn connection string.")

	flag.StringVar(&h, "h", "127.0.0.1", "Connect to host.")

	flag.StringVar(&u, "u", "root", "User for login if not root user.")

	flag.StringVar(&p, "p", "", "Password to use when connecting to server.")

	flag.StringVar(&P, "P", "3306", "Port number to use for connection.")

	flag.StringVar(&D, "D", "", "Database to use.")

	flag.StringVar(&params, "params", "", "Connection parameters.")

	flag.StringVar(&tables, "tables", "", "Generation range of tables, use \",\" separate multiple tables.")

	// Output config
	flag.StringVar(&outputDir, "o", "", "Output directory.")

	flag.Parse()
	// Validate flag vars
	if help {
		fmt.Println(cmdQuote)
		flag.PrintDefaults()
		os.Exit(0)
	}
	if version {
		fmt.Println(Version)
		os.Exit(0)
	}
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		if D == "" {
			if u != "root" {
				D = u
			} else {
				fmt.Printf("Error: missed -D or -database flag to set database name.\n")
				os.Exit(1)
			}
		}
		connParams := map[string]string{}
		if params != "" {
			values, err := url.ParseQuery(params)
			if err != nil {
				fmt.Printf("Error: Parse params failed, %v.\n", err)
				os.Exit(1)
			}
			for index, val := range values {
				connParams[index] = val[0]
			}
		}
		config := mysqlDrv.NewConfig()
		config.User = u
		config.Passwd = p
		config.Net = "tcp"
		config.Addr = net.JoinHostPort(h, P)
		config.DBName = D
		config.Params = connParams
		dsn = config.FormatDSN()
	}
	fmt.Println(dsn)
	outputDir = strings.TrimSpace(outputDir)
	var err error
	if outputDir == "" {
		outputDir, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error: get current path failed.\n")
			os.Exit(1)
		}
	}
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		fmt.Printf("Error: get absolute output path failed.\n")
		os.Exit(1)
	}
	tables = strings.TrimSpace(tables)
	if tables != "" {
		tablesList := strings.Split(tables, ",")
		for _, table := range tablesList {
			table = strings.TrimSpace(table)
			if table == "" {
				continue
			}
			tablesMap[table] = struct{}{}
		}
	}
}
