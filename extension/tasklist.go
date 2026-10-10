package extension

import (
	"regexp"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var taskListRegexp = regexp.MustCompile(`^\[([\sxX])\]\s*`)

type taskCheckBoxParser struct {
}

var defaultTaskCheckBoxParser = &taskCheckBoxParser{}

// NewTaskCheckBoxParser returns a new InlineParser that can parse
// checkboxes in list items.
// This parser must take precedence over the parser.LinkParser.
func NewTaskCheckBoxParser() parser.InlineParser {
	return defaultTaskCheckBoxParser
}

func (s *taskCheckBoxParser) Trigger() []byte {
	return []byte{'['}
}

func (s *taskCheckBoxParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	// Given AST structure must be like
	// - List
	//   - ListItem         : parent.Parent
	//     - TextBlock      : parent
	//       (current line)
	if parent.Parent() == nil || parent.Parent().FirstChild() != parent {
		return nil
	}

	if parent.HasChildren() {
		return nil
	}
	if _, ok := parent.Parent().(*gast.ListItem); !ok {
		return nil
	}
	line, _ := block.PeekLine()
	m := taskListRegexp.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}
	value := line[m[2]:m[3]][0]
	block.Advance(m[1])
	checked := value == 'x' || value == 'X'
	return ast.NewTaskCheckBox(checked)
}

func (s *taskCheckBoxParser) CloseBlock(parent gast.Node, pc parser.Context) {
	// nothing to do
}

type taskCheckBoxASTTransformer struct {
}

var defaultTaskCheckBoxASTTransformer = &taskCheckBoxASTTransformer{}

// NewTaskCheckBoxASTTransformer returns a new parser.ASTTransformer for task lists.
func NewTaskCheckBoxASTTransformer() parser.ASTTransformer {
	return defaultTaskCheckBoxASTTransformer
}

func (a *taskCheckBoxASTTransformer) Transform(node *gast.Document, reader text.Reader, pc parser.Context) {
	var checkBoxes []*ast.TaskCheckBox
	_ = gast.Walk(node, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if entering {
			if checkBox, ok := n.(*ast.TaskCheckBox); ok {
				checkBoxes = append(checkBoxes, checkBox)
			}
		}
		return gast.WalkContinue, nil
	})
	for _, checkBox := range checkBoxes {
		parent := checkBox.Parent()
		if parent != nil && parent.FirstChild() == checkBox {
			for c := checkBox.NextSibling(); c != nil; {
				next := c.NextSibling()
				parent.RemoveChild(parent, c)
				checkBox.AppendChild(checkBox, c)
				c = next
			}
		}
	}
}

// TaskListConfig struct holds options for the task list extension.
type TaskListConfig struct {
	html.Config

	// WrapLabel indicates whether checkboxes and task list item texts are wrapped in <label>.
	WrapLabel bool
}

// NewTaskListConfig returns a new TaskListConfig with defaults.
func NewTaskListConfig() TaskListConfig {
	return TaskListConfig{
		Config:    html.NewConfig(),
		WrapLabel: true,
	}
}

// SetOption implements renderer.SetOptioner.
func (c *TaskListConfig) SetOption(name renderer.OptionName, value any) {
	switch name {
	case optTaskListWrapLabel:
		c.WrapLabel = value.(bool)
	default:
		c.Config.SetOption(name, value)
	}
}

// TaskListOption interface is a functional option interface for the task list extension.
type TaskListOption interface {
	renderer.Option
	SetTaskListOption(*TaskListConfig)
}

const optTaskListWrapLabel renderer.OptionName = "TaskListWrapLabel"

type withTaskListWrapLabel struct {
	value bool
}

func (o *withTaskListWrapLabel) SetConfig(c *renderer.Config) {
	c.Options[optTaskListWrapLabel] = o.value
}

func (o *withTaskListWrapLabel) SetTaskListOption(c *TaskListConfig) {
	c.WrapLabel = o.value
}

// WithTaskListWrapLabel is a functional option that indicates whether checkboxes and
// task list item texts are wrapped in <label> elements to satisfy WCAG 4.1.2.
func WithTaskListWrapLabel(wrap bool) TaskListOption {
	return &withTaskListWrapLabel{wrap}
}

// WithWrapLabel is an alias for WithTaskListWrapLabel.
func WithWrapLabel(wrap bool) TaskListOption {
	return WithTaskListWrapLabel(wrap)
}

type withTaskListHTMLOptions struct {
	value []html.Option
}

func (o *withTaskListHTMLOptions) SetConfig(c *renderer.Config) {
	if o.value != nil {
		for _, v := range o.value {
			if ro, ok := v.(renderer.Option); ok {
				ro.SetConfig(c)
			}
		}
	}
}

func (o *withTaskListHTMLOptions) SetTaskListOption(c *TaskListConfig) {
	if o.value != nil {
		for _, v := range o.value {
			v.SetHTMLOption(&c.Config)
		}
	}
}

// WithTaskListHTMLOptions is a functional option that wraps goldmark HTMLRenderer options.
func WithTaskListHTMLOptions(opts ...html.Option) TaskListOption {
	return &withTaskListHTMLOptions{opts}
}

// TaskCheckBoxHTMLRenderer is a renderer.NodeRenderer implementation that
// renders checkboxes in list items.
type TaskCheckBoxHTMLRenderer struct {
	TaskListConfig
}

// NewTaskCheckBoxHTMLRenderer returns a new TaskCheckBoxHTMLRenderer.
func NewTaskCheckBoxHTMLRenderer(opts ...TaskListOption) renderer.NodeRenderer {
	r := &TaskCheckBoxHTMLRenderer{
		TaskListConfig: NewTaskListConfig(),
	}
	for _, opt := range opts {
		opt.SetTaskListOption(&r.TaskListConfig)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *TaskCheckBoxHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindTaskCheckBox, r.renderTaskCheckBox)
}

func (r *TaskCheckBoxHTMLRenderer) renderTaskCheckBox(
	w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*ast.TaskCheckBox)

	if entering {
		if r.WrapLabel {
			_, _ = w.WriteString("<label>")
		}
		if n.IsChecked {
			_, _ = w.WriteString(`<input checked="" disabled="" type="checkbox"`)
		} else {
			_, _ = w.WriteString(`<input disabled="" type="checkbox"`)
		}
		if r.XHTML {
			if n.HasChildren() || !r.WrapLabel {
				_, _ = w.WriteString(" /> ")
			} else {
				_, _ = w.WriteString(" />")
			}
		} else {
			if n.HasChildren() || !r.WrapLabel {
				_, _ = w.WriteString("> ")
			} else {
				_, _ = w.WriteString(">")
			}
		}
	} else {
		if r.WrapLabel {
			_, _ = w.WriteString("</label>")
		}
	}
	return gast.WalkContinue, nil
}

type taskList struct {
	options []TaskListOption
}

// TaskList is an extension that allow you to use GFM task lists.
var TaskList = NewTaskList()

// NewTaskList returns a new task list extension with given options.
func NewTaskList(opts ...TaskListOption) goldmark.Extender {
	return &taskList{
		options: opts,
	}
}

func (e *taskList) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewTaskCheckBoxParser(), 0),
		),
		parser.WithASTTransformers(
			util.Prioritized(NewTaskCheckBoxASTTransformer(), 500),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewTaskCheckBoxHTMLRenderer(e.options...), 500),
		),
	)
}
