package omniv2xml

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jf-tech/omniparser/errs"
)

func TestIsErrNodeReadingFailed(t *testing.T) {
	assert.True(t, IsErrNodeReadingFailed(ErrNodeReadingFailed("test")))
	assert.Equal(t, "test", ErrNodeReadingFailed("test").Error())
	assert.False(t, IsErrNodeReadingFailed(errors.New("test")))
}

func TestReader_Read_Success(t *testing.T) {
	r, err := NewReader(
		"test-input",
		strings.NewReader(`
			<Root>
				<Node>1</Node>
				<Node>2</Node>
			</Root>`),
		"Root/Node[. != '2']")
	assert.NoError(t, err)
	assert.Equal(t, 1, r.lineNumber())

	n, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, "1", n.InnerText())
	// xml.Decoder seems to keeps line at the end of whatever inside an element closing tag.
	assert.Equal(t, 3, r.lineNumber())

	n, err = r.Read()
	assert.Error(t, err)
	assert.Equal(t, errs.ErrEOF, err)
	assert.Nil(t, n)
}

func TestReader_Read_InvalidXML(t *testing.T) {
	r, err := NewReader(
		"test-input",
		strings.NewReader(`
			<Root>
				<Node>1<Node>
				<Node>2</Node>
			</Root>`),
		"Root/Node[. != '2']")
	assert.NoError(t, err)
	assert.Equal(t, 1, r.lineNumber())

	n, err := r.Read()
	assert.Error(t, err)
	assert.True(t, IsErrNodeReadingFailed(err))
	assert.Equal(t,
		`input 'test-input' near line 5: XML syntax error on line 5: element <Node> closed by </Root>`,
		err.Error())
	assert.Nil(t, n)
}

func TestReader_FmtErr(t *testing.T) {
	r, err := NewReader("test-input", strings.NewReader(""), "Root/Node")
	assert.NoError(t, err)
	err = r.FmtErr("golang is %s", "fun")
	assert.Error(t, err)
	assert.Equal(t, `input 'test-input' near line 1: golang is fun`, err.Error())
}

func TestReader_IsContinuableError(t *testing.T) {
	r, err := NewReader("test", strings.NewReader(""), "Root/Node")
	assert.NoError(t, err)
	assert.False(t, r.IsContinuableError(errs.ErrEOF))
	assert.False(t, r.IsContinuableError(ErrNodeReadingFailed("failure")))
	assert.True(t, r.IsContinuableError(errs.ErrTransformFailed("failure")))
	assert.True(t, r.IsContinuableError(errors.New("failure")))
}

func TestNewReader_InvalidXPath(t *testing.T) {
	r, err := NewReader("test-input", strings.NewReader(""), "[not-valid")
	assert.Error(t, err)
	assert.Equal(t,
		`invalid streamElementXPath '[not-valid', err: expression must evaluate to a node-set`,
		err.Error())
	assert.Nil(t, r)
}
