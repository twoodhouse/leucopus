package timeKeeper

import "testing"

func TestTimekeeper(t *testing.T) {
	tk := New([]string{"http://buttonlight01.env:5001/pushbutton", "http://buttonlight01.env:5001/resetbutton"}, []string{}, "http://buttonlight01.env:5001/checklight")
	tk.Begin()
}
