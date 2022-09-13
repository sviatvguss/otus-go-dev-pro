package hw09structvalidator

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		in              interface{}
		expectedValErrs []error
	}{
		{
			User{
				"0",
				"from Zero to Hero",
				20,
				"zero@hero.su",
				"admin",
				[]string{"79130000000", "7926", "75002222222222222222"},
				[]byte{},
			},
			[]error{ErrStrLength, ErrStrLength, ErrStrLength},
		},
		{
			User{
				"123456789012345678901234567890123456", // len = 36
				"Bobby",
				11,
				"bobby@hill",
				"professor",
				[]string{"123456789AB"},
				[]byte{},
			},
			[]error{ErrIntMin, ErrStrRegexp, ErrStrIn},
		},
		{
			User{
				"123456789012345678901234567890123456", // len = 36
				"Timotheus",
				22,
				"timon@server.ru",
				"admin",
				[]string{"123456789AB"},
				[]byte{},
			},
			[]error{},
		},
		{
			App{"12345"},
			nil,
		},
		{
			App{"1234"},
			[]error{ErrStrLength},
		},
		{
			Token{
				[]byte{},
				[]byte("empty"),
				make([]byte, 7),
			},
			nil,
		},
		{
			Response{200, "no matter"},
			nil,
		},
		{
			Response{202, "no matter"},
			[]error{ErrIntIn},
		},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			test := test
			t.Parallel()
			var validationErrors ValidationErrors
			errs := Validate(test.in)
			if len(test.expectedValErrs) == 0 {
				require.NoError(t, errs)
				return
			}
			require.ErrorAs(t, errs, &validationErrors)
			for i, e := range validationErrors {
				require.ErrorIs(t, e.Err, test.expectedValErrs[i])
			}
		})
	}
}
