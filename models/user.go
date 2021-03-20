package models

import (
	"beego-api-firebase-template/modules/global"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type UserType int

const (
	Admin UserType = iota
	General
)

type SignUpMode int

const (
	SignUpWithEmail SignUpMode = iota
)

type EUser struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Token      string `json:"token"`
	UID        string `json:"uid"`
	SignUpMode string `json:"signUpMode"` // move to enums later
}

type User struct {
	Id          int64      `json:"id"`
	Name        string     `json:"name"`
	ApiKey      string     `json:"apiKey"`
	Type        UserType   `json:"type"`
	Subscribed  bool       `json:"subscribed" orm:"default(false)"`
	Verified    bool       `json:"verified"`
	Username    string     `json:"username" orm:"unique"`
	Email       string     `json:"email" orm:"unique"`
	SignUpMode  SignUpMode `json:"signUpMode"`
	FirebaseUID string     `json:"firebaseUID" orm:"unique"`
	Created     int64      `json:"created"`
	Updated     int64      `json:"updated"`
}

func (u *User) Insert() error {

	u.Created = time.Now().UnixNano()
	u.Updated = u.Created

	if _, err := orm.NewOrm().Insert(u); err != nil {
		return err
	}
	return nil
}

func (u *User) Update(fields ...string) error {

	if len(fields) == 0 {
		u.Updated = time.Now().UnixNano()
	} else {
		fields = append(fields, "Updated")
		u.Updated = time.Now().UnixNano()
	}
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil

}

func (u *User) Delete() error {
	if _, err := orm.NewOrm().Delete(u); err != nil {
		return err
	}
	return nil
}

func AllUsers() orm.QuerySeter {
	return orm.NewOrm().QueryTable("user")
}

func GetUserFromID(id int64) (*User, error) {

	u := new(User)

	if user, found := global.UserCache.Get(id); found == false || user == nil {
		fmt.Println("user db hit")
		if err := AllUsers().Filter("id", id).RelatedSel().One(u); err != nil {
			fmt.Println("error while calling db to get user from user id", err.Error())
			return nil, err
		} else {
			fmt.Println("got from db", u)
			if global.UserCache.Set(id, *u, 1) == false {
				fmt.Println("setting user cache failed")
			} else {
				fmt.Println("setting user cache success")
			}
		}
	} else {
		fmt.Println("user cache hit")
		xu := user.(User)
		u = &xu
		fmt.Println(u)
	}

	return u, nil
}

func GetUserFromAPIKey(key string) (*User, error) {

	u := new(User)

	if user, found := global.UserApiCache.Get(key); found == false || user == nil {
		fmt.Println("user api db hit")
		if err := AllUsers().Filter("api_key", key).RelatedSel().One(u); err != nil {
			fmt.Println("error while calling db to get user from user api key", err.Error())
			return nil, err
		} else {
			fmt.Println("got from db", u)
			if global.UserApiCache.Set(key, *u, 1) == false {
				fmt.Println("setting user api cache failed")
			} else {
				fmt.Println("setting user api cache success")
			}
		}
	} else {
		fmt.Println("user api cache hit")
		xu := user.(User)
		u = &xu
		fmt.Println(u)
	}

	return u, nil
}

func init() {
	orm.RegisterModel(new(User))
}
