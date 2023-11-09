package access_control

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"testing"
	"time"
)

func TestAccessControl_CheckAccess(t *testing.T) {
	rsaPub := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAh7IigpveHCxCT0McMT7fgonmUoDTSAR9E/L+/A8fXQiuJETUhJ01iMAYA848EEx8BCPb0Fke+WkCOFPvLn4VtJm5k44Zdnq0cVA0St/gcm7OkGz/XaIOf0V9zl4riM6VaQ98V4VJIvDy75QTtVJktZsfgC5fsWMiNiGXU3YA8UmBSTqakEov65EFbSIXwAoLwa7Ql7G2uoHlr7xF81MGpDVX9WY5nV45UJQOx2hG2GO1P5ELO6/PI8XeXQM8cTxRE4OdBKZ8LUEvYhYEErOv7Eg7NgcgtzEcfhKaVr7jcLlGzxg3cNgiMBVFugonZyWFP0sKPe3s64Krc1vHhAWonQIDAQAB"
	rsaCert := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", rsaPub)
	type fields struct {
		jwtRSACert string
		roles      []*Role
	}
	type args struct {
		tokenString string
		action      *Action
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				jwtRSACert: rsaCert,
				roles: []*Role{
					{
						name: "main",
						privileges: &Privileges{
							"main-page": &Privilege{
								access:     Accessed,
								privileges: &Privileges{},
							},
						},
					},
				},
			},
			args: args{
				tokenString: "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI3YjRYbjBTVU1mdUlkWGlxemlCbmd0czBXemRwWlRiZnl5VFdmX1JCakdRIn0.eyJleHAiOjE2NzMyNTQ3MzMsImlhdCI6MTY3MzI1NDQzMywianRpIjoiNmY5YmZhZmYtYjI0My00MDYwLWJhNjgtNmE3MmYzODllOTM3IiwiaXNzIjoiaHR0cHM6Ly9rZXljbG9hay5mYXJtcGVyc3Bla3RpdmEuYWVyb2lkZWEucnUvYXV0aC9yZWFsbXMvVGVzdCIsInN1YiI6ImM1Njg3OTMzLTQ1YzEtNGI3ZS1iYmNkLWRkYTQwY2FmYWE1MyIsInR5cCI6IkJlYXJlciIsImF6cCI6ImZvY3VzIiwic2Vzc2lvbl9zdGF0ZSI6ImI2YWU1YzllLWVlNTktNGFlZi1hM2I1LTE3ZTViMzgyZmZjMiIsImFjciI6IjEiLCJzY29wZSI6IioiLCJzaWQiOiJiNmFlNWM5ZS1lZTU5LTRhZWYtYTNiNS0xN2U1YjM4MmZmYzIifQ.SLItlr9SP6KkFTduzlhilN01qR-t_vX7o1WaKKttMzshtomJ-MYW-zWe6EOPpTUPAy8BzI1ZH2Rm2edW8ElUfxq2P3yGW2w-rF7zymtT_ZecQBlqAHvzsbucf-YI_mziSvYTgrzeP7ZifeynzsZYsmYliwLPhlvEZbW-O8ALg-dG2UmJDymYA_kdOXzrWhUusv3hQOzB6_IVoFTLeHd0qRjQ9QTShTUkB_3X1mxHPSQUUVi7uKuiAguU3h-y_-I1xNh-vM48yGWLMkCqqHv7dHqFbAXC6-5e1z-z59CDJECevbRu57AXehAn2mypvtbK34wRw1WcYTF8JSzKU3tfrQ",
				action: &Action{
					Path:   "/settings",
					Method: "GET",
				},
			},
			wantErr: false,
		},
		{
			name: "wrong token",
			fields: fields{
				jwtRSACert: "",
				roles: []*Role{
					{
						name: "main",
						privileges: &Privileges{
							"main-page": &Privilege{
								access:     Accessed,
								privileges: &Privileges{},
							},
						},
					},
				},
			},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJyb2xlIjoibWFpbiJ9.gOWKpBz9dP6SEvfCrsYyPa3tkmCo5hNItKC00-cvn0Q",
				action: &Action{
					Path:   "/main-page",
					Method: "GET",
				},
			},
			wantErr: true,
		},
		{
			name: "expired",
			fields: fields{
				jwtRSACert: "",
				roles: []*Role{
					{
						name: "main",
						privileges: &Privileges{
							"main-page": &Privilege{
								access:     Accessed,
								privileges: &Privileges{},
							},
						},
					},
				},
			},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZXhwIjoxNTE2MjM5MDIyLCJpYXQiOjE1MTYyMzkwMjIsInJvbGUiOiJtYWluIn0.HIwhUzZ1xNMgNtwLNWZYpdf6XySmIARanXzwotYiTU8",
				action: &Action{
					Path:   "/main-page",
					Method: "GET",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwt.NewWithClaims(jwt.SigningMethodRS256, AccessClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
				Scope: "*",
			})
			jwtRsaCert, err := jwt.ParseRSAPublicKeyFromPEM([]byte(tt.fields.jwtRSACert))
			if err != nil {
				t.Errorf("Wrong rsa cert error = %v", err)
				return
			}

			s := AccessControl{
				jwtRSACert: jwtRsaCert,
			}
			if err := s.CheckAccess(tt.args.tokenString, tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("CheckAccess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
