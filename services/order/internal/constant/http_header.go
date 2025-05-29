package constant

const (
	// HTTP Header Standard
	RequestID      string = `X-Request-Id`
	RequestMethod  string = `x-request-method`
	RequestScheme  string = `x-request-scheme`
	KeyServerRoute string = `x-key-server-route`
	ForwardedFor   string = `x-forwarded-for`
	XClientID      string = `X-Client-Id`
	XClientVersion string = `X-Client-Version`
	XUserEmail     string = `x-user-email`

	// Custom HTTP Header
	AppToken string = `x-app-token`

	// Cache Control Header
	CacheControl        string = `Cache-Control`
	CacheNoCache        string = `no-cache`
	CacheNoStore        string = `no-store`
	CacheMustRevalidate string = `must-revalidate`
)
