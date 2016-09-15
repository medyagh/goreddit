package main

import(
	"os"
	"sort"
	"net/http"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
)

func init() {
	gothic.Store = sessions.NewFilesystemStore(os.TempDir(), []byte("goth-example"))
	userSession = sessions.NewSession(gothic.Store, "userSession")
}

func initAuth() {
	goth.UseProviders(
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://goreddit-v1.us-west-2.elasticbeanstalk.com/auth/facebook/callback"),
		// twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:4949/auth/twitter/callback"),
	)

	m := make(map[string]string)
	m["facebook"] = "Facebook"
	// m["twitter"] = "Twitter"

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
}

//begins authentication process for a new login
func beginAuth(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}
//complete user authentication and set ["user"]
//session value with gothic
func completeAuth(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	userSession.Values["user"] = user
	http.Redirect(w, r, "/", 301)
}

func getName() string {
	k, _ := userSession.Values["user"].(goth.User)
	return k.Name
}

type ScanData struct{
	Topics []Topic
	User interface{}
}

type GetData struct {
	Topic Topic
	User interface{}
}

type User struct {
	User interface{}
}
type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}
