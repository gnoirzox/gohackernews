package storyactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/gnoirzox/gohackernews/src/lib/session"
	"github.com/gnoirzox/gohackernews/src/stories"
)

// HandleUpdateShow renders the form to update a story.
func HandleUpdateShow(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the story
	story, err := stories.Find(params.GetInt(stories.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Authorise update story
	currentUser := session.CurrentUser(w, r)
	err = can.Update(story, currentUser)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("story", story)
	view.AddKey("currentUser", currentUser)
	return view.Render()
}

// HandleUpdate handles the POST of the form to update a story
func HandleUpdate(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the story
	story, err := stories.Find(params.GetInt(stories.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Check the authenticity token
	err = session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise update story
	currentUser := session.CurrentUser(w, r)
	err = can.Update(story, currentUser)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Clean params according to role
	accepted := stories.AllowedParams()
	if currentUser.Admin() {
		accepted = stories.AllowedParamsAdmin()
	}
	storyParams := story.ValidateParams(params.Map(), accepted)

	err = story.Update(storyParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to story
	return server.Redirect(w, r, story.ShowURL())
}
