package passkey

import (
	"github.com/SchoolAF/exodus/model"
	"github.com/SchoolAF/exodus/repository/db"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"time"
)

func Generate() int {
	rand.Seed(time.Now().UnixNano())
	return 100000 + rand.Intn(900000)
}

func Update(accountData model.Accounts, filter bson.M) error {
	if err := db.UpdateAccount(accountData, filter); err != nil {
		return fmt.Errorf("Failed to update account in the database: %w", err)
	}
	return nil
}
