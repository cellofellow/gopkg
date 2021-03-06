// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chai2010/gopkg/web"
)

const page = `
<html>
<meta charset="utf-8"/>
<body>
{{if .Value}}.
Hi {{.Value}}.
<form method="post" action="/logout">
<input type="submit" name="method" value="logout" />
</form>
You will logout after 10 seconds. Then try to reload.
{{else}}
<form method="post" action="/login">
<label for="name">Name:</label>
<input type="text" id="name" name="name" value="" />
<input type="submit" name="method" value="login" />
</form>
{{end}}
</body>
</html>
`

var tmpl = template.Must(template.New("x").Parse(page))

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	manager := web.NewSessionManager(logger)
	manager.OnStart(func(session *web.Session) {
		println("started new session:", session.Id)
	})
	manager.OnTouch(func(session *web.Session) {
		println("touch session:", session.Id)
	})
	manager.OnEnd(func(session *web.Session) {
		println("abandon:", session.Id)
	})
	manager.SetTimeout(10)

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		session := manager.GetSession(w, req)
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, session)
	}))
	http.Handle("/login", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		name := strings.Trim(req.FormValue("name"), " ")
		if name != "" {
			logger.Printf("User \"%s\" login", name)

			// XXX: set user own object.
			manager.GetSession(w, req).Value = name
		}
		http.Redirect(w, req, "/", http.StatusFound)
	}))
	http.Handle("/logout", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if manager.GetSession(w, req).Value != nil {
			// XXX: get user own object.
			name := manager.GetSession(w, req).Value.(string)
			logger.Printf("User \"%s\" logout", name)
			manager.GetSession(w, req).Abandon()
		}
		http.Redirect(w, req, "/", http.StatusFound)
	}))
	err := http.ListenAndServe(":6061", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
