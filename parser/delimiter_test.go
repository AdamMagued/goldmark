package parser

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yuin/goldmark/v2/ast"
)

func TestDelimiterOpenerCountsInline(t *testing.T) {
	var counts delimiterOpenerCounts

	if got := counts.get('*'); got != 0 {
		t.Errorf("expected 0 for unadded '*', got %d", got)
	}

	counts.add('*', 1)
	counts.add('_', 2)
	counts.add('~', 3)

	if got := counts.get('*'); got != 1 {
		t.Errorf("expected 1 for '*', got %d", got)
	}
	if got := counts.get('_'); got != 2 {
		t.Errorf("expected 2 for '_', got %d", got)
	}
	if got := counts.get('~'); got != 3 {
		t.Errorf("expected 3 for '~', got %d", got)
	}
	if got := counts.get('='); got != 0 {
		t.Errorf("expected 0 for '=', got %d", got)
	}

	counts.add('*', -1)
	if got := counts.get('*'); got != 0 {
		t.Errorf("expected 0 after decrementing '*', got %d", got)
	}
	if counts.overflow != nil {
		t.Errorf("expected nil overflow map for inline counts, got non-nil")
	}
}

func TestDelimiterOpenerCountsBeyondInlineCapacity(t *testing.T) {
	var counts delimiterOpenerCounts
	chars := []byte("*_~=^+")
	for _, c := range chars {
		counts.add(c, 1)
	}
	counts.add('*', 1)
	counts.add('+', 1)
	counts.add('_', -1)
	counts.add('^', -1)

	expectedMap := map[byte]int{
		'*': 2,
		'_': 0,
		'~': 1,
		'=': 1,
		'^': 0,
		'+': 2,
		'!': 0,
	}

	for c, expected := range expectedMap {
		if actual := counts.get(c); actual != expected {
			t.Errorf("%q: expected %d, got %d", c, expected, actual)
		}
	}

	if counts.overflow == nil {
		t.Errorf("expected non-nil overflow map after exceeding inline capacity")
	}
}

func TestDelimiterOpenerCountsTrackingInContext(t *testing.T) {
	pc := newContext(nil, false).(*parseContext)
	parent := ast.NewParagraph()

	d1 := NewDelimiter(true, false, 1, '*', defaultEmphasisDelimiterProcessor)
	d2 := NewDelimiter(true, false, 1, '*', defaultEmphasisDelimiterProcessor)
	d3 := NewDelimiter(false, true, 1, '*', defaultEmphasisDelimiterProcessor)
	d4 := NewDelimiter(true, false, 1, '_', defaultEmphasisDelimiterProcessor)

	parent.AppendChild(d1)
	parent.AppendChild(d2)
	parent.AppendChild(d3)
	parent.AppendChild(d4)

	pc.PushDelimiter(d1)
	pc.PushDelimiter(d2)
	pc.PushDelimiter(d3)
	pc.PushDelimiter(d4)

	if !d1.countedAsOpener || !d2.countedAsOpener {
		t.Errorf("expected d1 and d2 to have countedAsOpener=true")
	}
	if d3.countedAsOpener {
		t.Errorf("expected d3 (canOpen=false) to have countedAsOpener=false")
	}
	if !d4.countedAsOpener {
		t.Errorf("expected d4 to have countedAsOpener=true")
	}

	if got := pc.openerCounts.get('*'); got != 2 {
		t.Errorf("expected 2 openers for '*', got %d", got)
	}
	if got := pc.openerCounts.get('_'); got != 1 {
		t.Errorf("expected 1 opener for '_', got %d", got)
	}

	pc.RemoveDelimiter(d1)
	if d1.countedAsOpener {
		t.Errorf("expected d1.countedAsOpener=false after RemoveDelimiter")
	}
	if got := pc.openerCounts.get('*'); got != 1 {
		t.Errorf("expected 1 opener for '*' after removing d1, got %d", got)
	}

	pc.RemoveDelimiter(d3) // Was not an opener, should not affect counts
	if got := pc.openerCounts.get('*'); got != 1 {
		t.Errorf("expected 1 opener for '*' after removing non-opener d3, got %d", got)
	}

	pc.RemoveDelimiter(d2)
	if got := pc.openerCounts.get('*'); got != 0 {
		t.Errorf("expected 0 openers for '*' after removing d2, got %d", got)
	}
}

func TestProcessDelimitersUnmatchedDelimitersPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance regression test in short mode")
	}

	p := New()
	// Long run of unmatched closing delimiters (e.g. underscores in URLs or after words)
	n := 40000
	source := []byte(strings.Repeat("a_ ", n))

	start := time.Now()
	doc := p.Parse(source)
	elapsed := time.Since(start)

	if elapsed > 3*time.Second {
		t.Fatalf("parsing %d unmatched delimiters took %v, expected < 3s (O(N^2) backward scan detected)", n, elapsed)
	}

	if doc == nil {
		t.Fatalf("expected non-nil parsed document")
	}
}

func TestProcessDelimitersMixedDelimitersPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance regression test in short mode")
	}

	p := New()
	// Many unmatched delimiters followed by matched emphasis spans
	n := 40000
	source := []byte(strings.Repeat("a_ ", n) + strings.Repeat("b*c ", n))

	start := time.Now()
	doc := p.Parse(source)
	elapsed := time.Since(start)

	if elapsed > 3*time.Second {
		t.Fatalf("parsing mixed unmatched and matched delimiters took %v, expected < 3s", elapsed)
	}

	if doc == nil {
		t.Fatalf("expected non-nil parsed document")
	}
}

func TestProcessDelimitersCorrectnessWithUnmatchedDelimiters(t *testing.T) {
	p := New()
	source := []byte("foo_ bar_ *emphasized* baz_ _italic_")
	doc := p.Parse(source)

	var emphasisCount int
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if _, ok := n.(*ast.Emphasis); ok && entering {
			emphasisCount++
		}
		return ast.WalkContinue, nil
	})

	if emphasisCount != 2 {
		t.Errorf("expected 2 emphasis nodes, got %d", emphasisCount)
	}
}

func BenchmarkProcessDelimitersUnmatchedDelimiters(b *testing.B) {
	p := New()
	source := []byte(strings.Repeat("http://example.com/item_one_two_three/ ", 200))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = p.Parse(source)
	}
}

func BenchmarkProcessDelimitersMatchedDelimiters(b *testing.B) {
	p := New()
	var buf bytes.Buffer
	for i := 0; i < 200; i++ {
		buf.WriteString("*emphasized text* and _italic text_ ")
	}
	source := buf.Bytes()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = p.Parse(source)
	}
}
