package controllers

import (
	"database/sql"
	"fmt"

	"github.com/SelfServiceCo/api/pkg/models"
	_ "github.com/go-sql-driver/mysql"
)

func GetRestaurantOrders(rid int64) []models.Order {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	//
	results, err := db.Query("SELECT * FROM orders WHERE RID = ? AND STATUS != 'DONE' AND STATUS != 'DENY'", rid)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	orders := []models.Order{}
	for results.Next() {
		var order models.Order
		err = results.Scan(&order.ID, &order.ResID, &order.UserID, &order.TableID, &order.Details, &order.Status)
		if err != nil {
			panic(err.Error())
		}
		orders = append(orders, order)
	}

	return orders
}

func ChangeOrderStatus(oid int64, rid int64, status string) bool {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return false
	}

	_, err = db.Exec("UPDATE orders SET STATUS = ? WHERE RID = ? AND OID = ?", status, rid, oid)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return false
	}

	return true
}

func GetOrdersByUser(uid int64) []models.Order {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	results, err := db.Query("SELECT * FROM orders WHERE UID = ?", uid)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	orders := []models.Order{}
	for results.Next() {
		var order models.Order
		err = results.Scan(&order.ID, &order.ResID, &order.UserID, &order.TableID, &order.Details, &order.Status)
		if err != nil {
			panic(err.Error())
		}
		orders = append(orders, order)
	}

	return orders
}
