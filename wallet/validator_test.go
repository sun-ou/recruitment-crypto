package wallet

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestIsPositiveDecimal(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("is_positive_decimal", IsPositiveDecimal)

	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"PositiveDecimal", args{"1.5"}, false},
		{"NegativeDecimal", args{"-1.5"}, true},
		{"PositiveInteger", args{"1"}, false},
		{"NegativeInteger", args{"-1"}, true},
		{"InvalidType", args{"not a number"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := struct {
				Value string `validate:"is_positive_decimal"`
			}{
				Value: tt.args.value.(string),
			}

			err := validate.Struct(v)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsPositiveDecimal(%v) = %v, wantErr %v", tt.args.value, err, tt.wantErr)
			}
		})
	}
}
