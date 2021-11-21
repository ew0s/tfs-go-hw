package krakenFuturesWSSDK

import (
	"context"
	"os"
	"testing"
	"time"
	"trade-bot/pkg/krakenFuturesWSSDK/wsConfigs"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var apiClient *WSAPI
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

	config := wsConfigs.WSAPIConfig{
		Requests: wsConfigs.WSAPIRequestsConfig{
			WriteWait:      int(10 * time.Second),
			PongWait:       int(10 * time.Second),
			PingPeriod:     int(10 * time.Second),
			MaxMessageSize: 512,
		},
		Kraken: wsConfigs.KrakenWSAPIConfiguration{
			WSAPIURL:   apiURL,
			PublicKey:  publicKey,
			PrivateKey: privateKey,
		},
	}

	apiClient = NewWSAPI(config)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestWSAPI_CandlesTrade(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		feed       string
		productIDs []string
		wantErr    bool
	}{
		{name: "bad feed", feed: "asdflkjfasdlkj", productIDs: []string{"PI_XBTUSD"}, wantErr: true},
		{name: "bad products ids", feed: "candles_trade_1m", productIDs: []string{"sdlf;jk"}, wantErr: true},
		{name: "default", feed: "candles_trade_1m", productIDs: []string{"PI_XBTUSD"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			candlesTradesCh, err := apiClient.CandlesTrade(ctx, tt.feed, tt.productIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("CandlesTrade() error = %v, wantErr %v", err, tt.wantErr)
				cancel()
				return
			}

			if err == nil {
				go func() {
					time.Sleep(time.Second * 70)
					cancel()
				}()

				shapshotFeed := tt.feed + "_snapshot"
				for val := range candlesTradesCh {
					if val.Feed != tt.feed && val.Feed != shapshotFeed {
						t.Errorf("CandlesTrade() got = %v, want %v or %v", val.Feed, tt.feed, shapshotFeed)
					}
				}
			}
		})
	}
}

func TestWSAPI_Heartbeat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "default", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			dataCh, err := apiClient.Heartbeat(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Heartbeat() error = %v, wantErr %v", err, tt.wantErr)
				cancel()
				return
			}

			go func() {
				time.Sleep(time.Second * 20)
				cancel()
			}()

			for val := range dataCh {
				if val.Feed != "heartbeat" {
					t.Errorf("Heartbeat() got = %v, want %v", val, "heartbeat")
				}
			}
		})
	}
}

func TestWSAPI_serveWSHeartbeat(t *testing.T) {
	t.Parallel()
	type args struct {
		args KrakenSendMessageArguments
		typ  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "default", args: args{
			args: KrakenSendMessageArguments{
				Event: "subscribe",
				Feed:  "heartbeat",
			},
			typ: &HeartbeatSubscriptionData{},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			dataCh, errCh, err := apiClient.serveWS(ctx, tt.args.args, tt.args.typ)
			if (err != nil) != tt.wantErr {
				t.Errorf("serveWS() error = %v, wantErr %v", err, tt.wantErr)
				cancel()
				return
			}

			go func() {
				log.Warn(<-errCh)
			}()

			go func() {
				time.Sleep(time.Second * 20)
				cancel()
			}()

			for val := range dataCh {
				assertedVal, ok := val.(*HeartbeatSubscriptionData)
				if !ok {
					t.Errorf("could not assert to heartbeat")
				}
				if assertedVal.Feed != "heartbeat" {
					t.Errorf("Heartbeat() got = %v, want %v", assertedVal.Feed, "heartbeat")
				}
			}
		})
	}
}
