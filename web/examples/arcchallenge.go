// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/chai2010/gopkg/web"
)

var form = `
<form action="say" method="POST">
	<input name="said">
	<input type="submit">
</form>
`

var users = map[string]string{}

func main() {
	rand.Seed(time.Now().UnixNano())
	web.Config.CookieSecret = "7C19QRmwf3mHZ9CPAaPQ0hsWeufKd"
	web.Get("/", func(ctx *web.Context) string {
		ctx.Redirect(302, "/said")
		return ""
	})
	web.Get("/said", func() string { return form })
	web.Post("/say", func(ctx *web.Context) string {
		uid := fmt.Sprintf("%d\n", rand.Int63())
		ctx.SetSecureCookie("user", uid, 3600)
		users[uid] = ctx.Params["said"]
		return `<a href="/final">Click Here</a>`
	})
	web.Get("/final", func(ctx *web.Context) string {
		uid, _ := ctx.GetSecureCookie("user")
		return "You said " + users[uid]
	})
	web.Run("0.0.0.0:9999")
}
