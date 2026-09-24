package join

import "testing"

func TestJoinHorizontal(t *testing.T) {
	o := Options{
		Text:       []string{"Hello", "World"},
		Align:      "left",
		Horizontal: true,
	}
	if err := o.Run(); err != nil {
		t.Errorf("Horizontal join failed: %v", err)
	}
}

func TestJoinVertical(t *testing.T) {
	o := Options{
		Text:     []string{"Hello", "World"},
		Align:    "left",
		Vertical: true,
	}
	if err := o.Run(); err != nil {
		t.Errorf("Vertical join failed: %v", err)
	}
}

func TestJoinAlignMapping(t *testing.T) {
	for _, a := range []string{"left", "center", "right", "top", "bottom", "middle"} {
		o := Options{
			Text:       []string{"A", "B"},
			Align:      a,
			Horizontal: true,
		}
		if err := o.Run(); err != nil {
			t.Errorf("Join with align=%q failed: %v", a, err)
		}
	}
}
