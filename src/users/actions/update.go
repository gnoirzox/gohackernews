package useractions

import (
	"net/http"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/gnoirzox/gohackernews/src/lib/session"
	"github.com/gnoirzox/gohackernews/src/users"
)

// HandleUpdateShow renders the form to update a user.
func HandleUpdateShow(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the user
	user, err := users.Find(params.GetInt(users.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Authorise update user
	err = can.Update(user, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("user", user)
	return view.Render()
}

// HandleUpdate handles the POST of the form to update a user
func HandleUpdate(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the user
	user, err := users.Find(params.GetInt(users.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Check the authenticity token
	err = session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise update user
	err = can.Update(user, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Set the password hash from the password
	hash, err := auth.HashPassword(params.Get("password"))
	if err != nil {
		return server.InternalError(err)
	}
	// FIXME: For user update we should require the old password too, to match existing
	params.SetString("password_hash", hash)

	// Validate the params, removing any we don't accept
	userParams := user.ValidateParams(params.Map(), users.AllowedParams())

	err = user.Update(userParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to user
	return server.Redirect(w, r, user.ShowURL())
}
