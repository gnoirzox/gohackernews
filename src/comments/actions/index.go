package commentactions

import (
	"net/http"

	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/gnoirzox/gohackernews/src/comments"
	"github.com/gnoirzox/gohackernews/src/lib/session"
)

// HandleIndex displays a list of comments.
func HandleIndex(w http.ResponseWriter, r *http.Request) error {

	// Get the params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Build a query
	q := comments.Query().Order("created_at desc").Limit(100)

	// Filter on user id - we only show the actual user's comments
	// so not a nested view as in HN
	userID := params.GetInt("u")
	if userID > 0 {
		q.Where("user_id=?", userID)
	}

	// Filter if requested
	filter := params.Get("filter")
	if len(filter) > 0 {
		q.Where("name ILIKE ?", filter)
	}

	// Fetch the comments
	results, err := comments.FindAll(q)
	if err != nil {
		return server.InternalError(err)
	}

	// Get current user
	currentUser := session.CurrentUser(w, r)

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("filter", filter)
	view.AddKey("comments", results)
	view.AddKey("currentUser", currentUser)
	return view.Render()
}
