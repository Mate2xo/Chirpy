package auth_test

import (
	"testing"
	"time"

	"github.com/Mate2xo/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	id := uuid.New()
	validTokenString, _ := auth.MakeJWT(id, "sekret", time.Minute)
	validSecret := "sekret"

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tokenString string
		tokenSecret string
		want        uuid.UUID
		wantErr     bool
	}{
		{
			name:        "decodes and extract the User's uuid",
			tokenString: validTokenString,
			tokenSecret: validSecret,
			want:        id,
			wantErr:     false,
		},
		{
			name:        "errors with an incorrect tokenSecret",
			tokenString: validTokenString,
			tokenSecret: "walla",
			want:        uuid.UUID{},
			wantErr:     true,
		},
		{
			name:        "errors with an incorrect tokenString",
			tokenString: "walla",
			tokenSecret: validSecret,
			want:        uuid.UUID{},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() failed: got error %v (wanted an error?: %v)", gotErr, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ValidateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}
