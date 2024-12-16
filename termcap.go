package tea

import (
	"bytes"
	"encoding/hex"
	"strings"
)

// TermcapMsg represents a Termcap response event. Termcap responses are
// generated by the terminal in response to RequestTermcap (XTGETTCAP)
// requests.
//
// See: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
type TermcapMsg string

func parseTermcap(data []byte) TermcapMsg {
	// XTGETTCAP
	if len(data) == 0 {
		return TermcapMsg("")
	}

	var tc strings.Builder
	split := bytes.Split(data, []byte{';'})
	for _, s := range split {
		parts := bytes.SplitN(s, []byte{'='}, 2)
		if len(parts) == 0 {
			return TermcapMsg("")
		}

		name, err := hex.DecodeString(string(parts[0]))
		if err != nil || len(name) == 0 {
			continue
		}

		var value []byte
		if len(parts) > 1 {
			value, err = hex.DecodeString(string(parts[1]))
			if err != nil {
				continue
			}
		}

		if tc.Len() > 0 {
			tc.WriteByte(';')
		}
		tc.WriteString(string(name))
		if len(value) > 0 {
			tc.WriteByte('=')
			tc.WriteString(string(value))
		}
	}

	return TermcapMsg(tc.String())
}
