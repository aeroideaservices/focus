package access_control

import (
	"crypto/rsa"
	"fmt"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

type AccessClaims struct {
	jwt.RegisteredClaims
	Scope      string `json:"scope"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
}

type AccessControl struct {
	jwtRSACert *rsa.PublicKey
}

func NewAccessControl(jwtRSACert *rsa.PublicKey) *AccessControl {
	return &AccessControl{jwtRSACert: jwtRSACert}
}

// Action - действие пользователя
type Action struct {
	Path   string
	Method string
}

func (a Action) String() string {
	str := strings.ReplaceAll(strings.Trim(a.Path, "/"), "/", ".")
	return fmt.Sprintf("%s.%s", strings.Trim(str, "/"), strings.ToLower(a.Method))
}

func NewAction(route string, method string) *Action {
	return &Action{Path: route, Method: method}
}

// CheckAccess проверяет доступ пользователя к определенному действию
func (s AccessControl) CheckAccess(tokenString string, action *Action) error {
	claims, err := s.GetClaims(tokenString)
	if err != nil {
		return err
	}

	role := NewRole("", strings.Split(claims.Scope, " "))
	if role == nil || !role.HasAccess(action) {
		return errors.Forbidden.New("forbidden")
	}

	return nil
}

func (s AccessControl) GetClaims(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		switch token.Method.(type) {
		case *jwt.SigningMethodRSA:
			return s.jwtRSACert, nil
		default:
			return nil, errors.Unauthorized.New("unexpected signing method")
		}
	})
	if err != nil {
		return nil, errors.Unauthorized.Wrap(err, "authorization error")
	}

	return token.Claims.(*AccessClaims), nil
}
