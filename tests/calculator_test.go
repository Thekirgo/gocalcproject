package tests

import (
	"calculator-service/internal/calculator"
	"strings"
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool    //Ожидается ли ошибка
		errMsg  string
	}{
		{
			name:    "simple addition",
			input:   "2+2",
			want:    4,
			wantErr: false,
		},
		{
			name:    "complex expression",
			input:   "(2+2)*2",
			want:    8,
			wantErr: false,
		},
		{
			name:    "division",
			input:   "10/2",
			want:    5,
			wantErr: false,
		},
		{
			name:    "invalid character",
			input:   "2+a", //Буква в выражении
			wantErr: true,
			errMsg:  "invalid number",
		},
		{
			name:    "division by zero",
			input:   "1/0",	//Деление на ноль
			wantErr: true,
			errMsg:  "division by zero",
		},
		{
			name:    "mismatched parentheses",
			input:   "(1+2", //Нет скобки
			wantErr: true,
			errMsg:  "mismatched parentheses",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Calc(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Calc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Calc() error = %v, wantErr %v", err, tt.errMsg)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("Calc() = %v, want %v", got, tt.want)
			}
		})
	}
}
