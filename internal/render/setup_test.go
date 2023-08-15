package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/henryeffiong/bookings/internal/config"
	"github.com/henryeffiong/bookings/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	testApp.InProduction = false
	gob.Register(models.Reservation{})
	session = scs.New()
	session.Lifetime = 24 * time.Hour

	session.Cookie.Persist = false
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp
	os.Exit(m.Run())
}

type MyWriter struct{}

func (tw *MyWriter) Header() http.Header {
	var mh http.Header
	return mh
}

func (tw *MyWriter) WriteHeader(i int) {}

func (tw *MyWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
