package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"avito-tenders/internal/api/bids/dtos"
)

func (s *TestSuite) TestCreateBid() {
	type want struct {
		StatusCode int
	}
	type args struct {
		inputFileName string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Create bid from user",
			args: args{
				inputFileName: "bids/new/create_user.json",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Create bid from organization",
			args: args{
				inputFileName: "bids/new/create_organization.json",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Created tender",
			args: args{
				inputFileName: "bids/new/created_tender.json",
			},
			want: want{
				StatusCode: 403,
			},
		},
		{
			name: "Closed tender",
			args: args{
				inputFileName: "bids/new/closed_tender.json",
			},
			want: want{
				StatusCode: 403,
			},
		},
		{
			name: "Unauthorized organization",
			args: args{
				inputFileName: "bids/new/unauthorized_organization.json",
			},
			want: want{
				StatusCode: 403,
			},
		},
		{
			name: "Unknown user",
			args: args{
				inputFileName: "bids/new/unknown_user.json",
			},
			want: want{
				StatusCode: 401,
			},
		},
		{
			name: "Missing name",
			args: args{
				inputFileName: "bids/new/missing_name.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Missing description",
			args: args{
				inputFileName: "bids/new/missing_description.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Missing tender id",
			args: args{
				inputFileName: "bids/new/missing_tender_id.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Missing author type",
			args: args{
				inputFileName: "bids/new/missing_author_type.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Missing author id",
			args: args{
				inputFileName: "bids/new/missing_author_id.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Unknown tender id",
			args: args{
				inputFileName: "bids/new/unknown_tender.json",
			},
			want: want{
				StatusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			requestBody := s.loader.LoadString(fmt.Sprintf("%s/%s", fixturesPath, tt.args.inputFileName))

			res, err := s.server.Client().Post(fmt.Sprintf("%s/api/bids/new", s.server.URL), "", bytes.NewBufferString(requestBody))
			require.NoError(t, err)

			defer res.Body.Close()

			require.Equal(t, tt.want.StatusCode, res.StatusCode)
			if tt.want.StatusCode != 200 {
				return
			}

			// check response
			var response dtos.BidResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			require.NoError(t, err)

			expected := s.loader.LoadTemplate(fmt.Sprintf("%s/%s.result", fixturesPath, tt.args.inputFileName), map[string]interface{}{
				"id":        response.ID,
				"createdAt": response.CreatedAt,
			})

			JSONEq(t, expected, response)
		})
	}
}

func (s *TestSuite) TestGetMyBids() {
	type want struct {
		StatusCode int
		Len        int
	}
	type args struct {
		username       string
		limit          string
		offset         string
		outputFileName string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		// 550e8400-e29b-41d4-a716-44665544000a
		{
			name: "Get my bids user 9",
			args: args{
				username:       "user9",
				outputFileName: "bids/my/user_9_tenders.json.result",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get my bids user 10",
			args: args{
				username:       "user10",
				outputFileName: "bids/my/user_10_tenders.json.result",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get my bids from organization's employee account",
			args: args{
				username: "user8",
			},
			want: want{
				StatusCode: 200,
				Len:        0,
			},
		},
		{
			name: "Get my bids limit",
			args: args{
				username:       "user9",
				limit:          "1",
				outputFileName: "bids/my/user_9_tenders_limit_1.json.result",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get my tenders offset",
			args: args{
				username:       "user9",
				offset:         "1",
				outputFileName: "bids/my/user_9_tenders_offset_1.json.result",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get my tenders for user that didn't create them",
			args: args{
				username: "user30",
			},
			want: want{
				StatusCode: 200,
				Len:        0,
			},
		},
		{
			name: "Unknown user",
			args: args{
				username: "user40",
			},
			want: want{
				StatusCode: 401,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			baseURL := fmt.Sprintf("%s/api/bids/my", s.server.URL)
			v := url.Values{}
			v.Add("username", tt.args.username)
			if len(tt.args.limit) != 0 {
				v.Add("limit", tt.args.limit)
			}
			if len(tt.args.offset) != 0 {
				v.Add("offset", tt.args.offset)
			}

			res, err := s.server.Client().Get(fmt.Sprintf("%s?%s", baseURL, v.Encode()))
			require.NoError(t, err)

			defer res.Body.Close()

			require.Equal(t, tt.want.StatusCode, res.StatusCode)
			if tt.want.StatusCode != 200 {
				return
			}

			// check response
			var response []dtos.BidResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			require.NoError(t, err)

			if tt.args.outputFileName != "" {
				expected := s.loader.LoadString(fmt.Sprintf("%s/%s", fixturesPath, tt.args.outputFileName))

				JSONEq(t, expected, response)
			} else {
				assert.Equal(t, tt.want.Len, len(response))
			}
		})
	}
}

// func (s *TestSuite) TestGetMyTenders() {
// 	type want struct {
// 		StatusCode int
// 		Len        int
// 	}
// 	type args struct {
// 		username string
// 		limit    string
// 		offset   string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want want
// 	}{
// 		{
// 			name: "Get my tenders",
// 			args: args{
// 				username: "user1",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 				Len:        2,
// 			},
// 		},
// 		{
// 			name: "Get my tenders from organization's employee account",
// 			args: args{
// 				username: "user2",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 				Len:        2,
// 			},
// 		},
// 		{
// 			name: "Get my tenders limit",
// 			args: args{
// 				username: "user1",
// 				limit:    "1",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 				Len:        1,
// 			},
// 		},
// 		{
// 			name: "Get my tenders offset",
// 			args: args{
// 				username: "user1",
// 				offset:   "1",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 				Len:        1,
// 			},
// 		},
// 		{
// 			name: "Get my tenders limit and offset",
// 			args: args{
// 				username: "user1",
// 				offset:   "2",
// 				limit:    "5",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 				Len:        0,
// 			},
// 		},
// 		{
// 			name: "Get my tenders for user that didn't create them",
// 			args: args{
// 				username: "user30",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 				Len:        0,
// 			},
// 		},
// 		{
// 			name: "Unknown user",
// 			args: args{
// 				username: "user40",
// 			},
// 			want: want{
// 				StatusCode: 401,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		s.T().Run(tt.name, func(t *testing.T) {
// 			baseUrl := fmt.Sprintf("%s/api/tenders/my", s.server.URL)
// 			v := url.Values{}
// 			v.Add("username", tt.args.username)
// 			if len(tt.args.limit) != 0 {
// 				v.Add("limit", tt.args.limit)
// 			}
// 			if len(tt.args.offset) != 0 {
// 				v.Add("offset", tt.args.offset)
// 			}
//
// 			res, err := s.server.Client().Get(fmt.Sprintf("%s?%s", baseUrl, v.Encode()))
// 			require.NoError(t, err)
//
// 			defer res.Body.Close()
//
// 			require.Equal(t, tt.want.StatusCode, res.StatusCode)
// 			if tt.want.StatusCode != 200 {
// 				return
// 			}
//
// 			// check response
// 			var response []dtos.TenderResponse
// 			err = json.NewDecoder(res.Body).Decode(&response)
// 			require.NoError(t, err)
//
// 			assert.Equal(t, tt.want.Len, len(response))
// 		})
// 	}
// }
//
// func (s *TestSuite) TestGetTenderStatusByID() {
// 	type want struct {
// 		StatusCode int
// 	}
// 	type args struct {
// 		username       string
// 		tenderID       string
// 		outputFileName string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want want
// 	}{
// 		{
// 			name: "Get published tender from creator account",
// 			args: args{
// 				username:       "user4",
// 				tenderID:       "550e8400-e29b-41d4-a716-446655440041",
// 				outputFileName: "tenders/status/tender2.json.result",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 			},
// 		},
// 		{
// 			name: "Get published tender from not creator account",
// 			args: args{
// 				username:       "user1",
// 				tenderID:       "550e8400-e29b-41d4-a716-446655440041",
// 				outputFileName: "tenders/status/tender2.json.result",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 			},
// 		},
// 		{
// 			name: "Get created tender from creator account",
// 			args: args{
// 				username:       "user4",
// 				tenderID:       "550e8400-e29b-41d4-a716-446655440040",
// 				outputFileName: "tenders/status/tender1.json.result",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 			},
// 		},
// 		{
// 			name: "Get created tender from organization's employee account",
// 			args: args{
// 				username:       "user5",
// 				tenderID:       "550e8400-e29b-41d4-a716-446655440040",
// 				outputFileName: "tenders/status/tender1.json.result",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 			},
// 		},
// 		{
// 			name: "Get created tender from not creator account",
// 			args: args{
// 				username: "user1",
// 				tenderID: "550e8400-e29b-41d4-a716-446655440040",
// 			},
// 			want: want{
// 				StatusCode: 403,
// 			},
// 		},
// 		{
// 			name: "Get closed tender from creator account",
// 			args: args{
// 				username:       "user4",
// 				outputFileName: "tenders/status/tender3.json.result",
// 				tenderID:       "550e8400-e29b-41d4-a716-446655440042",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 			},
// 		},
// 		{
// 			name: "Get closed tender from organization's employee",
// 			args: args{
// 				username:       "user5",
// 				outputFileName: "tenders/status/tender3.json.result",
// 				tenderID:       "550e8400-e29b-41d4-a716-446655440042",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 			},
// 		},
// 		{
// 			name: "Get closed tender from not creator account",
// 			args: args{
// 				username: "user1",
// 				tenderID: "550e8400-e29b-41d4-a716-446655440042",
// 			},
// 			want: want{
// 				StatusCode: 403,
// 			},
// 		},
// 		{
// 			name: "Get closed tender with invalid user account",
// 			args: args{
// 				username: "user40",
// 				tenderID: "550e8400-e29b-41d4-a716-446655440042",
// 			},
// 			want: want{
// 				StatusCode: 401,
// 			},
// 		},
// 		{
// 			name: "Get unknown tender",
// 			args: args{
// 				tenderID: "550e8400-e29b-41d4-a716-446655440043",
// 			},
// 			want: want{
// 				StatusCode: 404,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		s.T().Run(tt.name, func(t *testing.T) {
// 			baseUrl := fmt.Sprintf("%s/api/tenders/%s/status", s.server.URL, tt.args.tenderID)
//
// 			v := url.Values{}
// 			if len(tt.args.username) != 0 {
// 				v.Add("username", tt.args.username)
// 			}
// 			res, err := s.server.Client().Get(fmt.Sprintf("%s?%s", baseUrl, v.Encode()))
//
// 			require.NoError(t, err)
//
// 			defer res.Body.Close()
//
// 			require.Equal(t, tt.want.StatusCode, res.StatusCode)
// 			if tt.want.StatusCode != 200 {
// 				return
// 			}
//
// 			expected := s.loader.LoadString(fmt.Sprintf("%s/%s", fixturesPath, tt.args.outputFileName))
//
// 			body, err := io.ReadAll(res.Body)
// 			require.NoError(t, err)
//
// 			assert.Equal(t, expected, string(body))
// 		})
// 	}
// }
//
// func (s *TestSuite) TestEditTenderStatusByID() {
// 	type want struct {
// 		StatusCode int
// 	}
// 	type args struct {
// 		username       string
// 		tenderID       string
// 		outputFileName string
// 		newStatus      string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want want
// 	}{
// 		{
// 			name: "Update status from creator account",
// 			args: args{
// 				username:       "user4",
// 				tenderID:       "550e8400-e29b-41d4-a716-446655440041",
// 				outputFileName: "tenders/status/tender2.json.result",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 			},
// 		},
// 		{
// 			name: "Get published tender from not creator account",
// 			args: args{
// 				username:       "user1",
// 				tenderID:       "550e8400-e29b-41d4-a716-446655440041",
// 				outputFileName: "tenders/status/tender2.json.result",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 			},
// 		},
// 		{
// 			name: "Get created tender from creator account",
// 			args: args{
// 				username:       "user4",
// 				tenderID:       "550e8400-e29b-41d4-a716-446655440040",
// 				outputFileName: "tenders/status/tender1.json.result",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 			},
// 		},
// 		{
// 			name: "Get created tender from not creator account",
// 			args: args{
// 				username: "user1",
// 				tenderID: "550e8400-e29b-41d4-a716-446655440040",
// 			},
// 			want: want{
// 				StatusCode: 403,
// 			},
// 		},
// 		{
// 			name: "Get closed tender from creator account",
// 			args: args{
// 				username:       "user4",
// 				outputFileName: "tenders/status/tender3.json.result",
// 				tenderID:       "550e8400-e29b-41d4-a716-446655440042",
// 			},
// 			want: want{
// 				StatusCode: 200,
// 			},
// 		},
// 		{
// 			name: "Get closed tender from not creator account",
// 			args: args{
// 				username: "user1",
// 				tenderID: "550e8400-e29b-41d4-a716-446655440042",
// 			},
// 			want: want{
// 				StatusCode: 403,
// 			},
// 		},
// 		{
// 			name: "Get closed tender with invalid user account",
// 			args: args{
// 				username: "user40",
// 				tenderID: "550e8400-e29b-41d4-a716-446655440042",
// 			},
// 			want: want{
// 				StatusCode: 401,
// 			},
// 		},
// 		{
// 			name: "Get unknown tender",
// 			args: args{
// 				tenderID: "550e8400-e29b-41d4-a716-446655440043",
// 			},
// 			want: want{
// 				StatusCode: 404,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		s.T().Run(tt.name, func(t *testing.T) {
// 			baseUrl := fmt.Sprintf("%s/api/tenders/%s/status", s.server.URL, tt.args.tenderID)
//
// 			v := url.Values{}
// 			if len(tt.args.username) != 0 {
// 				v.Add("username", tt.args.username)
// 			}
//
// 			res, err := s.server.Client().Get(fmt.Sprintf("%s?%s", baseUrl, v.Encode()))
// 			require.NoError(t, err)
//
// 			defer res.Body.Close()
//
// 			require.Equal(t, tt.want.StatusCode, res.StatusCode)
// 			if tt.want.StatusCode != 200 {
// 				return
// 			}
//
// 			expected := s.loader.LoadString(fmt.Sprintf("%s/%s", fixturesPath, tt.args.outputFileName))
//
// 			body, err := io.ReadAll(res.Body)
// 			require.NoError(t, err)
//
// 			assert.Equal(t, expected, string(body))
// 		})
// 	}
// }
