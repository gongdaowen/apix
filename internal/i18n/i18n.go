package i18n

import (
	"fmt"
	"os"
	"strings"
)

// Language represents a supported language
type Language string

const (
	EN Language = "en"
	ZH Language = "zh"
)

// Translator handles internationalization
type Translator struct {
	currentLang Language
	messages    map[Language]map[string]string
}

// NewTranslator creates a new translator with auto-detected language
func NewTranslator() *Translator {
	t := &Translator{
		currentLang: detectLanguage(),
		messages:    make(map[Language]map[string]string),
	}
	t.loadMessages()
	return t
}

// NewTranslatorWithLang creates a translator with specified language
func NewTranslatorWithLang(lang string) *Translator {
	t := &Translator{
		currentLang: parseLanguage(lang),
		messages:    make(map[Language]map[string]string),
	}
	t.loadMessages()
	return t
}

// SetLanguage changes the current language
func (t *Translator) SetLanguage(lang string) {
	t.currentLang = parseLanguage(lang)
}

// GetLanguage returns the current language
func (t *Translator) GetLanguage() Language {
	return t.currentLang
}

// T translates a key to the current language
func (t *Translator) T(key string) string {
	if msgs, ok := t.messages[t.currentLang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	// Fallback to English
	if msgs, ok := t.messages[EN]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	// Return key if translation not found
	return key
}

// TF translates a key with format arguments
func (t *Translator) TF(key string, args ...interface{}) string {
	msg := t.T(key)
	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}

// loadMessages loads all translation messages
func (t *Translator) loadMessages() {
	t.messages[EN] = loadEnglishMessages()
	t.messages[ZH] = loadChineseMessages()
}

// detectLanguage automatically detects system language
func detectLanguage() Language {
	// Check APIX_LANG environment variable first (project-specific)
	lang := os.Getenv("APIX_LANG")
	if lang != "" {
		return parseLanguage(lang)
	}

	// Check LANG environment variable (Linux/Mac/Git Bash)
	lang = os.Getenv("LANG")
	if lang != "" {
		return parseLanguage(lang)
	}

	// Check LC_ALL environment variable
	lang = os.Getenv("LC_ALL")
	if lang != "" {
		return parseLanguage(lang)
	}

	// Check LANGUAGE environment variable
	lang = os.Getenv("LANGUAGE")
	if lang != "" {
		return parseLanguage(lang)
	}

	// Default to English
	return EN
}

// parseLanguage parses a language string to Language type
func parseLanguage(lang string) Language {
	lang = strings.ToLower(lang)
	
	// Extract language code (e.g., "en_US.UTF-8" -> "en")
	if idx := strings.Index(lang, "_"); idx != -1 {
		lang = lang[:idx]
	}
	if idx := strings.Index(lang, "."); idx != -1 {
		lang = lang[:idx]
	}
	
	// Match supported languages
	switch {
	case strings.HasPrefix(lang, "zh"):
		return ZH
	case strings.HasPrefix(lang, "en"):
		return EN
	default:
		return EN // Default to English
	}
}

// formatMessage formats a message with arguments (simple implementation)
func formatMessage(msg string, args ...interface{}) string {
	// For now, just return the message
	// Can be enhanced to support sprintf-style formatting
	_ = args
	return msg
}
