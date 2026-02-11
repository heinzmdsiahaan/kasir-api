package middlewares

import "net/http"

// func (api key) func handler http.handler
func APIKey(validApiKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// before function running
			apiKey := r.Header.Get("X-API-Key")

			if apiKey == "" {
				http.Error(w, "API Key Required", http.StatusUnauthorized)
				return
			}

			if apiKey != validApiKey {
				http.Error(w, "Invalid API Key", http.StatusUnauthorized)
				return
			}

			next(w, r)

			// after function running
		}

	}
}
