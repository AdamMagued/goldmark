package extension

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/testutil"
)

func TestTaskList(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			TaskList,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/tasklist.txt", t, testutil.ParseCliCaseArg()...)
}

func TestTaskListWithWrapLabelFalse(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			NewTaskList(WithTaskListWrapLabel(false)),
		),
	)
	source := []byte("- [ ] foo\n- [x] bar\n")
	var buf bytes.Buffer
	if err := markdown.Convert(source, &buf); err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n<li><input disabled=\"\" type=\"checkbox\"> foo</li>\n<li><input checked=\"\" disabled=\"\" type=\"checkbox\"> bar</li>\n</ul>\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestTaskListAccessibleLabel(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			TaskList,
		),
	)
	source := []byte("- [x] Basic Test\n- [ ] More Tests\n  - [x] View\n  - [ ] Smell\n")
	var buf bytes.Buffer
	if err := markdown.Convert(source, &buf); err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n" +
		"<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> Basic Test</label></li>\n" +
		"<li><label><input disabled=\"\" type=\"checkbox\"> More Tests</label>\n" +
		"<ul>\n" +
		"<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> View</label></li>\n" +
		"<li><label><input disabled=\"\" type=\"checkbox\"> Smell</label></li>\n" +
		"</ul>\n" +
		"</li>\n" +
		"</ul>\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestTaskListFormatting(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			TaskList,
		),
	)
	source := []byte("- [x] **bold** text and `code` and [link](https://example.com)\n")
	var buf bytes.Buffer
	if err := markdown.Convert(source, &buf); err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> <strong>bold</strong> text and <code>code</code> and <a href=\"https://example.com\">link</a></label></li>\n</ul>\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestTaskListEmpty(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			TaskList,
		),
	)
	source := []byte("- [x]\n- [ ]\n")
	var buf bytes.Buffer
	if err := markdown.Convert(source, &buf); err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"></label></li>\n<li><label><input disabled=\"\" type=\"checkbox\"></label></li>\n</ul>\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestTaskListLoose(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			TaskList,
		),
	)
	source := []byte("- [x] item 1\n\n- [ ] item 2\n")
	var buf bytes.Buffer
	if err := markdown.Convert(source, &buf); err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n<li>\n<p><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> item 1</label></p>\n</li>\n<li>\n<p><label><input disabled=\"\" type=\"checkbox\"> item 2</label></p>\n</li>\n</ul>\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestTaskListXHTML(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithXHTML(),
		),
		goldmark.WithExtensions(
			TaskList,
		),
	)
	source := []byte("- [x] foo\n")
	var buf bytes.Buffer
	if err := markdown.Convert(source, &buf); err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\" /> foo</label></li>\n</ul>\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestTaskListMixed(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			TaskList,
		),
	)
	source := []byte("- regular item\n- [x] task item\n- another regular item\n")
	var buf bytes.Buffer
	if err := markdown.Convert(source, &buf); err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n<li>regular item</li>\n<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> task item</label></li>\n<li>another regular item</li>\n</ul>\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}
