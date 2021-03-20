package api

import (
	"beego-api-firebase-template/models"
	"beego-api-firebase-template/modules/utils"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	beego "github.com/beego/beego/v2/server/web"
	"net/http"
	"strings"
)

type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description Create New User
// @Param body body models.EUser true "User details"
// @Param X-Token header string true "Firebase token; consumed by middleware"
// @Success 200 {object} models.User "User created successfully"
// @Failure 400 Bad Request
// @Failure 401 Not Authorized
// @Failure 422 Request body is not in proper format
// @Failure 500 Internal Server Error
// @router / [post]
func (u *UserController) CreateUser() {

	var eu models.EUser
	var us models.User

	if err := json.Unmarshal(u.Ctx.Input.RequestBody, &eu); err == nil {

		eu.Email = strings.TrimSpace(eu.Email)
		if !govalidator.IsEmail(eu.Email) {
			u.Ctx.WriteString("invalid email address")
			u.Ctx.Output.SetStatus(http.StatusBadRequest)
			return
		}

		var mode models.SignUpMode
		switch eu.SignUpMode {
		case "signUpWithEmail":
			mode = models.SignUpWithEmail
		default:
			u.Ctx.WriteString("wrong signup mode")
			u.Ctx.Output.SetStatus(http.StatusBadRequest)
			return
		}

		us.Email = eu.Email
		us.ApiKey = utils.GenerateApiKey()
		us.SignUpMode = mode
		us.FirebaseUID = eu.UID
		us.Type = models.General
		// when we add a cache for usernames
		// right now, not returning custom error for unique for usernames
		us.Username = strings.TrimSpace(eu.Username)
		us.Name = ""

		if err := (&us).Insert(); err == nil {
			u.Ctx.Output.SetStatus(http.StatusOK)
			_ = u.Ctx.Output.JSON(us, true, false)
			return
		} else {
			u.Ctx.WriteString(err.Error())
			u.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}

	} else {
		u.Ctx.Output.SetStatus(http.StatusUnprocessableEntity)
		return
	}
}

// @Title GetUser
// @Description Get Existing User
// @Param X-Token header string true "Firebase token"
// @Param firebaseUid path string true "the firebase uid of the user you need to fetch"
// @Success 200 {object} models.User "User fetched successfully"
// @Failure 401 Not Authorized
// @Failure 422 Request body is not in proper format
// @Failure 500 Internal Server Error
// @router /:firebaseUid [get]
func (u *UserController) GetUser() {

	var us models.User
	var verified bool

	xv := u.Ctx.Input.GetData("verified")
	if xv != nil {
		verified = xv.(bool)
	}

	fuid := u.Ctx.Input.Param(":firebaseUid")
	if fuid != "" {
		if err := models.AllUsers().Filter("firebase_u_i_d", fuid).RelatedSel().One(&us); err == nil {
			if !us.Verified && verified {
				us.Verified = true
				go func() {
					if err := us.Update("verified"); err == nil {
						fmt.Println("user email verified updated", us.Email)
					}
				}()
			}
			u.Ctx.Output.SetStatus(http.StatusOK)
			_ = u.Ctx.Output.JSON(us, true, false)
			return
		} else {
			u.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}
	} else {
		u.Ctx.Output.SetStatus(http.StatusNotFound)
		return
	}
}
