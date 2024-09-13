package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"avito-tenders/internal/api/tenders/dtos"
)

func (s *TestSuite) TestCreateTender() {
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
			name: "Create tender",
			args: args{
				inputFileName: "tenders/new/create_tender.json",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Create published delivery tender",
			args: args{
				inputFileName: "tenders/new/create_tender_published.json",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Missing name",
			args: args{
				inputFileName: "tenders/new/missing_name.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Missing description",
			args: args{
				inputFileName: "tenders/new/missing_description.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Missing service type",
			args: args{
				inputFileName: "tenders/new/missing_service_type.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Missing organization",
			args: args{
				inputFileName: "tenders/new/missing_organization_id.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Missing creator username",
			args: args{
				inputFileName: "tenders/new/missing_creator_username.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Not from organization",
			args: args{
				inputFileName: "tenders/new/not_from_organization.json",
			},
			want: want{
				StatusCode: 403,
			},
		},
		{
			name: "Unknown user",
			args: args{
				inputFileName: "tenders/new/unknown_user.json",
			},
			want: want{
				StatusCode: 401,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			requestBody := s.loader.LoadString(fmt.Sprintf("%s/%s", fixturesPath, tt.args.inputFileName))

			res, err := s.server.Client().Post(fmt.Sprintf("%s/api/tenders/new", s.server.URL), "", bytes.NewBufferString(requestBody))
			require.NoError(t, err)

			defer res.Body.Close()

			require.Equal(t, tt.want.StatusCode, res.StatusCode)
			if tt.want.StatusCode != 200 {
				return
			}

			// check response
			var response dtos.TenderResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			require.NoError(t, err)

			expected := s.loader.LoadTemplate(fmt.Sprintf("%s/%s.result", fixturesPath, tt.args.inputFileName), map[string]interface{}{
				"id":        response.Id,
				"createdAt": response.CreatedAt,
			})

			JSONEq(t, expected, response)
		})
	}
}

func (s *TestSuite) TestGetMyTenders() {
	type want struct {
		StatusCode int
		Len        int
	}
	type args struct {
		username string
		limit    string
		offset   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Get my tenders",
			args: args{
				username: "user1",
			},
			want: want{
				StatusCode: 200,
				Len:        2,
			},
		},
		{
			name: "Get my tenders from organization's employee account",
			args: args{
				username: "user2",
			},
			want: want{
				StatusCode: 200,
				Len:        2,
			},
		},
		{
			name: "Get my tenders limit",
			args: args{
				username: "user1",
				limit:    "1",
			},
			want: want{
				StatusCode: 200,
				Len:        1,
			},
		},
		{
			name: "Get my tenders offset",
			args: args{
				username: "user1",
				offset:   "1",
			},
			want: want{
				StatusCode: 200,
				Len:        1,
			},
		},
		{
			name: "Get my tenders limit and offset",
			args: args{
				username: "user1",
				offset:   "2",
				limit:    "5",
			},
			want: want{
				StatusCode: 200,
				Len:        0,
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
			baseUrl := fmt.Sprintf("%s/api/tenders/my", s.server.URL)
			v := url.Values{}
			v.Add("username", tt.args.username)
			if len(tt.args.limit) != 0 {
				v.Add("limit", tt.args.limit)
			}
			if len(tt.args.offset) != 0 {
				v.Add("offset", tt.args.offset)
			}

			res, err := s.server.Client().Get(fmt.Sprintf("%s?%s", baseUrl, v.Encode()))
			require.NoError(t, err)

			defer res.Body.Close()

			require.Equal(t, tt.want.StatusCode, res.StatusCode)
			if tt.want.StatusCode != 200 {
				return
			}

			// check response
			var response []dtos.TenderResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Len, len(response))
		})
	}
}

func (s *TestSuite) TestGetTenderStatusByID() {
	type want struct {
		StatusCode int
	}
	type args struct {
		username       string
		tenderId       string
		outputFileName string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Get published tender from creator account",
			args: args{
				username:       "user4",
				tenderId:       "550e8400-e29b-41d4-a716-446655440041",
				outputFileName: "tenders/status/tender2.json.result",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get published tender from not creator account",
			args: args{
				username:       "user1",
				tenderId:       "550e8400-e29b-41d4-a716-446655440041",
				outputFileName: "tenders/status/tender2.json.result",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get created tender from creator account",
			args: args{
				username:       "user4",
				tenderId:       "550e8400-e29b-41d4-a716-446655440040",
				outputFileName: "tenders/status/tender1.json.result",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get created tender from organization's employee account",
			args: args{
				username:       "user5",
				tenderId:       "550e8400-e29b-41d4-a716-446655440040",
				outputFileName: "tenders/status/tender1.json.result",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get created tender from not creator account",
			args: args{
				username: "user1",
				tenderId: "550e8400-e29b-41d4-a716-446655440040",
			},
			want: want{
				StatusCode: 403,
			},
		},
		{
			name: "Get closed tender from creator account",
			args: args{
				username:       "user4",
				outputFileName: "tenders/status/tender3.json.result",
				tenderId:       "550e8400-e29b-41d4-a716-446655440042",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get closed tender from organization's employee",
			args: args{
				username:       "user5",
				outputFileName: "tenders/status/tender3.json.result",
				tenderId:       "550e8400-e29b-41d4-a716-446655440042",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get closed tender from not creator account",
			args: args{
				username: "user1",
				tenderId: "550e8400-e29b-41d4-a716-446655440042",
			},
			want: want{
				StatusCode: 403,
			},
		},
		{
			name: "Get closed tender with invalid user account",
			args: args{
				username: "user40",
				tenderId: "550e8400-e29b-41d4-a716-446655440042",
			},
			want: want{
				StatusCode: 401,
			},
		},
		{
			name: "Get unknown tender",
			args: args{
				tenderId: "550e8400-e29b-41d4-a716-446655440043",
			},
			want: want{
				StatusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			baseUrl := fmt.Sprintf("%s/api/tenders/%s/status", s.server.URL, tt.args.tenderId)

			v := url.Values{}
			if len(tt.args.username) != 0 {
				v.Add("username", tt.args.username)
			}
			res, err := s.server.Client().Get(fmt.Sprintf("%s?%s", baseUrl, v.Encode()))

			require.NoError(t, err)

			defer res.Body.Close()

			require.Equal(t, tt.want.StatusCode, res.StatusCode)
			if tt.want.StatusCode != 200 {
				return
			}

			expected := s.loader.LoadString(fmt.Sprintf("%s/%s", fixturesPath, tt.args.outputFileName))

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, expected, string(body))
		})
	}
}

func (s *TestSuite) TestEditTenderStatusByID() {
	type want struct {
		StatusCode int
	}
	type args struct {
		username       string
		tenderId       string
		outputFileName string
		newStatus      string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Update status from creator account",
			args: args{
				username:       "user4",
				tenderId:       "550e8400-e29b-41d4-a716-446655440041",
				outputFileName: "tenders/status/tender2.json.result",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get published tender from not creator account",
			args: args{
				username:       "user1",
				tenderId:       "550e8400-e29b-41d4-a716-446655440041",
				outputFileName: "tenders/status/tender2.json.result",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get created tender from creator account",
			args: args{
				username:       "user4",
				tenderId:       "550e8400-e29b-41d4-a716-446655440040",
				outputFileName: "tenders/status/tender1.json.result",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get created tender from not creator account",
			args: args{
				username: "user1",
				tenderId: "550e8400-e29b-41d4-a716-446655440040",
			},
			want: want{
				StatusCode: 403,
			},
		},
		{
			name: "Get closed tender from creator account",
			args: args{
				username:       "user4",
				outputFileName: "tenders/status/tender3.json.result",
				tenderId:       "550e8400-e29b-41d4-a716-446655440042",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Get closed tender from not creator account",
			args: args{
				username: "user1",
				tenderId: "550e8400-e29b-41d4-a716-446655440042",
			},
			want: want{
				StatusCode: 403,
			},
		},
		{
			name: "Get closed tender with invalid user account",
			args: args{
				username: "user40",
				tenderId: "550e8400-e29b-41d4-a716-446655440042",
			},
			want: want{
				StatusCode: 401,
			},
		},
		{
			name: "Get unknown tender",
			args: args{
				tenderId: "550e8400-e29b-41d4-a716-446655440043",
			},
			want: want{
				StatusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			baseUrl := fmt.Sprintf("%s/api/tenders/%s/status", s.server.URL, tt.args.tenderId)

			v := url.Values{}
			if len(tt.args.username) != 0 {
				v.Add("username", tt.args.username)
			}

			res, err := s.server.Client().Get(fmt.Sprintf("%s?%s", baseUrl, v.Encode()))
			require.NoError(t, err)

			defer res.Body.Close()

			require.Equal(t, tt.want.StatusCode, res.StatusCode)
			if tt.want.StatusCode != 200 {
				return
			}

			expected := s.loader.LoadString(fmt.Sprintf("%s/%s", fixturesPath, tt.args.outputFileName))

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, expected, string(body))
		})
	}
}
