package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// BindErrors mengubah error binding/validasi Gin menjadi daftar pesan yang ramah dibaca.
func BindErrors(err error) []string {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		out := make([]string, 0, len(verrs))
		for _, fe := range verrs {
			out = append(out, fmt.Sprintf("field '%s' failed validation '%s'", fe.Field(), fe.Tag()))
		}
		return out
	}
	return []string{"invalid request body"}
}
