// preprocessor_test.go
package preprocessor

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestDefineConstants(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name: "Basic defines",
			input: []string{
				"#define T_EOF 0x0000",
				"#define T_SPACE 0x0001",
				"#define T_TAB 0x0002",
				"# 1 \"test.c\"",
				"int main() { return 0; }",
			},
			expected: "const (\n\tT_EOF = 0x0000\n\tT_SPACE = 0x0001\n\tT_TAB = 0x0002\n)\n# 1 \"test.c\"\nint main() { return 0; }",
		},
		{
			name: "No defines",
			input: []string{
				"# 1 \"test.c\"",
				"int main() { return 0; }",
			},
			expected: "# 1 \"test.c\"\nint main() { return 0; }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := ioutil.TempFile("", "test*.c")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())
			_, err = tmpFile.WriteString(strings.Join(tt.input, "\n"))
			if err != nil {
				t.Fatal(err)
			}
			tmpFile.Close()

			f, err := NewFilePP([]string{tmpFile.Name()}, []string{}, false)
			if err != nil {
				t.Errorf("NewFilePP() error = %v", err)
				return
			}
			got := string(f.GetSource())
			if got != tt.expected {
				t.Errorf("GetSource() = %q, want %q", got, tt.expected)
			}
		})
	}
}
