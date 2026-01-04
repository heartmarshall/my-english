package model

// PartOfSpeech представляет часть речи.
type PartOfSpeech string

const (
	PartOfSpeechNoun         PartOfSpeech = "noun"
	PartOfSpeechVerb         PartOfSpeech = "verb"
	PartOfSpeechAdjective    PartOfSpeech = "adjective"
	PartOfSpeechAdverb       PartOfSpeech = "adverb"
	PartOfSpeechPronoun      PartOfSpeech = "pronoun"
	PartOfSpeechPreposition  PartOfSpeech = "preposition"
	PartOfSpeechConjunction  PartOfSpeech = "conjunction"
	PartOfSpeechInterjection PartOfSpeech = "interjection"
	PartOfSpeechPhrase       PartOfSpeech = "phrase"
	PartOfSpeechIdiom        PartOfSpeech = "idiom"
	PartOfSpeechOther        PartOfSpeech = "other"
)

func (p PartOfSpeech) String() string { return string(p) }

func (p PartOfSpeech) IsValid() bool {
	switch p {
	case PartOfSpeechNoun, PartOfSpeechVerb, PartOfSpeechAdjective,
		PartOfSpeechAdverb, PartOfSpeechPronoun, PartOfSpeechPreposition,
		PartOfSpeechConjunction, PartOfSpeechInterjection,
		PartOfSpeechPhrase, PartOfSpeechIdiom, PartOfSpeechOther:
		return true
	}
	return false
}

// AccentRegion представляет регион акцента.
type AccentRegion string

const (
	AccentRegionUS      AccentRegion = "us"
	AccentRegionUK      AccentRegion = "uk"
	AccentRegionAU      AccentRegion = "au"
	AccentRegionGeneral AccentRegion = "general"
)

func (a AccentRegion) String() string { return string(a) }

// RelationType представляет тип семантической связи.
type RelationType string

const (
	RelationTypeSynonym     RelationType = "synonym"
	RelationTypeAntonym     RelationType = "antonym"
	RelationTypeRelated     RelationType = "related"
	RelationTypeCollocation RelationType = "collocation"
)

func (r RelationType) String() string { return string(r) }

// MorphologicalType представляет тип морфологической формы.
type MorphologicalType string

const (
	MorphTypePlural            MorphologicalType = "plural"
	MorphTypePastTense         MorphologicalType = "past_tense"
	MorphTypePastParticiple    MorphologicalType = "past_participle"
	MorphTypePresentParticiple MorphologicalType = "present_participle"
	MorphTypeComparative       MorphologicalType = "comparative"
	MorphTypeSuperlative       MorphologicalType = "superlative"
)

func (m MorphologicalType) String() string { return string(m) }
