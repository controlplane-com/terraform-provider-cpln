package cpln

import (
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
)

/*** Unit Tests ***/

// TestClientRetriesTransientFailures verifies which failed requests the client sends again and how often.
func TestClientRetriesTransientFailures(t *testing.T) {
	// Shorten the wait between attempts so the cases run quickly
	previousDelay := client.TransientRetryBaseDelay
	client.TransientRetryBaseDelay = time.Millisecond
	t.Cleanup(func() { client.TransientRetryBaseDelay = previousDelay })

	// Outcomes the fake transport can replay
	gatewayConnectionFailure := fakeOutcome{status: http.StatusServiceUnavailable, body: "upstream connect error or disconnect/reset before headers. retried and the latest reset reason: remote connection failure, transport failure reason: delayed connect error: Connection refused"}
	gatewayConnectionTermination := fakeOutcome{status: http.StatusServiceUnavailable, body: "upstream connect error or disconnect/reset before headers. reset reason: connection termination"}
	connectionRefused := fakeOutcome{err: &net.OpError{Op: "dial", Net: "tcp", Err: &os.SyscallError{Syscall: "connect", Err: syscall.ECONNREFUSED}}}
	hostNotFound := fakeOutcome{err: &net.OpError{Op: "dial", Net: "tcp", Err: &net.DNSError{Err: "no such host", Name: "api.example.test", IsNotFound: true}}}
	connectionDropped := fakeOutcome{err: io.ErrUnexpectedEOF}
	timedOut := fakeOutcome{err: fakeTimeoutError{}}
	ok := fakeOutcome{status: http.StatusOK, body: "{}"}

	// Define the table of cases
	cases := []struct {
		name         string
		call         func(c *client.Client) error
		outcomes     []fakeOutcome
		wantAttempts int
		wantErr      bool
	}{
		{
			name:         "update is sent again while the gateway cannot reach the API",
			call:         updateResource,
			outcomes:     []fakeOutcome{gatewayConnectionFailure, gatewayConnectionFailure, ok},
			wantAttempts: 3,
		},
		{
			name:         "create is sent again after a refused connection",
			call:         createResource,
			outcomes:     []fakeOutcome{connectionRefused, {status: http.StatusCreated, body: "{}"}},
			wantAttempts: 2,
		},
		{
			name:         "delete is sent again while the gateway cannot reach the API",
			call:         deleteResource,
			outcomes:     []fakeOutcome{gatewayConnectionFailure, ok},
			wantAttempts: 2,
		},
		{
			name:         "update is not sent again after the gateway reached the API",
			call:         updateResource,
			outcomes:     []fakeOutcome{gatewayConnectionTermination},
			wantAttempts: 1,
			wantErr:      true,
		},
		{
			name:         "update is not sent again after a dropped connection",
			call:         updateResource,
			outcomes:     []fakeOutcome{connectionDropped},
			wantAttempts: 1,
			wantErr:      true,
		},
		{
			name:         "read is sent again after a dropped connection",
			call:         getResource,
			outcomes:     []fakeOutcome{connectionDropped, ok},
			wantAttempts: 2,
		},
		{
			name:         "read is sent again after gateway errors",
			call:         getResource,
			outcomes:     []fakeOutcome{{status: http.StatusBadGateway}, {status: http.StatusGatewayTimeout}, gatewayConnectionTermination, ok},
			wantAttempts: 4,
		},
		{
			name:         "read stops at the attempt limit",
			call:         getResource,
			outcomes:     []fakeOutcome{{status: http.StatusBadGateway}},
			wantAttempts: 5,
			wantErr:      true,
		},
		{
			name:         "read is not sent again when the host does not resolve",
			call:         getResource,
			outcomes:     []fakeOutcome{hostNotFound},
			wantAttempts: 1,
			wantErr:      true,
		},
		{
			name:         "read is not sent again after a timeout",
			call:         getResource,
			outcomes:     []fakeOutcome{timedOut},
			wantAttempts: 1,
			wantErr:      true,
		},
		{
			name:         "read is not sent again after an API error",
			call:         getResource,
			outcomes:     []fakeOutcome{{status: http.StatusInternalServerError}},
			wantAttempts: 1,
			wantErr:      true,
		},
	}

	// Run each case
	for _, tc := range cases {
		// Run the case as a subtest
		t.Run(tc.name, func(t *testing.T) {
			// Point a client at a transport that replays the case outcomes
			transport := &fakeTransport{outcomes: tc.outcomes}
			c := &client.Client{HostURL: "https://api.example.test", Org: "org", Token: "token", HTTPClient: &http.Client{Transport: transport}}

			// Invoke the request under test
			err := tc.call(c)

			// Verify the request failed or succeeded as expected
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tc.wantErr)
			}

			// Verify the request was sent the expected number of times
			if len(transport.bodies) != tc.wantAttempts {
				t.Fatalf("attempts = %d, want %d", len(transport.bodies), tc.wantAttempts)
			}

			// Verify every attempt carried the same body as the first
			for i, body := range transport.bodies {
				if body != transport.bodies[0] {
					t.Fatalf("attempt %d sent body %q, want %q", i+1, body, transport.bodies[0])
				}
			}
		})
	}
}

/*** Unit Test Helpers ***/

// fakeOutcome is the response or transport error the fake transport returns for one attempt.
type fakeOutcome struct {
	status int
	body   string
	err    error
}

// fakeTransport replays one outcome per attempt, repeating the last one, and records the body each attempt carried.
type fakeTransport struct {
	outcomes []fakeOutcome
	bodies   []string
}

// RoundTrip records the attempt and returns its outcome.
func (f *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Record the body this attempt carried
	body := ""
	if req.Body != nil {
		raw, _ := io.ReadAll(req.Body)
		body = string(raw)
	}
	f.bodies = append(f.bodies, body)

	// Select the outcome for this attempt, repeating the last one once the sequence runs out
	outcome := f.outcomes[min(len(f.bodies), len(f.outcomes))-1]

	// Return the transport error when the outcome carries one
	if outcome.err != nil {
		return nil, outcome.err
	}

	// Return the response the outcome describes
	return &http.Response{StatusCode: outcome.status, Body: io.NopCloser(strings.NewReader(outcome.body)), Header: make(http.Header), Request: req}, nil
}

// fakeTimeoutError is a network error that reports a timeout.
type fakeTimeoutError struct{}

// Error describes the timeout.
func (fakeTimeoutError) Error() string { return "i/o timeout" }

// Timeout reports that the error is a timeout.
func (fakeTimeoutError) Timeout() bool { return true }

// Temporary reports that the error is temporary.
func (fakeTimeoutError) Temporary() bool { return true }

// retryTestResource is a minimal resource body for create and update requests.
type retryTestResource struct {
	Name string `json:"name"`
}

// getResource issues a read through the client.
func getResource(c *client.Client) error {
	_, _, err := c.Get("gvc/example", &retryTestResource{})
	return err
}

// createResource issues a create through the client.
func createResource(c *client.Client) error {
	_, err := c.CreateResource("gvc", "example", &retryTestResource{Name: "example"})
	return err
}

// updateResource issues an update through the client.
func updateResource(c *client.Client) error {
	_, err := c.UpdateResource("gvc/example", &retryTestResource{Name: "example"})
	return err
}

// deleteResource issues a delete through the client.
func deleteResource(c *client.Client) error {
	return c.DeleteResource("gvc/example")
}
