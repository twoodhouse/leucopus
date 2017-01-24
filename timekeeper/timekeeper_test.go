package timeKeeper

import "testing"

func TestTimekeeper(t *testing.T) {
	tk := New()
	tk.InitInfo("A")
	tk.InitInfo("B")
	tk.Begin()
}
