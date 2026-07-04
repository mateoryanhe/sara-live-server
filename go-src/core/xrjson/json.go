package xrjson

import "github.com/bytedance/sonic"

var api = sonic.Config{SortMapKeys: true}.Froze()

func Marshal(v any) ([]byte, error) {
	return api.Marshal(v)
}

func MustMarshal(v any) []byte {
	data, err := api.Marshal(v)
	if err != nil {
		return []byte("{}")
	}
	return data
}

func Unmarshal(data []byte, v any) error {
	return api.Unmarshal(data, v)
}

func MarshalIndent(v any) ([]byte, error) {
	return api.MarshalIndent(v, "", " ")
}

func MustMarshalIndent(v any) []byte {
	data, err := api.MarshalIndent(v, "", " ")
	if err != nil {
		return []byte("{}")
	}
	return data
}
