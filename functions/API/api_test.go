package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/mark-ignacio/member-gentei/pkg/common"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

const (
	testToken = "b7a31354-test-test-test-c3c98d2828a2"
)

var mockRT = &mockRoundTripper{
	responses: map[string]mockRoundTripperResponse{},
}

type mockRoundTripperResponse struct {
	code int
	body string
}

type mockRoundTripper struct {
	responses map[string]mockRoundTripperResponse
	http.RoundTripper
}

func (m *mockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	requestURL := r.URL.String()
	log.Info().Str("path", requestURL).Send()
	if r, exists := m.responses[requestURL]; exists {
		return &http.Response{
			StatusCode: r.code,
			Body:       ioutil.NopCloser(strings.NewReader(r.body)),
		}, nil
	}
	log.Error().Msg("URL not found, returning 404")
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       ioutil.NopCloser(strings.NewReader("{some invalid json}")),
	}, nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	setMember(ctx, fs, "has-one", "123")
	setMember(ctx, fs, "has-two", "789")
	setMember(ctx, fs, "has-two", "456")
	fs.Collection(common.ChannelCollection).Doc("has-one").
		Collection("check").Doc("check").Set(ctx, common.ChannelCheck{VideoID: "test"})
	fs.Collection(common.ChannelCollection).Doc("has-two").
		Collection("check").Doc("check").Set(ctx, common.ChannelCheck{VideoID: "test2"})
	fakeToken := oauth2.Token{
		AccessToken: "test",
		Expiry:      time.Now().Add(time.Hour),
	}
	fs.Collection(common.UsersCollection).Doc("123").
		Collection("private").
		Doc("youtube").Set(ctx, fakeToken)
	fs.Collection(common.UsersCollection).Doc("456").
		Collection("private").
		Doc("youtube").Set(ctx, fakeToken)
	ytClientOptions = []option.ClientOption{
		option.WithHTTPClient(&http.Client{Transport: mockRT}),
	}
	os.Exit(m.Run())
}

func setMember(ctx context.Context, fs *firestore.Client, channelName, discordID string) {
	_, err := fs.Collection(common.ChannelCollection).Doc(channelName).
		Collection(common.ChannelMemberCollection).Doc(discordID).Set(ctx, common.ChannelMember{
		DiscordID: discordID,
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("WriteResult failed")
	}
}

func TestGetMembers(t *testing.T) {
	tests := []struct{ path, body, response string }{
		{path: "/v1/channel/has-one/members", response: `{"users":[{"id":"123"}]}`},
		{path: "/v1/channel/has-two/members", response: `{"users":[{"id":"456"},{"id":"789"}]}`},
		{path: "/v1/channel/has-none/members", response: `{"users":[]}`},
		// pagination
		{path: "/v1/channel/has-two/members?limit=1", response: `{"users":[{"id":"456"}],"after":"456"}`},
		{path: "/v1/channel/has-two/members?limit=1&after=456", response: `{"users":[{"id":"789"}]}`},
		// bogus
		{path: "/v1/channel/has-two/members?limit=-1", response: `{"users":[{"id":"456"},{"id":"789"}]}`},
		{path: "/v1/channel/has-two/members?limit=-1&after=9999", response: `{"users":[]}`},
	}
	for _, testData := range tests {
		req := httptest.NewRequest("GET", testData.path, strings.NewReader(testData.body))
		req.Header.Set("Authorization", "Bearer "+testToken)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		API(rr, req)
		if got := strings.TrimSpace(rr.Body.String()); got != testData.response {
			t.Errorf("API(%q, %q) = %q, want %q", testData.path, testData.body, got, testData.response)
		}
	}
}

func TestCheckMembership(t *testing.T) {
	tests := []struct{ path, body, response string }{
		{path: "/v1/channel/has-one/members/check", body: `{"snowflake": "123"}`, response: `{"member":true}`},
		{path: "/v1/channel/has-one/members/check", body: `{"snowflake": "069"}`, response: `{"member":false,"reason":"not connected"}`},
	}
	mockRT.responses["https://www.googleapis.com/youtube/v3/commentThreads?alt=json&part=id&prettyPrint=false&videoId=test"] = mockRoundTripperResponse{
		code: http.StatusOK,
		body: "{}",
	}
	for _, testData := range tests {
		req := httptest.NewRequest("POST", testData.path, strings.NewReader(testData.body))
		req.Header.Set("Authorization", "Bearer "+testToken)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		API(rr, req)
		if got := strings.TrimSpace(rr.Body.String()); got != testData.response {
			t.Errorf("API(%q, %q) = %q, want %q", testData.path, testData.body, got, testData.response)
		}
		if rr.Code != http.StatusOK {
			t.Errorf("API(%q, %q) code is %d", testData.path, testData.body, rr.Code)
		}
	}
}
