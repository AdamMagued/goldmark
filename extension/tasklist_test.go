package extension

import (
	"testing"

	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/testutil"
)

func TestTaskList(t *testing.T) {
	markdown := testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithExtensions(NewTaskListItemParser()),
		),
		html.New(
			html.WithUnsafe(),
			html.WithExtensions(NewTaskListItemHTMLRenderer()),
		),
	)
	testutil.DoTestCaseFile(markdown, "testdata/tasklist.txt", t, testutil.ParseCliCaseArg()...)
}

func TestTaskListWrapLabel(t *testing.T) {
	markdown := testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithExtensions(NewTaskListItemParser()),
		),
		html.New(
			html.WithUnsafe(),
			html.WithExtensions(NewTaskListItemHTMLRenderer(WithWrapLabel())),
		),
	)

	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name:   "basic checked and unchecked items",
			source: "- [ ] foo\n- [x] bar\n",
			expected: "<ul>\n" +
				"<li><label><input disabled=\"\" type=\"checkbox\"> foo</label></li>\n" +
				"<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> bar</label></li>\n" +
				"</ul>\n",
		},
		{
			name:   "nested task lists matching issue 609",
			source: "- [x] Basic Test\n- [ ] More Tests\n  - [x] View\n  - [ ] Smell\n",
			expected: "<ul>\n" +
				"<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> Basic Test</label></li>\n" +
				"<li><label><input disabled=\"\" type=\"checkbox\"> More Tests</label>\n" +
				"<ul>\n" +
				"<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> View</label></li>\n" +
				"<li><label><input disabled=\"\" type=\"checkbox\"> Smell</label></li>\n" +
				"</ul>\n" +
				"</li>\n" +
				"</ul>\n",
		},
		{
			name:   "inline formatting in task item",
			source: "- [x] **bold** text and `code` and [link](https://example.com)\n",
			expected: "<ul>\n" +
				"<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> <strong>bold</strong> text and <code>code</code> and <a href=\"https://example.com\">link</a></label></li>\n" +
				"</ul>\n",
		},
		{
			name:   "empty task item",
			source: "- [x]\n- [ ]\n",
			expected: "<ul>\n" +
				"<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"></label></li>\n" +
				"<li><label><input disabled=\"\" type=\"checkbox\"></label></li>\n" +
				"</ul>\n",
		},
		{
			name:   "loose task list",
			source: "- [x] item 1\n\n- [ ] item 2\n",
			expected: "<ul>\n" +
				"<li>\n" +
				"<p><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> item 1</label></p>\n" +
				"</li>\n" +
				"<li>\n" +
				"<p><label><input disabled=\"\" type=\"checkbox\"> item 2</label></p>\n" +
				"</li>\n" +
				"</ul>\n",
		},
		{
			name:   "mixed regular and task items",
			source: "- regular item\n- [x] task item\n- another regular item\n",
			expected: "<ul>\n" +
				"<li>regular item</li>\n" +
				"<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> task item</label></li>\n" +
				"<li>another regular item</li>\n" +
				"</ul>\n",
		},
		{
			name:   "task item with multiple paragraphs",
			source: "- [x] first paragraph\n\n  second paragraph\n",
			expected: "<ul>\n" +
				"<li>\n" +
				"<p><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> first paragraph</label></p>\n" +
				"<p>second paragraph</p>\n" +
				"</li>\n" +
				"</ul>\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := markdown(tc.source)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("got %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestTaskListWrapLabelFalse(t *testing.T) {
	markdown := testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithExtensions(NewTaskListItemParser()),
		),
		html.New(
			html.WithUnsafe(),
			html.WithExtensions(NewTaskListItemHTMLRenderer(WithWrapLabel(false))),
		),
	)

	source := "- [ ] foo\n- [x] bar\n"
	expected := "<ul>\n<li><input disabled=\"\" type=\"checkbox\"> foo</li>\n<li><input checked=\"\" disabled=\"\" type=\"checkbox\"> bar</li>\n</ul>\n"
	got, err := markdown(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestTaskListWrapLabelXHTML(t *testing.T) {
	markdown := testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithExtensions(NewTaskListItemParser()),
		),
		html.New(
			html.WithUnsafe(),
			html.WithExtensions(NewTaskListItemHTMLRenderer(WithWrapLabel(), WithXHTML())),
		),
	)

	source := "- [x] foo\n"
	expected := "<ul>\n<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\" /> foo</label></li>\n</ul>\n"
	got, err := markdown(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}
