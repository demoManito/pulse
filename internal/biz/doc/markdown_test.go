package doc

import (
	"testing"
)

func TestHTMLToMarkdown(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    string
		wantErr bool
	}{
		{
			name: "paragraph",
			html: "<p>hello world</p>",
			want: "hello world",
		},
		{
			name: "heading",
			html: "<h1>Title</h1><p>body</p>",
			want: "# Title\n\nbody",
		},
		{
			name: "bold and italic",
			html: "<p><strong>bold</strong> and <em>italic</em></p>",
			want: "**bold** and *italic*",
		},
		{
			name: "unordered list",
			html: "<ul><li>a</li><li>b</li></ul>",
			want: "- a\n- b",
		},
		{
			name: "link",
			html: `<p><a href="https://example.com">click</a></p>`,
			want: "[click](https://example.com)",
		},
		{
			name: "nested heading and list",
			html: "<h2>Section</h2><ul><li>item1</li><li>item2</li></ul>",
			want: "## Section\n\n- item1\n- item2",
		},
		{
			name: "empty input",
			html: "",
			want: "",
		},
		{
			name: "whitespace only",
			html: "   \n\t  ",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HTMLToMarkdown(tt.html)
			if (err != nil) != tt.wantErr {
				t.Fatalf("HTMLToMarkdown() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("HTMLToMarkdown() = %q, want %q", got, tt.want)
			}
		})
	}
}
