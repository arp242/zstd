// Package zxml provides functions for working with XML.
package zxml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Format the XML document xmlDoc.
//
// This is not especially efficient for (very) large XML documents.
func Format(xmlDoc []byte, prefix, indent string) ([]byte, error) {
	// Read all tokens in a slice; this makes it possible to peek to the next
	// one later on. Could make this more efficient, but meh.
	var (
		dec        = xml.NewDecoder(bytes.NewReader(xmlDoc))
		namespaces []string
		tokens     = make([]xml.Token, 0, 16)
	)
	for {
		token, err := dec.Token()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("format: %w", err)
			}
			break
		}
		// Skip empty character nodes; this is often just whitespace between
		// elements.
		if t, ok := token.(xml.CharData); ok {
			if len(bytes.TrimSpace(t)) == 0 {
				continue
			}
		}
		// Keep track if have xmlns is set; we don't really want to add the full
		// namespace for everything later on, which looks more than a little
		// ugly.
		if t, ok := token.(xml.StartElement); ok {
			for _, a := range t.Attr {
				if a.Name.Local == "xmlns" {
					namespaces = append(namespaces, a.Value)
				}
			}
		}
		tokens = append(tokens, xml.CopyToken(token))
	}

	var (
		indentLevel = 0
		b           = new(bytes.Buffer)
		doIndent    = true
		writeIndent = func() {
			for i := 0; i < indentLevel; i++ {
				b.WriteString(indent)
			}
		}
		writeNS = func(ns string) {
			if ns == "" {
				return
			}
			for _, n := range namespaces {
				if ns == n {
					return
				}
			}
			if i := strings.LastIndex(ns, "/"); i > -1 {
				ns = ns[i+1:]
			}
			if i := strings.LastIndex(ns, ":"); i > -1 {
				ns = ns[i+1:]
			}
			if i := strings.LastIndex(ns, "-"); i > -1 {
				ns = ns[:i]
			}

			b.WriteString(ns)
			b.WriteByte(':')
		}
	)
	b.Grow(len(xmlDoc))
	for i, token := range tokens {
		switch t := token.(type) {
		case xml.ProcInst:
			b.WriteString("<?")
			b.WriteString(t.Target)
			if len(t.Inst) > 0 {
				b.WriteByte(' ')
				b.Write(t.Inst)
			}
			b.WriteString("?>\n")
		case xml.Comment:
			b.WriteString("<!-- ")
			b.Write(t)
			b.WriteString(" -->")
		case xml.Directive:
			b.WriteString("<!")
			b.Write(t)
			b.WriteString(">")

		case xml.CharData:
			if doIndent {
				writeIndent()
			}
			b.Write(t)
			if doIndent {
				b.WriteByte('\n')
			}
		case xml.StartElement:
			writeIndent()

			b.WriteByte('<')
			writeNS(t.Name.Space)
			b.WriteString(t.Name.Local)
			if len(t.Attr) > 0 {
				for _, a := range t.Attr {
					b.WriteByte(' ')
					writeNS(a.Name.Space)
					b.WriteString(a.Name.Local)
					b.WriteString(`="`)
					b.WriteString(a.Value)
					b.WriteByte('"')
				}
			}
			b.WriteByte('>')
			indentLevel++

			// Print newline only if next element is a start element.
			if _, ok := tokens[i+1].(xml.StartElement); ok {
				doIndent = true
				b.WriteByte('\n')
			} else {
				doIndent = false
			}
		case xml.EndElement:
			indentLevel--
			if doIndent {
				writeIndent()
			}
			b.WriteString("</")
			writeNS(t.Name.Space)
			b.WriteString(t.Name.Local)
			b.WriteString(">\n")
			doIndent = true
		}
	}
	if b.Len() == 0 {
		return b.Bytes(), nil
	}
	return b.Bytes()[:b.Len()-1], nil // Eat last newline
}
