import { gql } from '@apollo/client'

// ============================================
// QUERIES
// ============================================

export const GET_WORDS = gql`
  query GetWords($filter: WordFilter, $first: Int, $after: String) {
    words(filter: $filter, first: $first, after: $after) {
      edges {
        node {
          id
          text
          transcription
          audioUrl
          frequencyRank
          meanings {
            id
            partOfSpeech
            definitionEn
            translationRu
            cefrLevel
            imageUrl
            status
            nextReviewAt
            reviewCount
            examples {
              id
              sentenceEn
              sentenceRu
              sourceName
            }
            tags {
              id
              name
            }
          }
        }
        cursor
      }
      pageInfo {
        hasNextPage
        endCursor
      }
      totalCount
    }
  }
`

export const GET_WORD = gql`
  query GetWord($id: ID!) {
    word(id: $id) {
      id
      text
      transcription
      audioUrl
      frequencyRank
      meanings {
        id
        wordId
        partOfSpeech
        definitionEn
        translationRu
        cefrLevel
        imageUrl
        status
        nextReviewAt
        reviewCount
        examples {
          id
          sentenceEn
          sentenceRu
          sourceName
        }
        tags {
          id
          name
        }
      }
    }
  }
`

export const GET_INBOX_ITEMS = gql`
  query GetInboxItems {
    inboxItems {
      id
      text
      sourceContext
      createdAt
    }
  }
`

export const GET_SUGGEST = gql`
  query Suggest($query: String!) {
    suggest(query: $query) {
      text
      transcription
      translations
      origin
      existingWordId
    }
  }
`

export const GET_STUDY_QUEUE = gql`
  query GetStudyQueue($limit: Int) {
    studyQueue(limit: $limit) {
      id
      wordId
      partOfSpeech
      definitionEn
      translationRu
      cefrLevel
      imageUrl
      status
      nextReviewAt
      reviewCount
      examples {
        id
        sentenceEn
        sentenceRu
        sourceName
      }
      tags {
        id
        name
      }
    }
  }
`

export const GET_STATS = gql`
  query GetStats {
    stats {
      totalWords
      masteredCount
      learningCount
      dueForReviewCount
    }
  }
`

// ============================================
// MUTATIONS
// ============================================

export const CREATE_WORD = gql`
  mutation CreateWord($input: CreateWordInput!) {
    createWord(input: $input) {
      word {
        id
        text
        transcription
        audioUrl
        frequencyRank
        meanings {
          id
          partOfSpeech
          definitionEn
          translationRu
          cefrLevel
          imageUrl
          status
          nextReviewAt
          reviewCount
          examples {
            id
            sentenceEn
            sentenceRu
            sourceName
          }
          tags {
            id
            name
          }
        }
      }
    }
  }
`

export const UPDATE_WORD = gql`
  mutation UpdateWord($id: ID!, $input: CreateWordInput!) {
    updateWord(id: $id, input: $input) {
      word {
        id
        text
        transcription
        audioUrl
        frequencyRank
        meanings {
          id
          partOfSpeech
          definitionEn
          translationRu
          cefrLevel
          imageUrl
          status
          nextReviewAt
          reviewCount
          examples {
            id
            sentenceEn
            sentenceRu
            sourceName
          }
          tags {
            id
            name
          }
        }
      }
    }
  }
`

export const DELETE_WORD = gql`
  mutation DeleteWord($id: ID!) {
    deleteWord(id: $id)
  }
`

export const ADD_TO_INBOX = gql`
  mutation AddToInbox($text: String!, $sourceContext: String) {
    addToInbox(text: $text, sourceContext: $sourceContext) {
      id
      text
      sourceContext
      createdAt
    }
  }
`

export const DELETE_INBOX_ITEM = gql`
  mutation DeleteInboxItem($id: ID!) {
    deleteInboxItem(id: $id)
  }
`

export const CONVERT_INBOX_ITEM = gql`
  mutation ConvertInboxItem($inboxId: ID!, $input: CreateWordInput!) {
    convertInboxItem(inboxId: $inboxId, input: $input) {
      word {
        id
        text
        transcription
        audioUrl
        frequencyRank
        meanings {
          id
          partOfSpeech
          definitionEn
          translationRu
          cefrLevel
          imageUrl
          status
          nextReviewAt
          reviewCount
          examples {
            id
            sentenceEn
            sentenceRu
            sourceName
          }
          tags {
            id
            name
          }
        }
      }
    }
  }
`

export const REVIEW_MEANING = gql`
  mutation ReviewMeaning($meaningId: ID!, $grade: Int!) {
    reviewMeaning(meaningId: $meaningId, grade: $grade) {
      id
      wordId
      partOfSpeech
      definitionEn
      translationRu
      cefrLevel
      imageUrl
      status
      nextReviewAt
      reviewCount
      examples {
        id
        sentenceEn
        sentenceRu
        sourceName
      }
      tags {
        id
        name
      }
    }
  }
`


