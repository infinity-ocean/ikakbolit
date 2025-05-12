package tests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/infinity-ocean/ikakbolit/pkg/rest"
)

func (s *Suite) TestGetSchedule() {
	const endpoint = "/schedule"

	rq := s.Require()
	ctx := context.Background()

	testCases := []struct {
		name       string
		userID     int
		scheduleID int
		expectCode int
	}{
		{
			name:       "success",
			userID:     4,
			scheduleID: 4,
			expectCode: http.StatusOK,
		},
		{
			name:       "invalid_user_id",
			userID:     0,
			scheduleID: 1,
			expectCode: http.StatusNotFound,
		},
		{
			name:       "not_found",
			userID:     1,
			scheduleID: 2,
			expectCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			url := fmt.Sprintf("%s?user_id=%d&schedule_id=%d", endpoint, tc.userID, tc.scheduleID)

			var respBody rest.Schedule
			var errBody rest.APIError

			resp, err := s.apiClient.Get(ctx, url, nil, &respBody, &errBody)
			rq.NoError(err)
			defer rq.NoError(resp.Body.Close())

			rq.Equal(tc.expectCode, resp.StatusCode, "status code should match")

			if resp.StatusCode == http.StatusOK {
				rq.Equal(tc.userID, respBody.UserID)
				rq.Equal(tc.scheduleID, respBody.ID)
			} else {
				rq.NotEmpty(errBody.Message, "error message should be present")
			}
		})
	}
}
