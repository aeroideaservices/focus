package main

import (
	"fmt"
	"github.com/aeroideaservices/focus/services/access_control"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"net/http/httptest"
)

func main() {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI3YjRYbjBTVU1mdUlkWGlxemlCbmd0czBXemRwWlRiZnl5VFdmX1JCakdRIn0.eyJleHAiOjE2Njg0OTMwODIsImlhdCI6MTY2ODQ5Mjc4MiwianRpIjoiMzVlYmQ5OTUtYjE0NS00MWQxLThlYmMtN2U5ZGYwYjM2Y2MxIiwiaXNzIjoiaHR0cHM6Ly9rZXljbG9hay5mYXJtcGVyc3Bla3RpdmEuYWVyb2lkZWEucnUvYXV0aC9yZWFsbXMvVGVzdCIsInN1YiI6IjAyMDU3Y2JlLTZlZWYtNDkxOC1iZDBmLWFiYzk0NzU1MTM3ZiIsInR5cCI6IkJlYXJlciIsImF6cCI6ImZvY3VzIiwic2Vzc2lvbl9zdGF0ZSI6Ijg1Y2ZhYjYwLWYwZmQtNGU5YS1iNGQxLWQ1N2Y1NmQ4ZmNiNyIsImFjciI6IjEiLCJyZXNvdXJjZV9hY2Nlc3MiOnsiZm9jdXMiOnsicm9sZXMiOlsiZm9jdXMgYWRtaW4iXX19LCJzY29wZSI6ImVtYWlsIHNlcnZpY2VzIHByaXZpbGVnZXMgcHJvZmlsZSBjb250ZW50IGNhdGFsb2cubW9kZWxzLnN0b3JlIiwic2lkIjoiODVjZmFiNjAtZjBmZC00ZTlhLWI0ZDEtZDU3ZjU2ZDhmY2I3IiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJmb2N1cy1hZG1pbiIsImVtYWlsIjoiZm9jdXNAYWVyb2lkZWEucnUifQ.IUKODAggB7InY_9LQkXDtd5ezXZj5vk3H6BFNawno_JWwzsrMWx3DeRn9TqytaWtRCQr7XswQq2hqnjY5mu6pzqdWBW_yDmkjbecaryiU9DPyY6tdzsuq6bPAw9onlafb4yyQep32sxz-RBDvw0NMKPTimlF0bhseIlIJLOPFh0DnfH2e26042l4sv3A2YszY2FFNSlXVtm7jjAPQpIKP627h-zgON9rinLAgqY6ry4s6rij-1Yory2zmQL0136jEZzM97sECT2Wh-X4M56AZjeMymVO1E6L5ZzvzkS3H_etmMwBW9THF8HDmmPepIIuK6D1NzgoijGi8KLUSqNIBA\n"

	rsaPub, _ := jwt.ParseRSAPublicKeyFromPEM([]byte("-----BEGIN PUBLIC KEY-----\n" +
		"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAh7IigpveHCxCT0McMT7fgonmUoDTSAR9E/L+/A8fXQiuJETUhJ01iMAYA848EEx8BCPb0Fke+WkCOFPvLn4VtJm5k44Zdnq0cVA0St/gcm7OkGz/XaIOf0V9zl4riM6VaQ98V4VJIvDy75QTtVJktZsfgC5fsWMiNiGXU3YA8UmBSTqakEov65EFbSIXwAoLwa7Ql7G2uoHlr7xF81MGpDVX9WY5nV45UJQOx2hG2GO1P5ELO6/PI8XeXQM8cTxRE4OdBKZ8LUEvYhYEErOv7Eg7NgcgtzEcfhKaVr7jcLlGzxg3cNgiMBVFugonZyWFP0sKPe3s64Krc1vHhAWonQIDAQAB" +
		"\n-----END PUBLIC KEY-----"))
	accessControl := access_control.NewAccessControl(rsaPub)

	req, _ := http.NewRequest("GET", "http://localhost/settings", nil)
	req.Header.Set("Authorization", token)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req

	jwt.GetAlgorithms()
	mw := access_control.NewAccessMiddleware("", accessControl)
	mw.CheckAccess(c)

	fmt.Println(c.Errors)
}
