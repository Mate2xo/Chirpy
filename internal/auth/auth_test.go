package auth_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Mate2xo/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	id := uuid.New()
	validTokenString, _ := auth.MakeJWT(id, "sekret")
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

func TestGetBearerToken(t *testing.T) {
	validTokenString, _ := auth.MakeJWT(uuid.New(), "sekret")
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name:    "retrieves the 'Bearer' token from the 'Authorization' header",
			headers: http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", validTokenString)}},
			want:    validTokenString,
			wantErr: false,
		},
		{
			name:    "errors when there is no Authorization header",
			headers: http.Header{"Chicken": []string{"Nanban"}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "errors when the Authorization header does not contain a 'Bearer' token",
			headers: http.Header{"Authorization": []string{"What is Love"}},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.GetBearerToken(tt.headers)
			if (gotErr != nil) != tt.wantErr {
				t.Fatalf("GetBearerToken() failed: got error %v (expecting error?: %v)", gotErr, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
