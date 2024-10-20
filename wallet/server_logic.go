package wallet

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"crypto.com/pkg"
)

func NewUser() *User {
	return &User{}
}

// Get get user by name
func (u *User) Get(UserName string) *User {
	if strings.TrimSpace(UserName) == "" {
		return nil
	}

	u.UserName = UserName
	rows, err := pkg.DBEngine.Query(`SELECT id, name, balance FROM "wallet"."user" WHERE "name" = $1`, UserName)
	if err != nil {
		log.Fatal("Error executing query:", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&u.Id, &u.UserName, &u.balance)
		if err != nil {
			log.Fatal("Error scanning row:", err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error with rows:", err)
	}

	return u
}

// Balance Get specify user balance
func (u *User) Balance() uint {
	return u.balance
}

// History Get specify user transaction history
func (u *User) History() []Transaction {
	selectUser := `SELECT id, user_id, receiver_id, action, money, balance, create_time FROM "wallet"."transation" WHERE user_id = $1 ORDER BY id DESC`
	rows, err := pkg.DBEngine.Query(selectUser, u.Id)
	if err != nil {
		log.Fatal("Error executing query:", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := Transaction{}
		err := rows.Scan(&t.Id, &t.SenderId, &t.ReceiverId, &t.Action, &t.Money, &t.Balance, &t.CreateDate)
		if err != nil {
			log.Fatal("Error scanning row:", err)
		} else {
			u.history = append(u.history, t)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error with rows:", err)
	}

	return u.history
}

// Deposit Deposit to specify user wallet
func (u *User) Deposit(money uint) uint {
	tx, err := pkg.DBEngine.Begin()
	if err != nil {
		log.Fatal("Error beginning transaction:", err)
	}
	defer func() {
		if err != nil {
			log.Println("Rolling back transaction due to error")
			_ = tx.Rollback()
		} else {
			log.Println("Committing transaction")
			_ = tx.Commit()
		}
	}()

	if u.Id == 0 {
		insertUser := `INSERT INTO "wallet"."user" ("balance", "name") VALUES ($1, $2) RETURNING id`
		err = tx.QueryRow(insertUser, money, u.UserName).Scan(&u.Id)
		if err != nil {
			log.Fatal("Deposit error in inserting user:", err)
			return 0
		}
		u.balance = money
	} else {
		updateUser := `UPDATE "wallet"."user" SET "balance" = "balance" + $1 WHERE "id" = $2 RETURNING balance`
		err = tx.QueryRow(updateUser, money, u.Id).Scan(&u.balance)
		if err != nil {
			log.Print("Deposit error in updating user:", err)
			return u.balance
		}
	}

	insertHistory := `INSERT INTO "wallet"."transation" (user_id, action, money, balance, create_time) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(insertHistory, u.Id, ActionDeposit, money, u.balance, time.Now().Local().Unix())
	if err != nil {
		log.Fatal("Error inserting transaction history:", err)
	}
	return u.balance
}

// Withdraw Withdraw from specify user wallet
func (u *User) Withdraw(money uint) (uint, bool) {
	tx, err := pkg.DBEngine.Begin()
	if err != nil {
		log.Fatal("Error beginning transaction:", err)
	}
	defer func() {
		if err != nil {
			log.Println("Rolling back transaction due to error")
			_ = tx.Rollback()
		} else {
			log.Println("Committing transaction")
			_ = tx.Commit()
		}
	}()

	if u.Id == 0 {
		return u.balance, false
	}

	updateUser := `UPDATE "wallet"."user" SET balance = balance - $1 WHERE id = $2 AND balance >= $1 RETURNING balance`
	err = tx.QueryRow(updateUser, money, u.Id).Scan(&u.balance)
	if err != nil {
		log.Print("Not enougth money to withdraw, user_id: ", u.Id, " name: ", u.UserName)
		return u.balance, false
	}

	insertHistory := `INSERT INTO "wallet"."transation" (user_id, action, money, balance, create_time) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(insertHistory, u.Id, ActionWithdraw, money, u.balance, time.Now().Local().Unix())
	if err != nil {
		log.Fatal("Error inserting transaction history:", err)
	}

	return u.balance, true
}

// Transfer Transfer from one user to another user
func (u *User) Transfer(receiver *User, money uint) (uint, uint, bool) {
	tx, err := pkg.DBEngine.Begin()
	if err != nil {
		log.Fatal("Error beginning transaction:", err)
	}
	defer func() {
		if err != nil {
			log.Println("Rolling back transaction due to error")
			_ = tx.Rollback()
		} else {
			log.Println("Committing transaction")
			_ = tx.Commit()
		}
	}()

	// update sender balance
	if u.Id == 0 {
		return u.balance, receiver.balance, false
	} else {
		updateUser := `UPDATE "wallet"."user" SET balance = balance - $1 WHERE id = $2 AND balance >= $1 RETURNING balance`
		err = tx.QueryRow(updateUser, money, u.Id).Scan(&u.balance)
		if err != nil {
			log.Print("Not enougth money to transfer, user_id: ", u.Id, " name: ", u.UserName)
			return u.balance, receiver.balance, false
		}
	}

	// update receiver balance
	if receiver.Id == 0 {
		insertUser := `INSERT INTO "wallet"."user" (balance, name) VALUES ($1, $2) RETURNING id`
		err = tx.QueryRow(insertUser, money, receiver.UserName).Scan(&receiver.Id)
		if err != nil {
			log.Fatal("Error inserting user:", err)
			return u.balance, receiver.balance, false
		}
		receiver.balance = money
	} else {
		updateUser := `UPDATE "wallet"."user" SET balance = balance + $1 WHERE id = $2 RETURNING balance`
		err = tx.QueryRow(updateUser, money, receiver.Id).Scan(&receiver.balance)
		if err != nil {
			log.Fatal("Error updating user:", err)
			return u.balance, receiver.balance, false
		}
	}

	// insert history
	now := time.Now().Local().Unix()
	insertUHistory := `INSERT INTO "wallet"."transation" (user_id, receiver_id, action, money, balance, create_time) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(insertUHistory, u.Id, receiver.Id, ActionSend, money, u.balance, now)
	if err != nil {
		log.Fatal("Error inserting transaction history:", err)
	}

	insertRHistory := `INSERT INTO "wallet"."transation" (user_id, receiver_id, action, money, balance, create_time) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(insertRHistory, receiver.Id, u.Id, ActionReceive, money, receiver.balance, now)
	if err != nil {
		log.Fatal("Error inserting transaction history:", err)
	}

	return u.balance, receiver.balance, true
}

// GetName Get user name mapping
func (u *User) GetName(idMap map[uint]string) {
	if len(idMap) == 0 {
		return
	}

	var userId []string
	for id := range idMap {
		userId = append(userId, strconv.Itoa(int(id)))
	}

	selectUser := fmt.Sprintf(`SELECT id, name FROM "wallet"."user" WHERE id IN (%s)`, strings.Join(userId, ","))
	rows, err := pkg.DBEngine.Query(selectUser)
	if err != nil {
		log.Fatal("Error executing query:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id uint
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal("Error scanning row:", err)
		} else {
			idMap[uint(id)] = name
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error with rows:", err)
	}
}

// Reset only for testing to clear data
func (u *User) Reset(mu ...*User) {
	if len(mu) > 1 {
		tx, _ := pkg.DBEngine.Begin()
		for i := range mu {
			if mu[i].Id > 0 {
				tx.Exec(`DELETE FROM "wallet"."user" WHERE id = $1`, mu[i].Id)
				tx.Exec(`DELETE FROM "wallet"."transation" WHERE user_id = $1`, mu[i].Id)
				mu[i].Id = 0
			}
		}
		tx.Commit()
	} else if u.Id > 0 {
		pkg.DBEngine.Exec(`DELETE FROM "wallet"."user" WHERE id = $1`, u.Id)
		pkg.DBEngine.Exec(`DELETE FROM "wallet"."transation" WHERE user_id = $1`, u.Id)
		u.Id = 0
	}
}
