package helper

import (
	"log"
	"github.com/unidoc/unidoc/common/license"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

const (
	CompanyName = "Digital Alpha Indonesia"
	Licen       = `
-----BEGIN UNIDOC LICENSE KEY-----
eyJsaWNlbnNlX2lkIjoiZWM2YjNkMjItYjNkOS00NjkyLTc3ODQtYTMzOTkyNzUyYjhhIiwiY3VzdG9tZXJfaWQiOiIzMmI0NzQ3Ni0yMTI2LTQ1NGYtN2E2MS0wMjFiZDMwZTc0YTQiLCJjdXN0b21lcl9uYW1lIjoiRGlnaXRhbCBBbHBoYSBJbmRvbmVzaWEiLCJjdXN0b21lcl9lbWFpbCI6ImFuZHJlLnByYXRhbWFAdWFuZ3RlbWFuLmNvbSIsInRpZXIiOiJidXNpbmVzcyIsImNyZWF0ZWRfYXQiOjE1NzUwMTIyOTIsImV4cGlyZXNfYXQiOjE2MDY2OTQzOTksImNyZWF0b3JfbmFtZSI6IlVuaURvYyBTdXBwb3J0IiwiY3JlYXRvcl9lbWFpbCI6InN1cHBvcnRAdW5pZG9jLmlvIiwidW5pcGRmIjp0cnVlLCJ1bmlvZmZpY2UiOnRydWUsInRyaWFsIjpmYWxzZX0=
+
VaLlarl3MB1oaoqp9OukhnBaZi91Kse1JJ27bVdjuN6CrpurccwNx5CU8kcLAxFNqKzlaiBJmbHkYvl2T7HGLYGnErrg2DhF2e2JfFnb7lc78uouEdjS5raf9dmjZsjK9FNDKGRUQxYhRPYU22/GKGrrHcaJH5FCsaBLRVBOMWIUraFktEcwzaLILneZ/Dgsdv+ymIHowHDdViM4wwdkWb7fIonuy/QsDgBygBUkmPvZzjFQ+HRnFnqTltP/DZwez9cOemVS5IUl4/80Y1kVnD+63K8kSBrr8lyBYQhJR0DTez0Z08s2KJeWCRwLJdnvR5EJzvkld0Q0NmcS+9x03g==
-----END UNIDOC LICENSE KEY-----
	`
)

func init() {
	LoadUnidoc()
}

func LoadUnidoc() {
	err := license.SetLicenseKey(Licen)

	if err != nil {
		log.Printf("Error loading license: %v\n", err)
	} else {
		log.Println("Success Loading unidocs")
	}
}

func LoadFontSansPro() *pdf.PdfFont {
	fontSansPro, _ := pdf.NewPdfFontFromTTFFile("files/font/sanspro.ttf")
	return fontSansPro
}

func LoadFontSansProBold() *pdf.PdfFont {
	fontSansProBold, _ := pdf.NewPdfFontFromTTFFile("files/font/sansprobold.ttf")
	return fontSansProBold
}
