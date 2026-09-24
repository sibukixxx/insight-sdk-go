package insight

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

type recordedTransport struct {
	method, path string
	status       int
	response     string
}

func (r *recordedTransport) RoundTrip(_ context.Context, method, path string, _ []byte) (int, []byte, error) {
	r.method, r.path = method, path
	return r.status, []byte(r.response), nil
}

func TestClientUsesCustomTransportWithoutHTTP(t *testing.T) {
	rt := &recordedTransport{status: http.StatusOK, response: `{"contractVersion":"1","contractSchema":"insight-lab.public-engine","supportedContractVersions":["1"],"engine":{"version":"v","commit":"c","dirty":"false"}}`}
	info, err := NewClient("unused", WithTransport(rt)).GetEngine(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if rt.method != http.MethodGet || rt.path != "/engine" || info.Engine.Version != "v" {
		t.Fatalf("method=%s path=%s info=%+v", rt.method, rt.path, info)
	}
}

func TestClientReturnsUnsupportedContractVersionWhenTransportAnswersOtherVersion(t *testing.T) {
	rt := &recordedTransport{status: http.StatusOK, response: `{"contractVersion":"2"}`}
	_, err := NewClient("unused", WithTransport(rt)).GetEngine(context.Background())
	if !errors.Is(err, &Error{Code: CodeUnsupportedContractVersion}) {
		t.Fatalf("got %v", err)
	}
}
