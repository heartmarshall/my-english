package model

type DictionaryWordData struct {
	Word     DictionaryWord
	Meanings []DictionaryMeaningData
}

type DictionaryMeaningData struct {
	Meaning      DictionaryMeaning
	Translations []DictionaryTranslation
}
