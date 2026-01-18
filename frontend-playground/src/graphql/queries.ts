import { gql } from '@apollo/client';

// ==============================================================================
// QUERIES
// ==============================================================================

export const GET_DASHBOARD_STATS = gql`
  query GetDashboardStats {
    dashboardStats {
      totalWords
      totalCards
      newCards
      learningCards
      reviewCards
      masteredCards
      dueToday
    }
  }
`;

export const GET_DICTIONARY = gql`
  query GetDictionary($filter: WordFilter) {
    dictionary(filter: $filter) {
      id
      text
      textNormalized
      cardEnabled
      card {
        id
        status
        nextReviewAt
        intervalDays
        easeFactor
      }
      senses {
        id
        definition
        partOfSpeech
        translations {
          id
          text
        }
        examples {
          id
          sentence
          translation
        }
      }
      images {
        id
        url
        caption
      }
      pronunciations {
        id
        audioUrl
        transcription
        region
      }
      createdAt
      updatedAt
    }
  }
`;

export const GET_DICTIONARY_ENTRY = gql`
  query GetDictionaryEntry($id: UUID!) {
    dictionaryEntry(id: $id) {
      id
      text
      textNormalized
      cardEnabled
      card {
        id
        status
        nextReviewAt
        intervalDays
        easeFactor
        reviewHistory(limit: 10) {
          id
          grade
          durationMs
          reviewedAt
        }
      }
      senses {
        id
        definition
        partOfSpeech
        sourceSlug
        translations {
          id
          text
          sourceSlug
        }
        examples {
          id
          sentence
          translation
          sourceSlug
          createdAt
        }
      }
      images {
        id
        url
        caption
        sourceSlug
      }
      pronunciations {
        id
        audioUrl
        transcription
        region
        sourceSlug
      }
      auditLog {
        id
        entityType
        action
        changes
        createdAt
      }
      createdAt
      updatedAt
    }
  }
`;

export const GET_STUDY_QUEUE = gql`
  query GetStudyQueue($limit: Int) {
    studyQueue(limit: $limit) {
      id
      text
      textNormalized
      card {
        id
        status
        nextReviewAt
        intervalDays
        easeFactor
      }
      senses {
        id
        definition
        partOfSpeech
        translations {
          id
          text
        }
        examples {
          id
          sentence
          translation
        }
      }
      images {
        id
        url
        caption
      }
      pronunciations {
        id
        audioUrl
        transcription
        region
      }
    }
  }
`;

export const GET_INBOX_ITEMS = gql`
  query GetInboxItems {
    inboxItems {
      id
      text
      context
      createdAt
    }
  }
`;

export const FETCH_SUGGESTIONS = gql`
  query FetchSuggestions($text: String!, $sources: [String!]!) {
    fetchSuggestions(text: $text, sources: $sources) {
      sourceSlug
      sourceName
      senses {
        definition
        partOfSpeech
        translations
        examples {
          sentence
          translation
        }
      }
      images {
        url
        thumbnailUrl
        caption
      }
      pronunciations {
        audioUrl
        transcription
        region
      }
    }
  }
`;

// ==============================================================================
// MUTATIONS
// ==============================================================================

export const CREATE_WORD = gql`
  mutation CreateWord($input: CreateWordInput!) {
    createWord(input: $input) {
      id
      text
      textNormalized
      cardEnabled
      card {
        id
        status
        nextReviewAt
      }
      senses {
        id
        definition
        partOfSpeech
        translations {
          id
          text
        }
      }
      createdAt
    }
  }
`;

export const UPDATE_WORD = gql`
  mutation UpdateWord($id: UUID!, $input: UpdateWordInput!) {
    updateWord(id: $id, input: $input) {
      id
      text
      textNormalized
      senses {
        id
        definition
        partOfSpeech
        translations {
          id
          text
        }
        examples {
          id
          sentence
          translation
        }
      }
      updatedAt
    }
  }
`;

export const DELETE_WORD = gql`
  mutation DeleteWord($id: UUID!) {
    deleteWord(id: $id)
  }
`;

export const ADD_TO_INBOX = gql`
  mutation AddToInbox($text: String!, $context: String) {
    addToInbox(text: $text, context: $context) {
      id
      text
      context
      createdAt
    }
  }
`;

export const DELETE_INBOX_ITEM = gql`
  mutation DeleteInboxItem($id: UUID!) {
    deleteInboxItem(id: $id)
  }
`;

export const CONVERT_INBOX_TO_WORD = gql`
  mutation ConvertInboxToWord($inboxId: UUID!, $input: CreateWordInput!) {
    convertInboxToWord(inboxId: $inboxId, input: $input) {
      id
      text
      textNormalized
      cardEnabled
      senses {
        id
        definition
        translations {
          id
          text
        }
      }
      createdAt
    }
  }
`;

export const REVIEW_CARD = gql`
  mutation ReviewCard($cardId: UUID!, $grade: ReviewGrade!, $timeTakenMs: Int) {
    reviewCard(cardId: $cardId, grade: $grade, timeTakenMs: $timeTakenMs) {
      entry {
        id
        text
        card {
          id
          status
          nextReviewAt
          intervalDays
          easeFactor
        }
      }
      nextReviewAt
    }
  }
`;

