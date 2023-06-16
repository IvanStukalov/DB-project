// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonFda65e30DecodeGithubComIvanStukalovDBProjectInternalModels(in *jlexer.Lexer, out *ErrMsg) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "msg":
			out.Msg = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonFda65e30EncodeGithubComIvanStukalovDBProjectInternalModels(out *jwriter.Writer, in ErrMsg) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"msg\":"
		out.RawString(prefix[1:])
		out.String(string(in.Msg))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ErrMsg) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonFda65e30EncodeGithubComIvanStukalovDBProjectInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ErrMsg) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonFda65e30EncodeGithubComIvanStukalovDBProjectInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ErrMsg) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonFda65e30DecodeGithubComIvanStukalovDBProjectInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ErrMsg) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonFda65e30DecodeGithubComIvanStukalovDBProjectInternalModels(l, v)
}
