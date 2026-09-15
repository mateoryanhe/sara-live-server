package recharge

import "testing"

func TestHaiPayCheckCollectPayAmount(t *testing.T) {
	tests := []struct {
		name           string
		queryData      *haiPayCollectQueryData
		orderCurrency  string
		orderPayAmount float64
		wantField      string
		wantErr        bool
	}{
		{
			name:           "same currency uses amount",
			queryData:      &haiPayCollectQueryData{Currency: "USD", Amount: "20.00"},
			orderCurrency:  "USD",
			orderPayAmount: 20,
			wantField:      "amount",
		},
		{
			name: "converted payment uses original amount",
			queryData: &haiPayCollectQueryData{
				Currency: "HKD", Amount: "316.82", OriginalCurrency: "USD", OriginalAmount: "39.99",
			},
			orderCurrency:  "USD",
			orderPayAmount: 39.99,
			wantField:      "originalAmount",
		},
		{
			name:           "original amount mismatch",
			queryData:      &haiPayCollectQueryData{Currency: "HKD", Amount: "316.82", OriginalCurrency: "USD", OriginalAmount: "40.00"},
			orderCurrency:  "USD",
			orderPayAmount: 39.99,
			wantErr:        true,
		},
		{
			name:           "query has no matching currency",
			queryData:      &haiPayCollectQueryData{Currency: "HKD", Amount: "316.82"},
			orderCurrency:  "USD",
			orderPayAmount: 39.99,
			wantErr:        true,
		},
		{
			name:           "query amount beyond cents",
			queryData:      &haiPayCollectQueryData{Currency: "USD", Amount: "20.001"},
			orderCurrency:  "USD",
			orderPayAmount: 20,
			wantErr:        true,
		},
		{
			name:           "order amount beyond cents",
			queryData:      &haiPayCollectQueryData{Currency: "USD", Amount: "20.00"},
			orderCurrency:  "USD",
			orderPayAmount: 20.001,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, _, err := haiPayCheckCollectPayAmount(tt.queryData, tt.orderCurrency, tt.orderPayAmount)
			if (err != nil) != tt.wantErr {
				t.Fatalf("haiPayCheckCollectPayAmount() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && field != tt.wantField {
				t.Fatalf("haiPayCheckCollectPayAmount() field = %q, want %q", field, tt.wantField)
			}
		})
	}
}
