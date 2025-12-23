package model

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

// Validation constants
const (
	MaxWordLength        = 100
	MaxTranslationLength = 500
	MaxDefinitionLength  = 1000
	MaxSentenceLength    = 1000
	MaxTagNameLength     = 50
	MaxURLLength         = 2048
)

// Validator — глобальный валидатор (thread-safe).
var (
	validate     *validator.Validate
	validateOnce sync.Once
)

// GetValidator возвращает настроенный валидатор.
func GetValidator() *validator.Validate {
	validateOnce.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())

		// Регистрируем кастомные валидаторы
		validate.RegisterValidation("part_of_speech", validatePartOfSpeech)
		validate.RegisterValidation("learning_status", validateLearningStatus)
		validate.RegisterValidation("example_source", validateExampleSource)
	})
	return validate
}

// --- Custom validators ---

func validatePartOfSpeech(fl validator.FieldLevel) bool {
	pos, ok := fl.Field().Interface().(PartOfSpeech)
	if !ok {
		return false
	}
	return pos.IsValid()
}

func validateLearningStatus(fl validator.FieldLevel) bool {
	status, ok := fl.Field().Interface().(LearningStatus)
	if !ok {
		return false
	}
	return status.IsValid()
}

func validateExampleSource(fl validator.FieldLevel) bool {
	field := fl.Field()

	// Проверяем тип поля
	switch field.Kind() {
	case reflect.Ptr:
		if field.IsNil() {
			return true // nil допустим
		}
		src, ok := field.Interface().(*ExampleSource)
		if !ok {
			return false
		}
		return src.IsValid()
	case reflect.String:
		src := ExampleSource(field.String())
		return src.IsValid()
	default:
		return false
	}
}

// --- ValidationError ---

// ValidationError представляет ошибку валидации.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   any    `json:"value,omitempty"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors — список ошибок валидации.
type ValidationErrors []*ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, e := range ve {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(e.Error())
	}
	return sb.String()
}

// translateError преобразует ошибку validator в понятное сообщение.
func translateError(fe validator.FieldError) *ValidationError {
	field := toSnakeCase(fe.Field())
	tag := fe.Tag()

	var message string
	switch tag {
	case "required":
		message = "is required"
	case "max":
		message = fmt.Sprintf("must be at most %s characters", fe.Param())
	case "min":
		message = fmt.Sprintf("must be at least %s characters", fe.Param())
	case "url":
		message = "must be a valid URL"
	case "part_of_speech":
		message = "must be one of: noun, verb, adjective, adverb, other"
	case "learning_status":
		message = "must be one of: new, learning, review, mastered"
	case "example_source":
		message = "must be one of: film, book, chat, video, podcast"
	default:
		message = fmt.Sprintf("failed on '%s' validation", tag)
	}

	return &ValidationError{
		Field:   field,
		Message: message,
		Tag:     tag,
		Value:   fe.Value(),
	}
}

// toSnakeCase конвертирует CamelCase в snake_case.
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// --- Validate functions ---

// Validate валидирует любую структуру.
func Validate(s any) error {
	err := GetValidator().Struct(s)
	if err == nil {
		return nil
	}

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		result := make(ValidationErrors, 0, len(validationErrs))
		for _, fe := range validationErrs {
			result = append(result, translateError(fe))
		}
		return result
	}

	return err
}

// --- Validatable structs with tags ---

// WordInput — структура для валидации входных данных слова.
type WordInput struct {
	Text          string  `validate:"required,min=1,max=100"`
	Transcription *string `validate:"omitempty,max=100"`
	AudioURL      *string `validate:"omitempty,url,max=2048"`
}

// MeaningInput — структура для валидации входных данных значения.
type MeaningInput struct {
	PartOfSpeech  PartOfSpeech `validate:"required,part_of_speech"`
	TranslationRu string       `validate:"required,min=1,max=500"`
	DefinitionEn  *string      `validate:"omitempty,max=1000"`
	CefrLevel     *string      `validate:"omitempty,max=10"`
	ImageURL      *string      `validate:"omitempty,url,max=2048"`
}

// ExampleInput — структура для валидации входных данных примера.
type ExampleInput struct {
	SentenceEn string         `validate:"required,min=1,max=1000"`
	SentenceRu *string        `validate:"omitempty,max=1000"`
	SourceName *ExampleSource `validate:"omitempty,example_source"`
}

// TagInput — структура для валидации входных данных тега.
type TagInput struct {
	Name string `validate:"required,min=1,max=50"`
}

// --- Model Validate methods ---

// Validate проверяет валидность Word.
func (w *Word) Validate() error {
	input := WordInput{
		Text:          w.Text,
		Transcription: w.Transcription,
		AudioURL:      w.AudioURL,
	}
	return Validate(input)
}

// Validate проверяет валидность Meaning.
func (m *Meaning) Validate() error {
	input := MeaningInput{
		PartOfSpeech:  m.PartOfSpeech,
		TranslationRu: m.TranslationRu,
		DefinitionEn:  m.DefinitionEn,
		CefrLevel:     m.CefrLevel,
		ImageURL:      m.ImageURL,
	}
	return Validate(input)
}

// Validate проверяет валидность Example.
func (e *Example) Validate() error {
	input := ExampleInput{
		SentenceEn: e.SentenceEn,
		SentenceRu: e.SentenceRu,
		SourceName: e.SourceName,
	}
	return Validate(input)
}

// Validate проверяет валидность Tag.
func (t *Tag) Validate() error {
	input := TagInput{
		Name: t.Name,
	}
	return Validate(input)
}
