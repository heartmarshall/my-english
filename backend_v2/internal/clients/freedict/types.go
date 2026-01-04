package freedict

// apiEntry описывает структуру одного элемента массива в ответе API.
type apiEntry struct {
	Word       string        `json:"word"`
	Phonetic   string        `json:"phonetic"`
	Phonetics  []apiPhonetic `json:"phonetics"`
	Meanings   []apiMeaning  `json:"meanings"`
	SourceUrls []string      `json:"sourceUrls"`
}

type apiPhonetic struct {
	Text  string `json:"text"`
	Audio string `json:"audio"`
}

type apiMeaning struct {
	PartOfSpeech string          `json:"partOfSpeech"`
	Definitions  []apiDefinition `json:"definitions"`
}

type apiDefinition struct {
	Definition string   `json:"definition"`
	Example    string   `json:"example"`
	Synonyms   []string `json:"synonyms"`
	Antonyms   []string `json:"antonyms"`
}
