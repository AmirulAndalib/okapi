/*
 * Copyright 2025 Jonas Kaninda
 *
 * Licensed under the MIT License.
 */

package okapi

import (
	"fmt"
	"net/http"
)

// mustRegister registers a handler and stops the program if the route table
// cannot accept it.
//
// njia reports a bad registration instead of panicking, which is what a gateway
// building routes from user configuration needs. Okapi's routes are written in
// Go and fixed at compile time, so a rejected one is a bug in the program —
// like the empty path addRoute already panics on — and failing at startup beats
// serving a table that is quietly missing a route.
//
// The registrations that can fail are worth naming, because gorilla accepted
// all of them:
//
//   - the same method and path registered twice, which gorilla resolved in
//     favour of whichever came first;
//   - two routes naming the same path position differently, such as
//     "/users/{id}" alongside "/users/{name}", which njia rejects because the
//     captured name would otherwise depend on registration order.
func (o *Okapi) mustRegister(method, path string, h http.HandlerFunc) {
	if err := o.router.njia.HandleFunc(method, path, h); err != nil {
		panic(fmt.Sprintf("okapi: cannot register route %s %s: %v", method, path, err))
	}
}

// mountPrefix serves a prefix and everything beneath it from one handler, for
// the given methods.
//
// A catch-all also matches the prefix itself with nothing after it, so
// "/static/{any...}" covers "/static", "/static/" and "/static/css/app.css" —
// which is the whole of what gorilla's PathPrefix did here, minus its habit of
// matching a bare string prefix and so swallowing "/staticky" too.
func (o *Okapi) mountPrefix(prefix string, h http.HandlerFunc, methods ...string) {
	pattern := mountPattern(prefix)
	for _, m := range methods {
		o.mustRegister(m, pattern, h)
	}
}
