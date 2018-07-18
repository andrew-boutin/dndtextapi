// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

// Authentication resources utilized:
// - https://medium.com/@hfogelberg/the-black-magic-of-oauth-in-golang-part-1-3cef05c28dde
// - https://skarlso.github.io/2016/11/02/google-signin-with-go-part2/

package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/andrew-boutin/dndtextapi/backends"
	"github.com/andrew-boutin/dndtextapi/configs"
	"github.com/andrew-boutin/dndtextapi/users"

	"github.com/dchest/uniuri"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	// userSessionKey is the key to look up a Users's session in the session store.
	userSessionStoreKey = "USER_SESSION_STORE_KEY"

	// userContextKey is the key to look up the authenticated User in the Context with.
	userContextKey = "USER_CONTEXT_KEY"
)

// googleOauthConfig is all of the config data required to authenticate a User with Google
var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/callback",
	ClientID:     "",
	ClientSecret: "",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint: google.Endpoint,
}

// store is the session store used for authentication
var store cookie.Store

func init() {
	store = cookie.NewStore([]byte(randToken(64)))
	store.Options(sessions.Options{
		Path: "/",
		// 1 week
		MaxAge: 86400 * 7,
	})
}

// InitAuthentication initializes authentication configuration that has
// to be read in from config files
func InitAuthentication(c configs.ClientConfiguration) {
	googleOauthConfig.ClientID = c.ID
	googleOauthConfig.ClientSecret = c.Secret
}

// RegisterAuthenticationRoutes adds the authentication routes
func RegisterAuthenticationRoutes(r *gin.Engine) {
	// Use the cookie store
	r.Use(sessions.Sessions("dndtextapisession", store))

	// r.GET("/", IndexHandler) IndexHandler just gives a link to LoginHandler so that seems not necessary for now
	r.GET("/login", LoginHandler)
	r.GET("/callback", CallbackHandler)
}

// LoginHandler handles redirecting the User to Google for authentication
func LoginHandler(c *gin.Context) {
	oauthStateString := uniuri.New()
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// CodeForm is used to pull the access code out of the request sent to the callback
// handler.
type CodeForm struct {
	Code string `form:"code" binding:"required"`
}

// CallbackHandler handles callbacks from Google after the User has logged in. An
// access code should be sent on successful login that will allow us to get an
// access token and then access their profile data in Google. If all of this suceeds
// then we can consider the User authenticated.
func CallbackHandler(c *gin.Context) {
	dbBackend := GetDBBackend(c)
	// TODO: Random state string is a query param in the callback url - where does it come into play?
	// TODO: Example callback url
	// http://localhost:8080/callback?state=TpHYsn544lycCIb5&code=4/AABxNAZm3-Z4LRRbHTVQvwdX74pVF1o9JG60ktNFGbSj_yTfOHYQDKyK0wcPIW6FBoDQjXnbpdl9UbXM48BnAlY#
	// Authentication provider returns an access code when the User has logged in
	var code string
	var form CodeForm
	if c.ShouldBind(&form) == nil {
		code = form.Code
	}

	// Exchange the access code for an access token that we can make calls with
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.WithError(err).WithField("code", code).Error("Failed to get Google access token using access code.")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Make sure the token is valid
	if !token.Valid() {
		log.WithField("token", token).Error("Google access token not valid.")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Attempt to get the User info from Google using the access token
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.WithError(err).Error("Failed to get Google User data.")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Convert the response from the Google endpoint into userful data
	googleUser, err := readInGoogleUser(response)
	if err != nil {
		log.WithError(err).Error("Failed to extract Google User data.")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// TODO: Could return some data to indicate a returning user or not
	user, err := getOrCreateUser(dbBackend, googleUser)
	if err != nil {
		log.WithError(err).Error("Failed to either look up or create User.")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Create a session for the User and put it in the session store so we can check if they're authenticated later
	err = createUserSession(c, user.Email)
	if err != nil {
		log.WithError(err).Error("Failed to create a User session.")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

// readInGoogleUser converts a Response into Google User data.
func readInGoogleUser(response *http.Response) (googleUser *users.GoogleUser, err error) {
	// Read out the data from the response
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	// Convert the body into a Google User representation
	err = json.Unmarshal(contents, &googleUser)
	return
}

// createUserSession attempts to create a new session in the session store
// for the given email that corresponds with a User.
func createUserSession(c *gin.Context, email string) error {
	session := sessions.Default(c)
	session.Set(userSessionStoreKey, email)
	return session.Save()
}

// getOrCreateUser attempts to lookup a User in the database and creates a new one if one
// doesn't already exist.
func getOrCreateUser(dbBackend backends.Backend, gu *users.GoogleUser) (user *users.User, err error) {
	// Attempt to look up the Google User in the database
	user, err = dbBackend.GetUser(gu.Email)
	if err == users.ErrUserNotFound {
		// Create a new User in the database for this new Google profile
		user, err = dbBackend.CreateUser(gu)
	}
	return
}

// randToken generates a random token of the given length.
func randToken(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// AuthenticationMiddleware requires that the User is authenticated or else they
// get access denied.
func AuthenticationMiddleware(c *gin.Context) {
	// If there is a valid session entry for the User then they're still authenticated
	session := sessions.Default(c)
	emailAsInterface := session.Get(userSessionStoreKey) // TODO: Potentially have the User in the Session
	if emailAsInterface == nil {
		// User doesn't have a session so deny access
		log.Error("No session data found denying access.")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Look up the User and set in the Context so all future middleware can have access
	email := emailAsInterface.(string)
	dbBackend := GetDBBackend(c)
	user, err := dbBackend.GetUser(email)
	if err != nil {
		log.WithError(err).Errorf("Failed to look up User using session email data %s.", email)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Set(userContextKey, user)
}

// GetAuthenticatedUser pulls out the authenticated User from the Context. Previous
// middleware should have set the User in the Context previously or aborted the request
// if there was an issue. Routes that don't require authentication won't have a User
// set in the Context, but they should not be attempting to get the authenticated User.
func GetAuthenticatedUser(c *gin.Context) *users.User {
	return c.MustGet(userContextKey).(*users.User)
}
