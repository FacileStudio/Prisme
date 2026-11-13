package auth

import "net/http"

// SessionUser is the user whose session or bearer token the request carries,
// or false when it carries none that resolves. It authorizes nothing: the
// rate limiter reads it to count a signed-in user's requests in a bucket of
// the user's own, and Require still authenticates the request by itself. The
// cookie a dead session would clear is dropped here, since Require clears it
// on the real response.
func (s *Service) SessionUser(r *http.Request) (int64, bool) {
	identity, err := s.sessions.Authenticate(discardedResponse{header: http.Header{}}, r)
	if err != nil || identity.UserID == 0 {
		return 0, false
	}
	return identity.UserID, true
}

// discardedResponse is a ResponseWriter that keeps nothing.
type discardedResponse struct{ header http.Header }

func (d discardedResponse) Header() http.Header { return d.header }

func (discardedResponse) Write(body []byte) (int, error) { return len(body), nil }

func (discardedResponse) WriteHeader(int) {}
