package log

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTemplateFormatter_Format(t *testing.T) {
	entry := &Entry{
		Logger: nil,
		Data: Fields{
			"my_key": "my_value",
		},
		Time:    time.Time{},
		Level:   2,
		Message: "hello",
		Buffer:  nil,
	}

	f, err := NewFormatterTemplateBuilder().
		WithTimestampFormat("2006/01/02 15:04:05").
		Build()
	require.NoError(t, err)
	data, err := f.Format(entry)
	require.NoError(t, err)
	assert.Equal(t, "0001/01/01 00:00:00 ERROR: hello DETAILS {\"my_key\":\"my_value\"}\n", string(data))
}
