package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["beego-api-firebase-template/controllers/api:UserController"] = append(beego.GlobalControllerRouter["beego-api-firebase-template/controllers/api:UserController"],
        beego.ControllerComments{
            Method: "CreateUser",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["beego-api-firebase-template/controllers/api:UserController"] = append(beego.GlobalControllerRouter["beego-api-firebase-template/controllers/api:UserController"],
        beego.ControllerComments{
            Method: "GetUser",
            Router: "/:firebaseUid",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
