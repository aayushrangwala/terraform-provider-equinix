package equinix

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/config"
)

func Test_waitUntilReservationProvisionable(t *testing.T) {
	type args struct {
		reservationId string
		instanceId    string
		handler       func(w http.ResponseWriter, r *http.Request)
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "error",
			args: args{
				reservationId: "reservationId",
				instanceId:    "instanceId",
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Content-Type", "application/json")
					w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
					w.WriteHeader(http.StatusInternalServerError)
				},
			},
			wantErr: true,
		},
		{
			name: "provisionable",
			args: args{
				reservationId: "reservationId",
				instanceId:    "instanceId",
				handler: (func() func(w http.ResponseWriter, r *http.Request) {
					invoked := new(int)

					responses := map[int]struct {
						id            string
						provisionable bool
					}{
						0: {"instanceId", false}, // should retry
						1: {"", true},            // should return success
					}

					return func(w http.ResponseWriter, r *http.Request) {
						response := responses[*invoked]
						*invoked++

						var device *metalv1.Device
						include := r.URL.Query().Get("include")
						if strings.Contains(include, "device") {
							device = &metalv1.Device{Id: &response.id}
						}
						reservation := metalv1.HardwareReservation{
							Device: device, Provisionable: &response.provisionable,
						}

						body, err := reservation.MarshalJSON()
						if err != nil {
							// This should never be reached and indicates a failure in the test itself
							panic(err)
						}

						w.Header().Add("Content-Type", "application/json")
						w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
						w.WriteHeader(http.StatusOK)
						w.Write(body)
					}
				})(),
			},
			wantErr: false,
		},
		{
			name: "reprovisioned",
			args: args{
				reservationId: "reservationId",
				instanceId:    "instanceId",
				handler: (func() func(w http.ResponseWriter, r *http.Request) {
					responses := map[int]struct {
						id            string
						provisionable bool
					}{
						0: {"instanceId", false},      // should retry
						1: {"new instance id", false}, // should return success
					}
					invoked := new(int)

					return func(w http.ResponseWriter, r *http.Request) {
						response := responses[*invoked]
						*invoked++

						var device *metalv1.Device
						include := r.URL.Query().Get("include")
						if strings.Contains(include, "device") {
							device = &metalv1.Device{Id: &response.id}
						}
						reservation := metalv1.HardwareReservation{
							Device: device, Provisionable: &response.provisionable,
						}

						body, err := reservation.MarshalJSON()
						if err != nil {
							// This should never be reached and indicates a failure in the test itself
							panic(err)
						}

						w.Header().Add("Content-Type", "application/json")
						w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
						w.WriteHeader(http.StatusOK)
						w.Write(body)
					}
				})(),
			},
			wantErr: false,
		},
		{
			name: "foreverDeprovisioning",
			args: args{
				reservationId: "reservationId",
				instanceId:    "instanceId",
				handler: func(w http.ResponseWriter, r *http.Request) {
					reservation := metalv1.HardwareReservation{
						Device: nil, Provisionable: metalv1.PtrBool(false),
					}

					body, err := reservation.MarshalJSON()
					if err != nil {
						// This should never be reached and indicates a failure in the test itself
						panic(err)
					}

					w.Header().Add("Content-Type", "application/json")
					w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
					w.WriteHeader(http.StatusOK)
					w.Write(body)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockAPI := httptest.NewServer(http.HandlerFunc(tt.args.handler))
			meta := &config.Config{
				BaseURL: mockAPI.URL,
				Token:   "fakeTokenForMock",
			}
			meta.Load(ctx)

			if err := waitUntilReservationProvisionable(ctx, meta.Metalgo, tt.args.reservationId, tt.args.instanceId, 50*time.Millisecond, 1*time.Second, 50*time.Millisecond); (err != nil) != tt.wantErr {
				t.Errorf("waitUntilReservationProvisionable() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockAPI.Close()
		})
	}
}
