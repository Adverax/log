package rex

import (
	"github.com/adverax/log"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParser(t *testing.T) {
	e := NewEngine()
	f, err := NewFrame(`(?P<time>\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}) (?P<level>\w+) (?P<message>.*)`)
	e.AddFrame(f)
	require.NoError(t, err)
	entry := log.NewEntry(nil)
	_, err = e.Parse([]byte(`2025/01/08 14:49:34 DEBUG #55300a10-ec6c-41b9-9866-3784983b7960: ASPECT.BOOTSTRAP: aspect.options`), entry)
	require.NoError(t, err)
}
