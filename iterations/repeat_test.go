package iterations

import (
	"fmt"
	"testing"
)

func TestRepeat(t* testing.T) {
	repeated := Repeat("a", 6)
	expected := "aaaaaa"

	if repeated != expected {
		t.Errorf("expected %q but got %q", expected, repeated)
	}
} 

func ExampleRepeat() {
	repeat := Repeat("a", 7)
	fmt.Println(repeat)
	// Output: aaaaaaa
}

func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a", 6)
	}
}