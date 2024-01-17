package nanoc

type SupportedLanguage string

const (
	LanguageCxx SupportedLanguage = "c++"

	LanguageSwift SupportedLanguage = "swift"

	LanguageTypeScript SupportedLanguage = "ts"
)

type Options struct {
	Language                  SupportedLanguage
	MessageFactoryAbsFilePath string
	InputFileAbsolutePaths    []string
	CodeFormatterPath         string
	CodeFormatterArgs         []string
}

func IsLanguageSupported(language string) bool {
	switch language {
	case string(LanguageCxx), string(LanguageSwift), string(LanguageTypeScript):
		return true
	default:
		return false
	}
}

func DefaultFormatter(lang SupportedLanguage) (string, []string) {
	switch lang {
	case LanguageCxx:
		return "clang-format", []string{"-i", "-style=LLVM"}

	case LanguageSwift:
		return "swift-format", []string{"--in-place"}

	default:
		return "", nil
	}
}
