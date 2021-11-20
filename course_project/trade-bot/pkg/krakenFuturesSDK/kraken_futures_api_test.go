package krakenFuturesSDK

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var apiClient *API
var publicKey string
var privateKey string

func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to load enviroment variables: ", err)
	}

	publicKey = os.Getenv("PUBLIC_API_KEY")
	if publicKey == "" {
		log.Fatal("empty public key")
	}
	privateKey = os.Getenv("PRIVATE_API_KEY")
	if privateKey == "" {
		log.Fatal("empty private key")
	}
	apiURL := os.Getenv("KRAKEN_API_URL")
	if apiURL == "" {
		log.Fatal("empty api url")
	}

	apiClient = NewAPI(publicKey, privateKey, apiURL)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestNewAPI(t *testing.T) {
	tests := []struct {
		name string
		want *API
	}{
		{name: "default", want: apiClient},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewAPI(publicKey, privateKey, apiClient.apiURL))
		})
	}
}

func TestNewAPIWithClient(t *testing.T) {
	type args struct {
		client *http.Client
	}
	tests := []struct {
		name string
		args args
		want *API
	}{
		{name: "default", args: args{client: http.DefaultClient}, want: apiClient},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewAPIWithClient(publicKey, privateKey, apiClient.apiURL, tt.args.client))
		})
	}
}

func TestAPI_(t *testing.T) {

}

func TestAPI_createSignature(t *testing.T) {
	api := NewAPI("h9GBs5aZb2ec2dpMR0g6gz8ih3hq0c+GW4LxZtLCgVYD0wdcVb+S5+vP",
		"xQvgl2eOV/nYzd0ok0KN7S1S4Yv5GRAC7k7HIp1WUE4ypYYWKf/xm3fHyXc3/asGgEzY7vYTkg6EpIijr4AU8g1s",
		apiClient.apiURL)
	type args struct {
		endPoint string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "default",
			args: args{
				endPoint: "/api/v3/cancelallorders",
			}, want: "bDmgmXqpb1rwGAyH6Nar3QXsCj+eeU1FMvkB1thfeXrz1FwLxizt1MbDg7d4BvZa+8yBmD5iE/0PLVEkODLxmA==",
			wantErr: false},
		{name: "with prefix",
			args: args{
				endPoint: "/derivatives/api/v3/cancelallorders",
			}, want: "bDmgmXqpb1rwGAyH6Nar3QXsCj+eeU1FMvkB1thfeXrz1FwLxizt1MbDg7d4BvZa+8yBmD5iE/0PLVEkODLxmA==",
			wantErr: false},
		{name: "with prefix and post data",
			args: args{
				endPoint: "/derivatives/api/v3/cancelallorders",
			}, want: "bDmgmXqpb1rwGAyH6Nar3QXsCj+eeU1FMvkB1thfeXrz1FwLxizt1MbDg7d4BvZa+8yBmD5iE/0PLVEkODLxmA==",
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := api.createSignature(tt.args.endPoint, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("createSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createSignature() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_FeeSchedules(t *testing.T) {
	_, err := apiClient.FeeSchedules()
	if err != nil {
		t.Errorf("FeeSchedules() error = %v", err)
	}
}

func TestAPI_OrderBook(t *testing.T) {
	type args struct {
		symbol string
	}
	tests := []struct {
		name           string
		args           args
		expectedResult string
		wantErr        bool
	}{
		{name: "default", args: args{symbol: "pi_xbtusd"}, expectedResult: "success", wantErr: false},
		{name: "with invalid symbol", args: args{symbol: "skldfj"}, expectedResult: "error", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := apiClient.OrderBook(tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.expectedResult, got.Result)
		})
	}
}

func TestAPI_Tickers(t *testing.T) {
	_, err := apiClient.Tickers()
	if err != nil {
		t.Errorf("Tickers() error = %v", err)
	}
}

func TestAPI_Instruments(t *testing.T) {
	_, err := apiClient.Instruments()
	if err != nil {
		t.Errorf("Instruments() error = %v", err)
	}
}

func TestAPI_CancelAllOrders(t *testing.T) {
	type args struct {
		symbol string
	}
	tests := []struct {
		name           string
		args           args
		expectedResult string
		wantErr        bool
	}{
		{name: "without symbol", args: args{symbol: ""}, expectedResult: "success", wantErr: false},
		{name: "with symbol", args: args{symbol: "PI_XBTUSD"}, expectedResult: "success", wantErr: false},
		{name: "with incorrect symbol", args: args{symbol: "skdjfakj"}, expectedResult: "error", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := apiClient.CancelAllOrders(tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("CancelAllOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.expectedResult, got.Result)
		})
	}
}

func TestAPI_SendOrder(t *testing.T) {
	sendSTPArguments := SendOrderArguments{
		OrderType:     "stp",
		Symbol:        "pi_xbtusd",
		Side:          "buy",
		Size:          1000,
		LimitPrice:    3900,
		StopPrice:     10000,
		TriggerSignal: "mark",
		CliOrderID:    "",
		ReduceOnly:    true,
	}
	sendWithInvalidSymbol := SendOrderArguments{
		OrderType:     "stp",
		Symbol:        "askldfl;aksdjf",
		Side:          "buy",
		Size:          1000,
		LimitPrice:    3900,
		StopPrice:     10000,
		TriggerSignal: "mark",
		CliOrderID:    "",
		ReduceOnly:    true,
	}
	type args struct {
		args SendOrderArguments
	}
	tests := []struct {
		name           string
		args           args
		expectedResult string
		wantErr        bool
	}{
		{name: "send stp order", args: args{args: sendSTPArguments}, expectedResult: "success", wantErr: false},
		{name: "send stp order with invalid symbol", args: args{args: sendWithInvalidSymbol}, expectedResult: "error", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := apiClient.SendOrder(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.expectedResult, got.Result)
		})
	}
}

func TestAPI_EditOrder(t *testing.T) {
	sendOrderResponse, err := apiClient.SendOrder(SendOrderArguments{
		OrderType:     "stp",
		Symbol:        "PI_XBTUSD",
		Side:          "buy",
		Size:          100,
		LimitPrice:    9400,
		StopPrice:     10000,
		TriggerSignal: "mark",
		CliOrderID:    "",
		ReduceOnly:    false,
	})
	if err != nil {
		t.Errorf("unable to send order")
	}

	type args struct {
		args EditOrderArguments
	}
	tests := []struct {
		name           string
		args           args
		expectedResult string
		wantErr        bool
	}{
		{name: "default", args: args{args: EditOrderArguments{
			OrderID:    sendOrderResponse.SendStatus.OrderID,
			Size:       1000,
			LimitPrice: 9300,
			StopPrice:  11000,
			CliOrdID:   "",
		}}, expectedResult: "success", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := apiClient.EditOrder(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EditOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.expectedResult, got.Result)
		})
	}
}

func TestAPI_CancelOrder(t *testing.T) {
	sendOrderResponse, err := apiClient.SendOrder(SendOrderArguments{
		OrderType:     "stp",
		Symbol:        "PI_XBTUSD",
		Side:          "buy",
		Size:          100,
		LimitPrice:    9400,
		StopPrice:     10000,
		TriggerSignal: "mark",
		CliOrderID:    "",
		ReduceOnly:    false,
	})
	if err != nil {
		t.Errorf("unable to send order")
	}

	type args struct {
		args CancelOrderArguments
	}
	tests := []struct {
		name           string
		args           args
		expectedResult string
		wantErr        bool
	}{
		{name: "default", args: args{args: CancelOrderArguments{
			OrderID:  sendOrderResponse.SendStatus.OrderID,
			CliOrdID: sendOrderResponse.SendStatus.CliOrderID,
		}}, expectedResult: "success", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := apiClient.CancelOrder(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("CancelOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.expectedResult, got.Result)
		})
	}
}

func TestAPI_doRequest(t *testing.T) {
	type args struct {
		reqType string
		reqURL  string
		headers map[string]string
		typ     interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "default request", args: args{
			reqType: http.MethodGet,
			reqURL:  fmt.Sprintf("%s%s", apiClient.apiURL, "/derivatives/api/v3/feeschedules"),
			typ:     &FeeSchedulesResponse{},
		}, wantErr: false},
		{name: "invalid url request", args: args{
			reqType: http.MethodGet,
			reqURL:  fmt.Sprintf("%s%s", apiClient.apiURL, "//////^&%^&%////afasdf%^&$%/"),
			typ:     &FeeSchedulesResponse{},
		}, wantErr: true},
		{name: "nil typ", wantErr: true},
		{name: "invalid typ type", args: args{
			reqType: http.MethodGet,
			reqURL:  fmt.Sprintf("%s%s", apiClient.apiURL, "/derivatives/api/v3/feeschedules"),
			typ:     map[int]int{},
		}, wantErr: true},
		{name: "invalid request type", args: args{
			reqType: "jkasdfjkh",
			reqURL:  fmt.Sprintf("%s%s", apiClient.apiURL, "/derivatives/api/v3/feeschedules"),
			typ:     &FeeSchedulesResponse{},
		}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := apiClient.doRequest(tt.args.reqType, tt.args.reqURL, tt.args.headers, tt.args.typ)
			if (err != nil) != tt.wantErr {
				t.Errorf("doRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
