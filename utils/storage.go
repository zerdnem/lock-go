package utils

import (
	"math/rand"

	"github.com/jinzhu/gorm"
)

var db, _ = gorm.Open("postgres", "")

type Lock struct {
	gorm.Model
	LockID  string `gorm:"type:varchar(15);not null;unique"`
	Message string
}

func RandomStr() string {
	var chars = []rune("abcdefghjkmnopqrstvwxyz01234567890")
	id := make([]rune, 15)
	for i := range id {
		id[i] = chars[rand.Intn(len(chars))]
	}
	return string(id)

}

func Gen() (id string) {
	lockid := RandomStr()
	return lockid
}

func GetLock(id string) (data Lock, err error) {
	var lock Lock

	err = db.Where(&Lock{LockID: id}).First(&lock).Error
	return lock, err

}

func AddLock(message string) (id string) {
	lockid := Gen()
	db.Create(&Lock{LockID: lockid, Message: message})
	return lockid
}

func DeleteLock(id string) {
	var lock Lock
	db.Where(&Lock{LockID: id}).Delete(&lock)
}
