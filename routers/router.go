// @APIVersion 1.0.0
// @Title Sample API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact mail@sampledomain.com
// @TermsOfServiceUrl http://sampledomain.com/tnc
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"beego-api-firebase-template/controllers/api"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/users",
			beego.NSInclude(
				&api.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
