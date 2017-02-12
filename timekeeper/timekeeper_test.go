package timeKeeper

import (
	"net/http"
	"testing"
)

func TestTimekeeper(t *testing.T) {
	http.Get("http://buttonlight01.env:5001/resetall")
	tk := New([]string{"http://buttonlight01.env:5001/pushbutton", "http://buttonlight01.env:5001/resetbutton"}, []string{}, "http://buttonlight01.env:5001/checklight")
	tk.Begin()
}
