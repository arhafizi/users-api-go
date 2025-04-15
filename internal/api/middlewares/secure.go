package middlewares

import "github.com/gin-contrib/secure"

var Secure = secure.New(secure.Config{
	// AllowedHosts: []string{"example.com", "ssl.example.com"},
	// SSLRedirect:           true,
	// STSSeconds:            315360000,
	// STSIncludeSubdomains:  true,
	IsDevelopment:         true,
	FrameDeny:             true,                              // clickjacking
	ContentTypeNosniff:    true,                              // MIME-sniffing
	BrowserXssFilter:      true,                              // XSS p
	IENoOpen:              true,                              // file execution in IE
	ContentSecurityPolicy: "default-src 'self'",              // Basic CSP for development
	ReferrerPolicy:        "strict-origin-when-cross-origin", // Referrer policy
})
