package tests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/infinity-ocean/ikakbolit/pkg/rest"
)

func (s *Suite) TestGetScheduleIDs() {
	const endpoint = "/schedules"

	rq := s.Require()
	ctx := context.Background()

	testCases := []struct {
		name            string
		userID          int
		expectCode      int
		expectSchedules []int
	}{
		{
			name:            "success",
			userID:          6,
			expectCode:      http.StatusOK,
			expectSchedules: []int{6},
		},
		{
			name:       "invalid_user_id",
			userID:     0,
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			url := fmt.Sprintf("%s?user_id=%d", endpoint, tc.userID)

			var errBody rest.APIError

			var ids []int
			resp, err := s.apiClient.Get(ctx, url, nil, &ids, &errBody)
			rq.NoError(err)
			defer rq.NoError(resp.Body.Close())

			rq.Equal(tc.expectCode, resp.StatusCode, "status code should match")

			if resp.StatusCode == http.StatusOK {
				rq.Equal(tc.expectSchedules, ids, "schedules list should match expected IDs")
			} else {
				rq.NotEmpty(errBody.Message, "error message should be present")
			}
		})
	}
}
