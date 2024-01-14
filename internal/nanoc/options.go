package nanoc

type SupportedLanguage string

const (
	LanguageCxx SupportedLanguage = "c++"

	LanguageSwift SupportedLanguage = "swift"

	LanguageTypeScript SupportedLanguage = "ts"
)

func IsLanguageSupported(language string) bool {
	switch language {
	case string(LanguageCxx), string(LanguageSwift), string(LanguageTypeScript):
		return true
	default:
		return false
	}
}

type Options struct {
	Language                  SupportedLanguage
	MessageFactoryAbsFilePath string
	InputFileAbsolutePaths    []string
}
