package flamewriter

import (
	"encoding/json"
	"io"
)

// JSONWriter write flame data in JSON format
type JSONWriter struct {
	w io.Writer
}

// NewJSONWriter creates a new JSONWriter
func NewJSONWriter(w io.Writer) *JSONWriter {
	return &JSONWriter{
		w: w,
	}
}

func toJSONObject(root *Record) map[string]interface{} {
	children := make([]map[string]interface{}, 0, len(root.children))
	for _, child := range root.children {
		children = append(children, toJSONObject(child))
	}
	return map[string]interface{}{
		"name":     root.name,
		"value":    root.value,
		"children": children,
	}
}

// Write writes flame data
func (w *JSONWriter) Write(root *Record) error {
	return json.NewEncoder(w.w).Encode(toJSONObject(root))
}
