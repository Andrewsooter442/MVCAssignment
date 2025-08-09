package handler

import (
	"net/url"
	"testing"
)

func TestValidatePaymentForm(t *testing.T) {
	testCases := []struct {
		name        string
		inputForm   url.Values
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid Data",
			inputForm: url.Values{
				"paymentMethod": {"card"},
				"orderId":       {"123"},
				"total":         {"5000"},
			},
			expectError: false,
		},
		{
			name: "Missing paymentMethod",
			inputForm: url.Values{
				"orderId": {"123"},
				"total":   {"5000"},
			},
			expectError: true,
			errorMsg:    "missing required field: paymentMethod",
		},
		{
			name: "Missing orderId",
			inputForm: url.Values{
				"paymentMethod": {"card"},
				"total":         {"5000"},
			},
			expectError: true,
			errorMsg:    "missing required field: orderId",
		},
		{
			name: "Invalid total (not a number)",
			inputForm: url.Values{
				"paymentMethod": {"card"},
				"orderId":       {"123"},
				"total":         {"abc"},
			},
			expectError: true,
			errorMsg:    "invalid format for total, must be an integer",
		},
		{
			name: "Invalid total (zero)",
			inputForm: url.Values{
				"paymentMethod": {"card"},
				"orderId":       {"123"},
				"total":         {"0"},
			},
			expectError: true,
			errorMsg:    "total must be a positive number",
		},
		{
			name: "Unsupported payment method",
			inputForm: url.Values{
				"paymentMethod": {"upi"},
				"orderId":       {"123"},
				"total":         {"100"},
			},
			expectError: true,
			errorMsg:    "payment method 'bitcoin' is not supported",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePaymentData(tc.inputForm)

			if !tc.expectError && err != nil {
				t.Fatalf("expected no error, but got: %v", err)
			}

			if tc.expectError && err == nil {
				t.Fatalf("expected an error, but got nil")
			}

		})
	}

}

func TestValidateOrderForm(t *testing.T) {
	testCases := []struct {
		name        string
		inputForm   url.Values
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid Order Form",
			inputForm: url.Values{
				"itemId":      {"101", "102"},
				"quantity":    {"2", "1"},
				"instruction": {"extra cheese", "no onions"},
				"tableNumber": {"15"},
			},
			expectError: false,
		},
		{
			name: "No Items in Order",
			inputForm: url.Values{
				"tableNumber": {"15"},
			},
			expectError: true,
			errorMsg:    "order must contain at least one item",
		},
		{
			name: "Mismatched Items and Quantities",
			inputForm: url.Values{
				"itemId":      {"101", "102"},
				"quantity":    {"2"},
				"instruction": {"extra cheese", "no onions"},
				"tableNumber": {"15"},
			},
			expectError: true,
			errorMsg:    "mismatched number of order items, quantities, and instructions",
		},
		{
			name: "Missing Table Number",
			inputForm: url.Values{
				"itemId":      {"101"},
				"quantity":    {"1"},
				"instruction": {"none"},
			},
			expectError: true,
			errorMsg:    "tableNumber is a required field",
		},
		{
			name: "Invalid Table Number (Not a Number)",
			inputForm: url.Values{
				"itemId":      {"101"},
				"quantity":    {"1"},
				"instruction": {"none"},
				"tableNumber": {"abc"},
			},
			expectError: true,
			errorMsg:    "invalid table number",
		},
		{
			name: "Invalid Table Number (Zero)",
			inputForm: url.Values{
				"itemId":      {"101"},
				"quantity":    {"1"},
				"instruction": {"none"},
				"tableNumber": {"0"},
			},
			expectError: true,
			errorMsg:    "invalid table number",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateOrderForm(tc.inputForm)

			if !tc.expectError && err != nil {
				t.Fatalf("expected no error, but got: %v", err)
			}

			if tc.expectError && err == nil {
				t.Fatalf("expected an error, but got nil")
			}

		})
	}
}
