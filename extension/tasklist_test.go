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

func TestTaskListWrapLabelFalse(t *testing.T) {
	markdown := testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithExtensions(NewTaskListItemParser()),
		),
		html.New(
			html.WithUnsafe(),
			html.WithExtensions(NewTaskListItemHTMLRenderer(WithTaskListItemWrapLabel(false))),
		),
	)
	source := "- [ ] foo\n- [x] bar\n"
	got, err := markdown(source)
	if err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n<li><input disabled=\"\" type=\"checkbox\"> foo</li>\n<li><input checked=\"\" disabled=\"\" type=\"checkbox\"> bar</li>\n</ul>\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestTaskListAccessibleLabel(t *testing.T) {
	markdown := testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithExtensions(NewTaskListItemParser()),
		),
		html.New(
			html.WithUnsafe(),
			html.WithExtensions(NewTaskListItemHTMLRenderer()),
		),
	)
	source := "- [x] Basic Test\n- [ ] More Tests\n  - [x] View\n  - [ ] Smell\n"
	got, err := markdown(source)
	if err != nil {
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
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestTaskListFormatting(t *testing.T) {
	markdown := testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithExtensions(NewTaskListItemParser()),
		),
		html.New(
			html.WithUnsafe(),
			html.WithExtensions(NewTaskListItemHTMLRenderer()),
		),
	)
	source := "- [x] **bold** text and `code` and [link](https://example.com)\n"
	got, err := markdown(source)
	if err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> <strong>bold</strong> text and <code>code</code> and <a href=\"https://example.com\">link</a></label></li>\n</ul>\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestTaskListEmpty(t *testing.T) {
	markdown := testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithExtensions(NewTaskListItemParser()),
		),
		html.New(
			html.WithUnsafe(),
			html.WithExtensions(NewTaskListItemHTMLRenderer()),
		),
	)
	source := "- [x]\n- [ ]\n"
	got, err := markdown(source)
	if err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"></label></li>\n<li><label><input disabled=\"\" type=\"checkbox\"></label></li>\n</ul>\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestTaskListLoose(t *testing.T) {
	markdown := testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithExtensions(NewTaskListItemParser()),
		),
		html.New(
			html.WithUnsafe(),
			html.WithExtensions(NewTaskListItemHTMLRenderer()),
		),
	)
	source := "- [x] item 1\n\n- [ ] item 2\n"
	got, err := markdown(source)
	if err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n<li>\n<p><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> item 1</label></p>\n</li>\n<li>\n<p><label><input disabled=\"\" type=\"checkbox\"> item 2</label></p>\n</li>\n</ul>\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestTaskListXHTML(t *testing.T) {
	markdown := testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithExtensions(NewTaskListItemParser()),
		),
		html.New(
			html.WithXHTML(),
			html.WithExtensions(NewTaskListItemHTMLRenderer(WithXHTML())),
		),
	)
	source := "- [x] foo\n"
	got, err := markdown(source)
	if err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\" /> foo</label></li>\n</ul>\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestTaskListMixed(t *testing.T) {
	markdown := testutil.NewMarkdownToStringFunc(
		parser.New(
			parser.WithExtensions(NewTaskListItemParser()),
		),
		html.New(
			html.WithUnsafe(),
			html.WithExtensions(NewTaskListItemHTMLRenderer()),
		),
	)
	source := "- regular item\n- [x] task item\n- another regular item\n"
	got, err := markdown(source)
	if err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n<li>regular item</li>\n<li><label><input checked=\"\" disabled=\"\" type=\"checkbox\"> task item</label></li>\n<li>another regular item</li>\n</ul>\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}
