package model

import (
	"strings"
	"testing"
)

func TestWordValidate(t *testing.T) {
	tests := []struct {
		name    string
		word    Word
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid word",
			word:    Word{Text: "hello"},
			wantErr: false,
		},
		{
			name:    "empty text",
			word:    Word{Text: ""},
			wantErr: true,
			errMsg:  "text",
		},
		{
			name:    "text too long",
			word:    Word{Text: strings.Repeat("a", 101)},
			wantErr: true,
			errMsg:  "text",
		},
		{
			name:    "valid with transcription",
			word:    Word{Text: "hello", Transcription: ptr("həˈloʊ")},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.word.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error should contain %q, got %q", tt.errMsg, err.Error())
				}
			}
		})
	}
}

func TestMeaningValidate(t *testing.T) {
	tests := []struct {
		name    string
		meaning Meaning
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid meaning",
			meaning: Meaning{
				PartOfSpeech:  PartOfSpeechNoun,
				TranslationRu: "привет",
			},
			wantErr: false,
		},
		{
			name: "empty translation",
			meaning: Meaning{
				PartOfSpeech:  PartOfSpeechNoun,
				TranslationRu: "",
			},
			wantErr: true,
			errMsg:  "translation_ru",
		},
		{
			name: "invalid part of speech",
			meaning: Meaning{
				PartOfSpeech:  PartOfSpeech("invalid"),
				TranslationRu: "тест",
			},
			wantErr: true,
			errMsg:  "part_of_speech",
		},
		{
			name: "translation too long",
			meaning: Meaning{
				PartOfSpeech:  PartOfSpeechVerb,
				TranslationRu: strings.Repeat("а", 501),
			},
			wantErr: true,
			errMsg:  "translation_ru",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.meaning.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error should contain %q, got %q", tt.errMsg, err.Error())
				}
			}
		})
	}
}

func TestExampleValidate(t *testing.T) {
	tests := []struct {
		name    string
		example Example
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid example",
			example: Example{SentenceEn: "Hello, world!"},
			wantErr: false,
		},
		{
			name:    "empty sentence",
			example: Example{SentenceEn: ""},
			wantErr: true,
			errMsg:  "sentence_en",
		},
		{
			name:    "with valid source",
			example: Example{SentenceEn: "Hello!", SourceName: ptrSource(ExampleSourceFilm)},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.example.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error should contain %q, got %q", tt.errMsg, err.Error())
				}
			}
		})
	}
}

func TestTagValidate(t *testing.T) {
	tests := []struct {
		name    string
		tag     Tag
		wantErr bool
	}{
		{
			name:    "valid tag",
			tag:     Tag{Name: "greetings"},
			wantErr: false,
		},
		{
			name:    "empty name",
			tag:     Tag{Name: ""},
			wantErr: true,
		},
		{
			name:    "name too long",
			tag:     Tag{Name: strings.Repeat("a", 51)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tag.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidationErrors(t *testing.T) {
	// Тест множественных ошибок
	meaning := Meaning{
		PartOfSpeech:  PartOfSpeech("invalid"),
		TranslationRu: "",
	}

	err := meaning.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	validationErrs, ok := err.(ValidationErrors)
	if !ok {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}

	if len(validationErrs) < 2 {
		t.Errorf("expected at least 2 errors, got %d", len(validationErrs))
	}
}

// Helper functions
func ptr(s string) *string {
	return &s
}

func ptrSource(s ExampleSource) *ExampleSource {
	return &s
}
