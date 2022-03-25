package middlewares

import (
	"bytes"
	"net/url"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/authelia/authelia/v4/internal/utils"
)

// CORSMiddleware is a special middleware which provides CORS headers via handlers and middleware methods which can be
// configured. It aims to simplify CORS configurations.
type CORSMiddleware struct {
	enabled     bool
	varyOnly    bool
	methods     []byte
	headers     []string
	origins     [][]byte
	credentials []byte
	vary        []byte
	maxAge      []byte
}

// NewCORSMiddleware generates a new automatic CORS policy which adds the Vary header with a value reflecting that the
// Origin header will Vary this response, then if the Origin header has a https scheme it makes the following additional
// adjustments: copies the Origin header to the Access-Control-Allow-Origin header effectively allowing all origins,
// sets the Access-Control-Allow-Credentials header to false which disallows CORS requests from sending cookies etc,
// sets the Access-Control-Allow-Headers header to the value specified by Access-Control-Request-Headers in the request
// excluding the Cookie/Authorization/Proxy-Authorization and special * values, sets Access-Control-Allow-Methods to
// the value specified by the Access-Control-Request-Method header, sets the Access-Control-Max-Age header to 100.
//
// These behaviours can be overridden by the With methods on the returned policy.
//
// The CORS policy can either be used as a middleware using the Middleware method, or as an OPTIONS handler using the
// HandleOPTIONS method.
func NewCORSMiddleware() (policy *CORSMiddleware) {
	return &CORSMiddleware{
		enabled:     true,
		vary:        headerValueVary,
		maxAge:      headerValueMaxAge,
		credentials: headerValueFalse,
	}
}

// WithEnabled changes the enabled state of the middleware. If the middleware is initialized with NewCORSMiddleware this
// value will be true but this function can override the value. Setting it to false prevents the middleware from adding
// any CORS headers. The only effect this middleware has after disabling this is the HandleOPTIONS and HandleOnlyOPTIONS
// handlers still function to return a HTTP 204 No Content, with the Allow header communicating the available HTTP
// method verbs. The main benefit of this option is that you don't have to implement complex logic to add/remove the
// middleware, you can just add it with the Middleware method, and adjust it using the WithEnabled method.
func (p *CORSMiddleware) WithEnabled(enabled bool) (policy *CORSMiddleware) {
	p.enabled = enabled

	return p
}

// WithAllowedMethods takes a list or HTTP methods and adjusts the Access-Control-Allow-Methods header to respond with
// that value.
func (p *CORSMiddleware) WithAllowedMethods(methods ...string) (policy *CORSMiddleware) {
	if len(methods) == 0 {
		p.methods = nil

		return p
	}

	p.methods = []byte(strings.Join(methods, ", "))

	return p
}

// WithAllowedOrigins takes a list of origin strings and only applies the CORS policy if the origin matches one of these.
func (p *CORSMiddleware) WithAllowedOrigins(origins ...string) (policy *CORSMiddleware) {
	if len(origins) == 0 {
		p.origins = nil

		return p
	}

	originsValue := make([][]byte, len(origins))

	for i, origin := range origins {
		if origin == "*" {
			p.origins = [][]byte{[]byte(origin)}

			return p
		}

		originsValue[i] = []byte(origin)
	}

	p.origins = originsValue

	return p
}

// WithAllowedHeaders takes a list of header strings and alters the default Access-Control-Allow-Headers header.
func (p *CORSMiddleware) WithAllowedHeaders(headers ...string) (policy *CORSMiddleware) {
	if len(headers) == 0 {
		p.headers = nil

		return p
	}

	p.headers = headers

	return p
}

// WithAllowCredentials takes bool and alters the default Access-Control-Allow-Credentials header.
func (p *CORSMiddleware) WithAllowCredentials(allow bool) (policy *CORSMiddleware) {
	p.credentials = []byte(strconv.FormatBool(allow))

	return p
}

// WithVary takes a list of header strings and alters the default Vary header.
func (p *CORSMiddleware) WithVary(headers ...string) (policy *CORSMiddleware) {
	if len(headers) == 0 {
		p.vary = nil

		return p
	}

	p.vary = []byte(strings.Join(headers, ", "))

	return p
}

// WithVaryOnly just adds the Vary header.
func (p *CORSMiddleware) WithVaryOnly(varyOnly bool) (policy *CORSMiddleware) {
	p.varyOnly = varyOnly

	return p
}

// WithMaxAge takes an integer and alters the default Access-Control-Max-Age header.
func (p *CORSMiddleware) WithMaxAge(age int) (policy *CORSMiddleware) {
	if age == 0 {
		p.maxAge = nil

		return p
	}

	p.maxAge = []byte(strconv.Itoa(age))

	return p
}

// HandleOPTIONS is an OPTIONS handler that just adds CORS headers, the Allow header, and sets the status code to 204
// without a body. This handler should generally not be used without using WithAllowedMethods.
func (p CORSMiddleware) HandleOPTIONS(ctx *fasthttp.RequestCtx) {
	p.handleOPTIONS(ctx)
	p.handle(ctx)
}

// HandleOnlyOPTIONS is an OPTIONS handler that just handles the Allow header, and sets the status code to 204
// without a body. This handler should generally not be used without using WithAllowedMethods.
func (p CORSMiddleware) HandleOnlyOPTIONS(ctx *fasthttp.RequestCtx) {
	p.handleOPTIONS(ctx)
}

// Middleware provides a middleware that adds the appropriate CORS headers for this CORSMiddleware.
func (p CORSMiddleware) Middleware(next fasthttp.RequestHandler) (handler fasthttp.RequestHandler) {
	return func(ctx *fasthttp.RequestCtx) {
		p.handle(ctx)

		next(ctx)
	}
}

func (p CORSMiddleware) handle(ctx *fasthttp.RequestCtx) {
	if !p.enabled {
		return
	}

	p.handleVary(ctx)

	if !p.varyOnly {
		p.handleCORS(ctx)
	}
}

func (p CORSMiddleware) handleOPTIONS(ctx *fasthttp.RequestCtx) {
	ctx.Response.ResetBody()

	ctx.SetStatusCode(fasthttp.StatusNoContent)

	if len(p.methods) != 0 {
		ctx.Response.Header.SetBytesKV(headerAllow, p.methods)
	}
}

func (p CORSMiddleware) handleVary(ctx *fasthttp.RequestCtx) {
	if len(p.vary) != 0 {
		ctx.Response.Header.SetBytesKV(headerVary, p.vary)
	}
}

func (p CORSMiddleware) handleCORS(ctx *fasthttp.RequestCtx) {
	var (
		originURL *url.URL
		err       error
	)

	origin := ctx.Request.Header.PeekBytes(headerOrigin)

	// Skip processing of any `https` scheme URL that has not expressly been configured.
	if originURL, err = url.Parse(string(origin)); err != nil || (originURL.Scheme != "https" && p.origins == nil) {
		return
	}

	var allowedOrigin []byte

	switch len(p.origins) {
	case 0:
		allowedOrigin = origin
	default:
		for i := 0; i < len(p.origins); i++ {
			if bytes.Equal(p.origins[i], originValueWildcard) {
				allowedOrigin = originValueWildcard
			} else if bytes.Equal(p.origins[i], origin) {
				allowedOrigin = origin
			}
		}

		if len(allowedOrigin) == 0 {
			return
		}
	}

	ctx.Response.Header.SetBytesKV(headerAccessControlAllowOrigin, allowedOrigin)
	ctx.Response.Header.SetBytesKV(headerAccessControlAllowCredentials, p.credentials)

	if p.maxAge != nil {
		ctx.Response.Header.SetBytesKV(headerAccessControlMaxAge, p.maxAge)
	}

	p.handleAllowedHeaders(ctx)

	p.handleAllowedMethods(ctx)
}

func (p CORSMiddleware) handleAllowedMethods(ctx *fasthttp.RequestCtx) {
	switch len(p.methods) {
	case 0:
		if requestMethods := ctx.Request.Header.PeekBytes(headerAccessControlRequestMethod); requestMethods != nil {
			ctx.Response.Header.SetBytesKV(headerAccessControlAllowMethods, requestMethods)
		}
	default:
		ctx.Response.Header.SetBytesKV(headerAccessControlAllowMethods, p.methods)
	}
}

func (p CORSMiddleware) handleAllowedHeaders(ctx *fasthttp.RequestCtx) {
	switch len(p.headers) {
	case 0:
		if headers := ctx.Request.Header.PeekBytes(headerAccessControlRequestHeaders); headers != nil {
			requestedHeaders := strings.Split(string(headers), ",")
			allowHeaders := make([]string, 0, len(requestedHeaders))

			for i := 0; i < len(requestedHeaders); i++ {
				headerTrimmed := strings.Trim(requestedHeaders[i], " ")

				if headerTrimmed == "*" {
					continue
				}

				if bytes.Equal(p.credentials, headerValueTrue) ||
					(!strings.EqualFold(fasthttp.HeaderCookie, headerTrimmed) &&
						!strings.EqualFold(fasthttp.HeaderAuthorization, headerTrimmed) &&
						!strings.EqualFold(fasthttp.HeaderProxyAuthorization, headerTrimmed)) {
					allowHeaders = append(allowHeaders, headerTrimmed)
				}
			}

			if len(allowHeaders) != 0 {
				ctx.Response.Header.SetBytesKV(headerAccessControlAllowHeaders, []byte(strings.Join(allowHeaders, ", ")))
			}
		}
	default:
		headers := p.headers

		if bytes.Equal(p.credentials, headerValueTrue) {
			if !utils.IsStringInSliceFold(fasthttp.HeaderCookie, headers) {
				headers = append(headers, fasthttp.HeaderCookie)
			}

			if !utils.IsStringInSliceFold(fasthttp.HeaderAuthorization, headers) {
				headers = append(headers, fasthttp.HeaderAuthorization)
			}

			if !utils.IsStringInSliceFold(fasthttp.HeaderProxyAuthorization, headers) {
				headers = append(headers, fasthttp.HeaderProxyAuthorization)
			}
		}

		ctx.Response.Header.SetBytesKV(headerAccessControlAllowHeaders, []byte(strings.Join(headers, ", ")))
	}
}
