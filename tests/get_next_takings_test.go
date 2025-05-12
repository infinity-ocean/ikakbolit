package tests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/infinity-ocean/ikakbolit/pkg/rest"
)

func (s *Suite) TestGetNextTakings() {
	const endpoint = "/next_takings"

	rq := s.Require()
	ctx := context.Background()

	testCases := []struct {
		name       string
		userID     int
		expectCode int
		expectLen  int
	}{
		{
			name:       "success",
			userID:     5,
			expectCode: http.StatusOK,
			expectLen:  1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			url := fmt.Sprintf("%s?user_id=%d", endpoint, tc.userID)

			var respBody rest.SchedulesInWindow
			var errBody rest.APIError

			resp, err := s.apiClient.Get(ctx, url, nil, &respBody, &errBody)
			rq.NoError(err)
			defer rq.NoError(resp.Body.Close())

			rq.Equal(tc.expectCode, resp.StatusCode, "status code should match")

			if resp.StatusCode == http.StatusOK {
				rq.Len(respBody.Schedules, tc.expectLen, "number of schedules should match")
				for _, sch := range respBody.Schedules {
					rq.Equal(tc.userID, sch.UserID)
					rq.NotEmpty(sch.CureName)
					rq.NotEmpty(sch.Intakes)
					rq.GreaterOrEqual(sch.DosesPerDay, 1)
					rq.GreaterOrEqual(sch.DurationDays, 1)
				}
			} else {
				rq.NotEmpty(errBody.Message, "error message should be present")
			}
		})
	}
}
