package access_control

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/sarulabs/di/v2"
)

var Definitions = []di.Def{
	{
		Build: func(ctn di.Container) (interface{}, error) {
			accessControl := ctn.Get("focus.accessControl").(*AccessControl)
			newAction := ctn.Get("focus.access.createActionRule").(CreateActionRule)
			return NewAccessMiddleware(accessControl, newAction), nil
		},
		Name: "focus.accessMiddleware",
	},
	{
		Build: func(ctn di.Container) (interface{}, error) {
			jwtRSACert := ctn.Get("focus.access.jwtRSACert").(string)
			rsaPub, err := jwt.ParseRSAPublicKeyFromPEM([]byte(jwtRSACert))
			if err != nil {
				return nil, err
			}
			return NewAccessControl(rsaPub), nil
		},
		Name: "focus.accessControl",
	},
	{
		Build: func(ctn di.Container) (interface{}, error) {
			rolesFile := ctn.Get("focus.access.rolesFile").(string)
			return RolesFromFile(rolesFile), nil
		},
		Name: "focus.access.roles",
	},
}
