# go-dao-code-gen

A command-line tool for generating Go data access layer (DAO) code from MySQL database tables.

## Overview

`go-dao-code-gen` is a useful code generation tool that automatically creates Go DAO (Data Access Object) layer code based on your MySQL database schema. It generates type-safe Go structs and SQL operations for each table in your database, making database interactions more efficient and less error-prone.

## Features

- **Automatic Code Generation**: Generates Go structs and DAO methods from MySQL table schemas
- **Type Safety**: Converts MySQL data types to appropriate Go types
- **SQL Builder Integration**: All SQL operations are built using [go-sqlbuilder](https://github.com/huandu/go-sqlbuilder) for safe and efficient query construction
- **Index Support**: Automatically detects and handles database indexes
- **Flexible Output**: Generate code for all tables or specify specific tables
- **Connection Flexibility**: Support for both individual connection parameters and DSN strings

## Dependencies

This project leverages several excellent open-source libraries:

- **[huandu/go-sqlbuilder](https://github.com/huandu/go-sqlbuilder)**: Used for building type-safe SQL queries and operations
- **[xo/dbtpl](https://github.com/xo/dbtpl)**: Referenced for type conversion patterns and best practices
- **[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)**: MySQL driver for Go
- **[kenshaw/snaker](https://github.com/kenshaw/snaker)**: Utility for converting between snake_case and CamelCase naming conventions

## Installation

Install the tool using `go install`:

```bash
go install github.com/wadasann/go-dao-code-gen@latest
```

## Usage

### Basic Usage

Generate DAO code for all tables in a database:

```bash
# Using individual connection parameters
go-dao-code-gen -h 127.0.0.1 -u root -p 123456 -D dbname -o ./dao

# Using DSN connection string
go-dao-code-gen -dsn='user:passwd@tcp(127.0.0.1:3306)/dbname?charset=utf8' -o ./dao
```

### Generate Code for Specific Tables

Generate DAO code only for specified tables:

```bash
go-dao-code-gen -h 127.0.0.1 -u user -p 123456 -D dbname -o ./dao -tables "users,orders,products"
```

### Command Line Options

```bash
$ go-dao-code-gen -help

Usage 1:
	go-dao-code-gen -h 127.0.0.1 -u user -p 123456 -D dbname -o ./dao

	or

	go-dao-code-gen -dsn='user:passwd@tcp(127.0.0.1:3306)/dbname?charset=utf8' -o ./dao

Usage 2(Specified tables):
	go-dao-code-gen -h 127.0.0.1 -u user -p passwd -D dbname -o ./dao -tables "table1,table2"

  -D string
    	Database to use.
  -P string
    	Port number to use for connection. (default "3306")
  -dsn string
    	Mysql dsn connection string.
  -h string
    	Connect to host. (default "127.0.0.1")
  -help
    	Show command usage.
  -o string
    	Output directory.
  -p string
    	Password to use when connecting to server.
  -params string
    	Connection parameters.
  -tables string
    	Generation range of tables, use "," separate multiple tables.
  -u string
    	User for login if not root user. (default "root")
  -v	Show command version.
```

## Generated Code Usage

After generating the DAO code, you can use it in your Go application like this:

```go
package main

import (
	"context"
	"github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "user:passwd@(127.0.0.1:3306)/dbname"
	mysqlConfig, err := mysql.ParseDSN(dsn)
	if err != nil {
		// Handle error
	}
	
	ctx := context.Background()
	if err = dao.Init(ctx, mysqlConfig); err != nil {
		// Handle error
	}
	
	defer dao.Close()
	
	// Use generated DAO methods here
	// example: 
	userDao := dao.NewUserDao()
	user, err := userDao.Get(ctx, dao.SetUserEmail("foo@bar.com"))
	if err != nil {
		// Handle error
	}
	if user == nil {
		// User not found
	}
	println(user.Email)

}
```

## Generated Files

The tool generates the following files in your output directory:

- `dao.go`: Main DAO initialization and connection management
- `{table_name}.go`: Individual table DAO with CRUD operations
- `{table_name}_conds.go`: Condition builders for complex queries

## Type Conversion

The type conversion logic is inspired by the excellent [dbtpl](https://github.com/xo/dbtpl) project, which provides robust MySQL to Go type mapping. 

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.

## License

This package is licensed under the MIT license. For more information, refer to the LICENSE file.