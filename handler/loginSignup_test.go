package handler

import (
	"github.com/Andrewsooter442/MVCAssignment/config"
	"testing"
)

func TestValidateLoginRequest(t *testing.T) {
	testCases := []struct {
		name        string
		input       config.LoginRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid Login",
			input: config.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			expectError: false,
		},
		{
			name: "Missing Username",
			input: config.LoginRequest{
				Password: "password123",
			},
			expectError: true,
			errorMsg:    "username is a required field",
		},
		{
			name: "Missing Password",
			input: config.LoginRequest{
				Username: "testuser",
			},
			expectError: true,
			errorMsg:    "password is a required field",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateLoginRequest(tc.input)

			if !tc.expectError && err != nil {
				t.Fatalf("expected no error, but got: %v", err)
			}
			if tc.expectError && err == nil {
				t.Fatalf("expected an error, but got nil")
			}
			if tc.expectError && err != nil && err.Error() != tc.errorMsg {
				t.Fatalf("expected error '%s', but got '%s'", tc.errorMsg, err.Error())
			}
		})
	}
}

func TestValidateSignupRequest(t *testing.T) {
	testCases := []struct {
		name        string
		input       config.SignupRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid Signup",
			input: config.SignupRequest{
				Username: "newuser",
				Password: "password123",
				Email:    "test@example.com",
				Phone:    "1234567890",
			},
			expectError: false,
		},
		{
			name: "Password Too Short",
			input: config.SignupRequest{
				Username: "newuser",
				Password: "123",
				Email:    "test@example.com",
				Phone:    "1234567890",
			},
			expectError: true,
			errorMsg:    "password must be at least 6 characters long",
		},
		{
			name: "Invalid Email (No @)",
			input: config.SignupRequest{
				Username: "newuser",
				Password: "password123",
				Email:    "testexample.com",
				Phone:    "1234567890",
			},
			expectError: true,
			errorMsg:    "invalid email format",
		},
		{
			name: "Invalid Phone (Too Short)",
			input: config.SignupRequest{
				Username: "newuser",
				Password: "password123",
				Email:    "test@example.com",
				Phone:    "123",
			},
			expectError: true,
			errorMsg:    "phone number seems too short",
		},
		{
			name: "Missing Username",
			input: config.SignupRequest{
				Password: "password123",
				Email:    "test@example.com",
				Phone:    "1234567890",
			},
			expectError: true,
			errorMsg:    "username is a required field",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateSignupRequest(tc.input)

			if !tc.expectError && err != nil {
				t.Fatalf("expected no error, but got: %v", err)
			}
			if tc.expectError && err == nil {
				t.Fatalf("expected an error, but got nil")
			}
			if tc.expectError && err != nil && err.Error() != tc.errorMsg {
				t.Fatalf("expected error '%s', but got '%s'", tc.errorMsg, err.Error())
			}
		})
	}
}
