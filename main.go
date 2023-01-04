// Package main is an example application.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/bool64/sqluct"
	_ "github.com/lib/pq"
)

// User is a sample DB row structure.
type User struct {
	ID     string `db:"id"`
	Name   string `db:"name"`
	Weight int    `db:"weight"`
}

// Order is a sample DB row structure.
type Order struct {
	ID     int    `db:"id"`
	UserID string `db:"user_id"`
	Amount int    `db:"amount"`
}

func main() {
	var foo int

	flag.IntVar(&foo, "foo", 123, "Foo is an example flag.")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: dbcli <user_id> <order_id>")
		flag.PrintDefaults()

		return
	}

	userID := flag.Arg(0)

	orderID, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.Fatal("invalid order_id: ", err)
	}

	connStr := "postgres://pqgotest:password@localhost/pqgotest?sslmode=disable"

	st, err := sqluct.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to connect: ", err)
	}

	ctx := context.Background()

	user, err := sqluct.Get[User](ctx, st, sqluct.Stmt(`SELECT * FROM "user" WHERE "id" = $1`, userID))
	if err != nil {
		log.Fatal("user query failed: ", err)
	}

	fmt.Println(user.Name)

	orders, err := sqluct.List[Order](ctx, st,
		sqluct.Stmt(`SELECT * FROM "order" WHERE "user_id" = $1 AND "id" != $2`, userID, orderID))
	if err != nil {
		log.Fatal("orders query failed:", err.Error())
	}

	for _, o := range orders {
		fmt.Println(o.ID, o.Amount)
	}
}
