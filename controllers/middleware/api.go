package middleware

import (
	"beego-api-firebase-template/models"
	"fmt"
	"net/http"
	"strings"

	"beego-api-firebase-template/modules/global"
	"github.com/beego/beego/v2/server/web/context"
	gcontext "golang.org/x/net/context"
)

func AuthFilter(ctx *context.Context) {

	path := ctx.Request.URL.Path

	// X-Token
	if strings.Contains(path, "/users") {

		xtoken := ctx.Request.Header.Get("X-Token")
		xctx := gcontext.Background()
		client, err := global.FirebaseApp.Auth(xctx)
		if err != nil {
			fmt.Printf("error getting auth client: %v\n\n", err)
		}

		token, err := client.VerifyIDToken(xctx, xtoken)
		if err != nil {
			ctx.Output.SetStatus(http.StatusUnauthorized)
			return
		}
		fmt.Println("verified x token for uid", token.UID)

		u, err := client.GetUser(xctx, token.UID)
		if err != nil {
			fmt.Printf("error getting user %s: %v\n", token.UID, err)
		} else {
			ctx.Input.SetData("verified", u.EmailVerified)
		}
	} else {
		// NOTE: use X-User-ApiKey -> for other endpoints to not hit firebase everytime
		// easier for testing endpoints that require auth

		key := strings.TrimSpace(ctx.Request.Header.Get("X-User-ApiKey"))
		if u, err := models.GetUserFromAPIKey(key); err == nil {
			ctx.Input.SetData("user", *u)
		}
	}

}
