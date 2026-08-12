package randomnick

import "testing"

func TestBuiltinEnglishNicknames(t *testing.T) {
	list := buildBuiltinEnglishNicknames()
	if len(list) != 1000 {
		t.Fatalf("expected 1000 builtin nicknames, got %d", len(list))
	}
	seen := make(map[string]struct{}, len(list))
	for _, name := range list {
		if _, ok := seen[name]; ok {
			t.Fatalf("duplicate builtin nickname: %s", name)
		}
		seen[name] = struct{}{}
	}
}

func TestParseImportText(t *testing.T) {
	names := ParseImportText("  Alpha \nBeta\nAlpha\n\n")
	if len(names) != 2 || names[0] != "Alpha" || names[1] != "Beta" {
		t.Fatalf("unexpected parse result: %#v", names)
	}
}

func TestPickRandomUsesMemory(t *testing.T) {
	list := buildBuiltinEnglishNicknames()
	poolCache.Store(&nicknamePoolSnapshot{
		useDB:    false,
		byLang:   map[uint8][]string{LangEN: list},
		fallback: list,
	})
	name := PickRandom(LangEN)
	if name == "" {
		t.Fatal("empty nickname")
	}
}

func TestParseAcceptLanguage(t *testing.T) {
	if ParseAcceptLanguage("es-ES,es;q=0.9") != LangES {
		t.Fatal("expected spanish")
	}
	if ParseAcceptLanguage("hi-IN") != LangHI {
		t.Fatal("expected hindi")
	}
	if ParseAcceptLanguage("pt-BR") != LangPT {
		t.Fatal("expected portuguese")
	}
	if ParseAcceptLanguage("id-ID") != LangID {
		t.Fatal("expected indonesian")
	}
	if ParseAcceptLanguage("") != DefaultLang {
		t.Fatal("expected default english")
	}
}
