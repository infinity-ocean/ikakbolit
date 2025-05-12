package tests

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/infinity-ocean/ikakbolit/pkg/rest"
)

func (s *Suite) TestAddSchedule() {
	const endpoint = "/schedule"

	rq := s.Require()
	ctx := context.Background()

	testCases := []struct {
		name       string
		reqBody    map[string]any
		expectCode int
	}{
		{
			name: "success",
			reqBody: map[string]any{
				"user_id":       1,
				"cure_name":     "Arbidol",
				"doses_per_day": 5,
				"duration_days": 5,
			},
			expectCode: http.StatusCreated,
		},
		{
			name: "zero_duration",
			reqBody: map[string]any{
				"user_id":       1,
				"cure_name":     "Paracetamol",
				"doses_per_day": 5,
				"duration_days": 0,
			},
			expectCode: http.StatusCreated,
		},
		{
			name: "invalid_doses_0",
			reqBody: map[string]any{
				"user_id":       1,
				"cure_name":     "Paracetamol",
				"doses_per_day": 0,
				"duration_days": 5,
			},
			expectCode: http.StatusBadRequest,
		},
		{
			name: "invalid_doses_100",
			reqBody: map[string]any{
				"user_id":       1,
				"cure_name":     "Paracetamol",
				"doses_per_day": 100,
				"duration_days": 5,
			},
			expectCode: http.StatusBadRequest,
		},
		{
			name: "invalid_user_id",
			reqBody: map[string]any{
				"user_id":       0,
				"cure_name":     "Paracetamol",
				"doses_per_day": 5,
				"duration_days": 5,
			},
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var respBody rest.ResponseScheduleID
			var errBody rest.APIError
			log.Printf("Request body: %+v", tc.reqBody)
			resp, err := s.apiClient.Post(ctx, endpoint, nil, tc.reqBody, &respBody, &errBody)
			rq.NoError(err)

			defer func() { rq.NoError(resp.Body.Close()) }()

			rq.Equal(tc.expectCode, resp.StatusCode, "status code should match")

			if resp.StatusCode == http.StatusCreated {
				id, convErr := strconv.Atoi(respBody.ScheduleID)
				rq.NoError(convErr, "schedule_id should parse as integer")
				rq.GreaterOrEqual(id, 10, "schedule_id >= 10")
				rq.LessOrEqual(id, 40, "schedule_id <= 40")
			} else {
				rq.NotEmpty(errBody.Message, "error message should be present")
			}
		})
	}
}
