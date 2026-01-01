/* eslint-disable */
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = T | null | undefined;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Time: { input: any; output: any; }
};

export enum CefrLevel {
  A1 = 'A1',
  A2 = 'A2',
  B1 = 'B1',
  B2 = 'B2',
  C1 = 'C1',
  C2 = 'C2'
}

export type CreateWordInput = {
  audioUrl?: InputMaybe<Scalars['String']['input']>;
  meanings?: InputMaybe<Array<MeaningInput>>;
  sourceContext?: InputMaybe<Scalars['String']['input']>;
  text: Scalars['String']['input'];
  transcription?: InputMaybe<Scalars['String']['input']>;
};

export type CreateWordPayload = {
  __typename?: 'CreateWordPayload';
  word: Word;
};

export type DashboardStats = {
  __typename?: 'DashboardStats';
  dueForReviewCount: Scalars['Int']['output'];
  learningCount: Scalars['Int']['output'];
  masteredCount: Scalars['Int']['output'];
  totalWords: Scalars['Int']['output'];
};

export type Example = {
  __typename?: 'Example';
  id: Scalars['ID']['output'];
  sentenceEn: Scalars['String']['output'];
  sentenceRu?: Maybe<Scalars['String']['output']>;
  sourceName?: Maybe<ExampleSource>;
};

export type ExampleInput = {
  sentenceEn: Scalars['String']['input'];
  sentenceRu?: InputMaybe<Scalars['String']['input']>;
  sourceName?: InputMaybe<Scalars['String']['input']>;
};

export enum ExampleSource {
  Book = 'BOOK',
  Chat = 'CHAT',
  Film = 'FILM',
  Podcast = 'PODCAST',
  Video = 'VIDEO'
}

export type InboxItem = {
  __typename?: 'InboxItem';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  sourceContext?: Maybe<Scalars['String']['output']>;
  text: Scalars['String']['output'];
};

export enum LearningStatus {
  Learning = 'LEARNING',
  Mastered = 'MASTERED',
  New = 'NEW',
  Review = 'REVIEW'
}

export type Meaning = {
  __typename?: 'Meaning';
  cefrLevel?: Maybe<Scalars['String']['output']>;
  definitionEn?: Maybe<Scalars['String']['output']>;
  examples?: Maybe<Array<Example>>;
  id: Scalars['ID']['output'];
  imageUrl?: Maybe<Scalars['String']['output']>;
  nextReviewAt?: Maybe<Scalars['Time']['output']>;
  partOfSpeech: PartOfSpeech;
  reviewCount: Scalars['Int']['output'];
  status: LearningStatus;
  synonymsAntonyms?: Maybe<Array<SynonymAntonym>>;
  tags?: Maybe<Array<Tag>>;
  translationRu?: Maybe<Array<Scalars['String']['output']>>;
  wordId: Scalars['ID']['output'];
};

export type MeaningInput = {
  definitionEn?: InputMaybe<Scalars['String']['input']>;
  examples?: InputMaybe<Array<ExampleInput>>;
  imageUrl?: InputMaybe<Scalars['String']['input']>;
  partOfSpeech: PartOfSpeech;
  tags?: InputMaybe<Array<Scalars['String']['input']>>;
  translationRu: Scalars['String']['input'];
};

export type Mutation = {
  __typename?: 'Mutation';
  addToInbox: InboxItem;
  convertInboxItem: CreateWordPayload;
  createWord: CreateWordPayload;
  deleteInboxItem: Scalars['Boolean']['output'];
  deleteWord: Scalars['Boolean']['output'];
  reviewMeaning: Meaning;
  updateWord: UpdateWordPayload;
};


export type MutationAddToInboxArgs = {
  sourceContext?: InputMaybe<Scalars['String']['input']>;
  text: Scalars['String']['input'];
};


export type MutationConvertInboxItemArgs = {
  inboxId: Scalars['ID']['input'];
  input: CreateWordInput;
};


export type MutationCreateWordArgs = {
  input: CreateWordInput;
};


export type MutationDeleteInboxItemArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteWordArgs = {
  id: Scalars['ID']['input'];
};


export type MutationReviewMeaningArgs = {
  grade: Scalars['Int']['input'];
  meaningId: Scalars['ID']['input'];
};


export type MutationUpdateWordArgs = {
  id: Scalars['ID']['input'];
  input: CreateWordInput;
};

export type PageInfo = {
  __typename?: 'PageInfo';
  endCursor?: Maybe<Scalars['String']['output']>;
  hasNextPage: Scalars['Boolean']['output'];
};

export enum PartOfSpeech {
  Adjective = 'ADJECTIVE',
  Adverb = 'ADVERB',
  Noun = 'NOUN',
  Other = 'OTHER',
  Verb = 'VERB'
}

export type Query = {
  __typename?: 'Query';
  inboxItems: Array<InboxItem>;
  stats: DashboardStats;
  studyQueue: Array<Meaning>;
  suggest: Array<Suggestion>;
  word?: Maybe<Word>;
  wordByForm?: Maybe<Word>;
  words: WordConnection;
};


export type QueryStudyQueueArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
};


export type QuerySuggestArgs = {
  query: Scalars['String']['input'];
};


export type QueryWordArgs = {
  id: Scalars['ID']['input'];
};


export type QueryWordByFormArgs = {
  formText: Scalars['String']['input'];
};


export type QueryWordsArgs = {
  after?: InputMaybe<Scalars['String']['input']>;
  filter?: InputMaybe<WordFilter>;
  first?: InputMaybe<Scalars['Int']['input']>;
};

export enum RelationType {
  Antonym = 'ANTONYM',
  Synonym = 'SYNONYM'
}

export type Suggestion = {
  __typename?: 'Suggestion';
  existingWordId?: Maybe<Scalars['ID']['output']>;
  origin: SuggestionOrigin;
  text: Scalars['String']['output'];
  transcription?: Maybe<Scalars['String']['output']>;
  translations: Array<Scalars['String']['output']>;
};

export enum SuggestionOrigin {
  Dictionary = 'DICTIONARY',
  Local = 'LOCAL'
}

export type SynonymAntonym = {
  __typename?: 'SynonymAntonym';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  relatedMeaningId: Scalars['ID']['output'];
  relationType: RelationType;
  updatedAt: Scalars['Time']['output'];
};

export type Tag = {
  __typename?: 'Tag';
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
};

export type UpdateWordPayload = {
  __typename?: 'UpdateWordPayload';
  word: Word;
};

export type Word = {
  __typename?: 'Word';
  audioUrl?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['Time']['output'];
  forms?: Maybe<Array<WordForm>>;
  frequencyRank?: Maybe<Scalars['Int']['output']>;
  id: Scalars['ID']['output'];
  meanings?: Maybe<Array<Meaning>>;
  text: Scalars['String']['output'];
  transcription?: Maybe<Scalars['String']['output']>;
};

export type WordConnection = {
  __typename?: 'WordConnection';
  edges: Array<WordEdge>;
  pageInfo: PageInfo;
  totalCount: Scalars['Int']['output'];
};

export type WordEdge = {
  __typename?: 'WordEdge';
  cursor: Scalars['String']['output'];
  node: Word;
};

export type WordFilter = {
  search?: InputMaybe<Scalars['String']['input']>;
  status?: InputMaybe<LearningStatus>;
  tags?: InputMaybe<Array<Scalars['String']['input']>>;
};

export type WordForm = {
  __typename?: 'WordForm';
  createdAt: Scalars['Time']['output'];
  formText: Scalars['String']['output'];
  formType?: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type GetWordsQueryVariables = Exact<{
  filter?: InputMaybe<WordFilter>;
  first?: InputMaybe<Scalars['Int']['input']>;
  after?: InputMaybe<Scalars['String']['input']>;
}>;


export type GetWordsQuery = { __typename?: 'Query', words: { __typename?: 'WordConnection', totalCount: number, edges: Array<{ __typename?: 'WordEdge', cursor: string, node: { __typename?: 'Word', id: string, text: string, transcription?: string | null, audioUrl?: string | null, frequencyRank?: number | null, createdAt: any, forms?: Array<{ __typename?: 'WordForm', id: string, formText: string, formType?: string | null }> | null, meanings?: Array<{ __typename?: 'Meaning', id: string, partOfSpeech: PartOfSpeech, definitionEn?: string | null, translationRu?: Array<string> | null, cefrLevel?: string | null, imageUrl?: string | null, status: LearningStatus, nextReviewAt?: any | null, reviewCount: number, examples?: Array<{ __typename?: 'Example', id: string, sentenceEn: string, sentenceRu?: string | null, sourceName?: ExampleSource | null }> | null, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null, synonymsAntonyms?: Array<{ __typename?: 'SynonymAntonym', id: string, relatedMeaningId: string, relationType: RelationType }> | null }> | null } }>, pageInfo: { __typename?: 'PageInfo', hasNextPage: boolean, endCursor?: string | null } } };

export type GetWordQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetWordQuery = { __typename?: 'Query', word?: { __typename?: 'Word', id: string, text: string, transcription?: string | null, audioUrl?: string | null, frequencyRank?: number | null, createdAt: any, forms?: Array<{ __typename?: 'WordForm', id: string, formText: string, formType?: string | null }> | null, meanings?: Array<{ __typename?: 'Meaning', id: string, wordId: string, partOfSpeech: PartOfSpeech, definitionEn?: string | null, translationRu?: Array<string> | null, cefrLevel?: string | null, imageUrl?: string | null, status: LearningStatus, nextReviewAt?: any | null, reviewCount: number, examples?: Array<{ __typename?: 'Example', id: string, sentenceEn: string, sentenceRu?: string | null, sourceName?: ExampleSource | null }> | null, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null, synonymsAntonyms?: Array<{ __typename?: 'SynonymAntonym', id: string, relatedMeaningId: string, relationType: RelationType }> | null }> | null } | null };

export type GetInboxItemsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetInboxItemsQuery = { __typename?: 'Query', inboxItems: Array<{ __typename?: 'InboxItem', id: string, text: string, sourceContext?: string | null, createdAt: any }> };

export type SuggestQueryVariables = Exact<{
  query: Scalars['String']['input'];
}>;


export type SuggestQuery = { __typename?: 'Query', suggest: Array<{ __typename?: 'Suggestion', text: string, transcription?: string | null, translations: Array<string>, origin: SuggestionOrigin, existingWordId?: string | null }> };

export type GetStudyQueueQueryVariables = Exact<{
  limit?: InputMaybe<Scalars['Int']['input']>;
}>;


export type GetStudyQueueQuery = { __typename?: 'Query', studyQueue: Array<{ __typename?: 'Meaning', id: string, wordId: string, partOfSpeech: PartOfSpeech, definitionEn?: string | null, translationRu?: Array<string> | null, cefrLevel?: string | null, imageUrl?: string | null, status: LearningStatus, nextReviewAt?: any | null, reviewCount: number, examples?: Array<{ __typename?: 'Example', id: string, sentenceEn: string, sentenceRu?: string | null, sourceName?: ExampleSource | null }> | null, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null, synonymsAntonyms?: Array<{ __typename?: 'SynonymAntonym', id: string, relatedMeaningId: string, relationType: RelationType }> | null }> };

export type GetStatsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetStatsQuery = { __typename?: 'Query', stats: { __typename?: 'DashboardStats', totalWords: number, masteredCount: number, learningCount: number, dueForReviewCount: number } };

export type CreateWordMutationVariables = Exact<{
  input: CreateWordInput;
}>;


export type CreateWordMutation = { __typename?: 'Mutation', createWord: { __typename?: 'CreateWordPayload', word: { __typename?: 'Word', id: string, text: string, transcription?: string | null, audioUrl?: string | null, frequencyRank?: number | null, forms?: Array<{ __typename?: 'WordForm', id: string, formText: string, formType?: string | null }> | null, meanings?: Array<{ __typename?: 'Meaning', id: string, partOfSpeech: PartOfSpeech, definitionEn?: string | null, translationRu?: Array<string> | null, cefrLevel?: string | null, imageUrl?: string | null, status: LearningStatus, nextReviewAt?: any | null, reviewCount: number, examples?: Array<{ __typename?: 'Example', id: string, sentenceEn: string, sentenceRu?: string | null, sourceName?: ExampleSource | null }> | null, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null, synonymsAntonyms?: Array<{ __typename?: 'SynonymAntonym', id: string, relatedMeaningId: string, relationType: RelationType }> | null }> | null } } };

export type UpdateWordMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  input: CreateWordInput;
}>;


export type UpdateWordMutation = { __typename?: 'Mutation', updateWord: { __typename?: 'UpdateWordPayload', word: { __typename?: 'Word', id: string, text: string, transcription?: string | null, audioUrl?: string | null, frequencyRank?: number | null, forms?: Array<{ __typename?: 'WordForm', id: string, formText: string, formType?: string | null }> | null, meanings?: Array<{ __typename?: 'Meaning', id: string, partOfSpeech: PartOfSpeech, definitionEn?: string | null, translationRu?: Array<string> | null, cefrLevel?: string | null, imageUrl?: string | null, status: LearningStatus, nextReviewAt?: any | null, reviewCount: number, examples?: Array<{ __typename?: 'Example', id: string, sentenceEn: string, sentenceRu?: string | null, sourceName?: ExampleSource | null }> | null, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null, synonymsAntonyms?: Array<{ __typename?: 'SynonymAntonym', id: string, relatedMeaningId: string, relationType: RelationType }> | null }> | null } } };

export type DeleteWordMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteWordMutation = { __typename?: 'Mutation', deleteWord: boolean };

export type AddToInboxMutationVariables = Exact<{
  text: Scalars['String']['input'];
  sourceContext?: InputMaybe<Scalars['String']['input']>;
}>;


export type AddToInboxMutation = { __typename?: 'Mutation', addToInbox: { __typename?: 'InboxItem', id: string, text: string, sourceContext?: string | null, createdAt: any } };

export type DeleteInboxItemMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteInboxItemMutation = { __typename?: 'Mutation', deleteInboxItem: boolean };

export type ConvertInboxItemMutationVariables = Exact<{
  inboxId: Scalars['ID']['input'];
  input: CreateWordInput;
}>;


export type ConvertInboxItemMutation = { __typename?: 'Mutation', convertInboxItem: { __typename?: 'CreateWordPayload', word: { __typename?: 'Word', id: string, text: string, transcription?: string | null, audioUrl?: string | null, frequencyRank?: number | null, forms?: Array<{ __typename?: 'WordForm', id: string, formText: string, formType?: string | null }> | null, meanings?: Array<{ __typename?: 'Meaning', id: string, partOfSpeech: PartOfSpeech, definitionEn?: string | null, translationRu?: Array<string> | null, cefrLevel?: string | null, imageUrl?: string | null, status: LearningStatus, nextReviewAt?: any | null, reviewCount: number, examples?: Array<{ __typename?: 'Example', id: string, sentenceEn: string, sentenceRu?: string | null, sourceName?: ExampleSource | null }> | null, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null, synonymsAntonyms?: Array<{ __typename?: 'SynonymAntonym', id: string, relatedMeaningId: string, relationType: RelationType }> | null }> | null } } };

export type ReviewMeaningMutationVariables = Exact<{
  meaningId: Scalars['ID']['input'];
  grade: Scalars['Int']['input'];
}>;


export type ReviewMeaningMutation = { __typename?: 'Mutation', reviewMeaning: { __typename?: 'Meaning', id: string, wordId: string, partOfSpeech: PartOfSpeech, definitionEn?: string | null, translationRu?: Array<string> | null, cefrLevel?: string | null, imageUrl?: string | null, status: LearningStatus, nextReviewAt?: any | null, reviewCount: number, examples?: Array<{ __typename?: 'Example', id: string, sentenceEn: string, sentenceRu?: string | null, sourceName?: ExampleSource | null }> | null, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null, synonymsAntonyms?: Array<{ __typename?: 'SynonymAntonym', id: string, relatedMeaningId: string, relationType: RelationType }> | null } };


export const GetWordsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetWords"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"WordFilter"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"first"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"after"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"words"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}},{"kind":"Argument","name":{"kind":"Name","value":"first"},"value":{"kind":"Variable","name":{"kind":"Name","value":"first"}}},{"kind":"Argument","name":{"kind":"Name","value":"after"},"value":{"kind":"Variable","name":{"kind":"Name","value":"after"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"edges"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"node"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"transcription"}},{"kind":"Field","name":{"kind":"Name","value":"audioUrl"}},{"kind":"Field","name":{"kind":"Name","value":"frequencyRank"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"forms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"formText"}},{"kind":"Field","name":{"kind":"Name","value":"formType"}}]}},{"kind":"Field","name":{"kind":"Name","value":"meanings"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"partOfSpeech"}},{"kind":"Field","name":{"kind":"Name","value":"definitionEn"}},{"kind":"Field","name":{"kind":"Name","value":"translationRu"}},{"kind":"Field","name":{"kind":"Name","value":"cefrLevel"}},{"kind":"Field","name":{"kind":"Name","value":"imageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"nextReviewAt"}},{"kind":"Field","name":{"kind":"Name","value":"reviewCount"}},{"kind":"Field","name":{"kind":"Name","value":"examples"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceEn"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceRu"}},{"kind":"Field","name":{"kind":"Name","value":"sourceName"}}]}},{"kind":"Field","name":{"kind":"Name","value":"tags"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"synonymsAntonyms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"relatedMeaningId"}},{"kind":"Field","name":{"kind":"Name","value":"relationType"}}]}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"cursor"}}]}},{"kind":"Field","name":{"kind":"Name","value":"pageInfo"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"hasNextPage"}},{"kind":"Field","name":{"kind":"Name","value":"endCursor"}}]}},{"kind":"Field","name":{"kind":"Name","value":"totalCount"}}]}}]}}]} as unknown as DocumentNode<GetWordsQuery, GetWordsQueryVariables>;
export const GetWordDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetWord"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"word"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"transcription"}},{"kind":"Field","name":{"kind":"Name","value":"audioUrl"}},{"kind":"Field","name":{"kind":"Name","value":"frequencyRank"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"forms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"formText"}},{"kind":"Field","name":{"kind":"Name","value":"formType"}}]}},{"kind":"Field","name":{"kind":"Name","value":"meanings"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"wordId"}},{"kind":"Field","name":{"kind":"Name","value":"partOfSpeech"}},{"kind":"Field","name":{"kind":"Name","value":"definitionEn"}},{"kind":"Field","name":{"kind":"Name","value":"translationRu"}},{"kind":"Field","name":{"kind":"Name","value":"cefrLevel"}},{"kind":"Field","name":{"kind":"Name","value":"imageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"nextReviewAt"}},{"kind":"Field","name":{"kind":"Name","value":"reviewCount"}},{"kind":"Field","name":{"kind":"Name","value":"examples"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceEn"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceRu"}},{"kind":"Field","name":{"kind":"Name","value":"sourceName"}}]}},{"kind":"Field","name":{"kind":"Name","value":"tags"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"synonymsAntonyms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"relatedMeaningId"}},{"kind":"Field","name":{"kind":"Name","value":"relationType"}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetWordQuery, GetWordQueryVariables>;
export const GetInboxItemsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetInboxItems"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"inboxItems"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"sourceContext"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<GetInboxItemsQuery, GetInboxItemsQueryVariables>;
export const SuggestDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Suggest"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"query"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"suggest"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"query"},"value":{"kind":"Variable","name":{"kind":"Name","value":"query"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"transcription"}},{"kind":"Field","name":{"kind":"Name","value":"translations"}},{"kind":"Field","name":{"kind":"Name","value":"origin"}},{"kind":"Field","name":{"kind":"Name","value":"existingWordId"}}]}}]}}]} as unknown as DocumentNode<SuggestQuery, SuggestQueryVariables>;
export const GetStudyQueueDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetStudyQueue"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"limit"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"studyQueue"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"limit"},"value":{"kind":"Variable","name":{"kind":"Name","value":"limit"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"wordId"}},{"kind":"Field","name":{"kind":"Name","value":"partOfSpeech"}},{"kind":"Field","name":{"kind":"Name","value":"definitionEn"}},{"kind":"Field","name":{"kind":"Name","value":"translationRu"}},{"kind":"Field","name":{"kind":"Name","value":"cefrLevel"}},{"kind":"Field","name":{"kind":"Name","value":"imageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"nextReviewAt"}},{"kind":"Field","name":{"kind":"Name","value":"reviewCount"}},{"kind":"Field","name":{"kind":"Name","value":"examples"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceEn"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceRu"}},{"kind":"Field","name":{"kind":"Name","value":"sourceName"}}]}},{"kind":"Field","name":{"kind":"Name","value":"tags"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"synonymsAntonyms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"relatedMeaningId"}},{"kind":"Field","name":{"kind":"Name","value":"relationType"}}]}}]}}]}}]} as unknown as DocumentNode<GetStudyQueueQuery, GetStudyQueueQueryVariables>;
export const GetStatsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetStats"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"stats"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"totalWords"}},{"kind":"Field","name":{"kind":"Name","value":"masteredCount"}},{"kind":"Field","name":{"kind":"Name","value":"learningCount"}},{"kind":"Field","name":{"kind":"Name","value":"dueForReviewCount"}}]}}]}}]} as unknown as DocumentNode<GetStatsQuery, GetStatsQueryVariables>;
export const CreateWordDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateWord"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateWordInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createWord"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"word"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"transcription"}},{"kind":"Field","name":{"kind":"Name","value":"audioUrl"}},{"kind":"Field","name":{"kind":"Name","value":"frequencyRank"}},{"kind":"Field","name":{"kind":"Name","value":"forms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"formText"}},{"kind":"Field","name":{"kind":"Name","value":"formType"}}]}},{"kind":"Field","name":{"kind":"Name","value":"meanings"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"partOfSpeech"}},{"kind":"Field","name":{"kind":"Name","value":"definitionEn"}},{"kind":"Field","name":{"kind":"Name","value":"translationRu"}},{"kind":"Field","name":{"kind":"Name","value":"cefrLevel"}},{"kind":"Field","name":{"kind":"Name","value":"imageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"nextReviewAt"}},{"kind":"Field","name":{"kind":"Name","value":"reviewCount"}},{"kind":"Field","name":{"kind":"Name","value":"examples"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceEn"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceRu"}},{"kind":"Field","name":{"kind":"Name","value":"sourceName"}}]}},{"kind":"Field","name":{"kind":"Name","value":"tags"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"synonymsAntonyms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"relatedMeaningId"}},{"kind":"Field","name":{"kind":"Name","value":"relationType"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<CreateWordMutation, CreateWordMutationVariables>;
export const UpdateWordDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateWord"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateWordInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateWord"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"word"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"transcription"}},{"kind":"Field","name":{"kind":"Name","value":"audioUrl"}},{"kind":"Field","name":{"kind":"Name","value":"frequencyRank"}},{"kind":"Field","name":{"kind":"Name","value":"forms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"formText"}},{"kind":"Field","name":{"kind":"Name","value":"formType"}}]}},{"kind":"Field","name":{"kind":"Name","value":"meanings"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"partOfSpeech"}},{"kind":"Field","name":{"kind":"Name","value":"definitionEn"}},{"kind":"Field","name":{"kind":"Name","value":"translationRu"}},{"kind":"Field","name":{"kind":"Name","value":"cefrLevel"}},{"kind":"Field","name":{"kind":"Name","value":"imageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"nextReviewAt"}},{"kind":"Field","name":{"kind":"Name","value":"reviewCount"}},{"kind":"Field","name":{"kind":"Name","value":"examples"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceEn"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceRu"}},{"kind":"Field","name":{"kind":"Name","value":"sourceName"}}]}},{"kind":"Field","name":{"kind":"Name","value":"tags"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"synonymsAntonyms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"relatedMeaningId"}},{"kind":"Field","name":{"kind":"Name","value":"relationType"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<UpdateWordMutation, UpdateWordMutationVariables>;
export const DeleteWordDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteWord"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteWord"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}]}]}}]} as unknown as DocumentNode<DeleteWordMutation, DeleteWordMutationVariables>;
export const AddToInboxDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"AddToInbox"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"text"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"sourceContext"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"addToInbox"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"text"},"value":{"kind":"Variable","name":{"kind":"Name","value":"text"}}},{"kind":"Argument","name":{"kind":"Name","value":"sourceContext"},"value":{"kind":"Variable","name":{"kind":"Name","value":"sourceContext"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"sourceContext"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<AddToInboxMutation, AddToInboxMutationVariables>;
export const DeleteInboxItemDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteInboxItem"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteInboxItem"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}]}]}}]} as unknown as DocumentNode<DeleteInboxItemMutation, DeleteInboxItemMutationVariables>;
export const ConvertInboxItemDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"ConvertInboxItem"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"inboxId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateWordInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"convertInboxItem"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"inboxId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"inboxId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"word"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"transcription"}},{"kind":"Field","name":{"kind":"Name","value":"audioUrl"}},{"kind":"Field","name":{"kind":"Name","value":"frequencyRank"}},{"kind":"Field","name":{"kind":"Name","value":"forms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"formText"}},{"kind":"Field","name":{"kind":"Name","value":"formType"}}]}},{"kind":"Field","name":{"kind":"Name","value":"meanings"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"partOfSpeech"}},{"kind":"Field","name":{"kind":"Name","value":"definitionEn"}},{"kind":"Field","name":{"kind":"Name","value":"translationRu"}},{"kind":"Field","name":{"kind":"Name","value":"cefrLevel"}},{"kind":"Field","name":{"kind":"Name","value":"imageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"nextReviewAt"}},{"kind":"Field","name":{"kind":"Name","value":"reviewCount"}},{"kind":"Field","name":{"kind":"Name","value":"examples"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceEn"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceRu"}},{"kind":"Field","name":{"kind":"Name","value":"sourceName"}}]}},{"kind":"Field","name":{"kind":"Name","value":"tags"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"synonymsAntonyms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"relatedMeaningId"}},{"kind":"Field","name":{"kind":"Name","value":"relationType"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<ConvertInboxItemMutation, ConvertInboxItemMutationVariables>;
export const ReviewMeaningDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"ReviewMeaning"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"meaningId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"grade"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"reviewMeaning"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"meaningId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"meaningId"}}},{"kind":"Argument","name":{"kind":"Name","value":"grade"},"value":{"kind":"Variable","name":{"kind":"Name","value":"grade"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"wordId"}},{"kind":"Field","name":{"kind":"Name","value":"partOfSpeech"}},{"kind":"Field","name":{"kind":"Name","value":"definitionEn"}},{"kind":"Field","name":{"kind":"Name","value":"translationRu"}},{"kind":"Field","name":{"kind":"Name","value":"cefrLevel"}},{"kind":"Field","name":{"kind":"Name","value":"imageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"nextReviewAt"}},{"kind":"Field","name":{"kind":"Name","value":"reviewCount"}},{"kind":"Field","name":{"kind":"Name","value":"examples"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceEn"}},{"kind":"Field","name":{"kind":"Name","value":"sentenceRu"}},{"kind":"Field","name":{"kind":"Name","value":"sourceName"}}]}},{"kind":"Field","name":{"kind":"Name","value":"tags"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"synonymsAntonyms"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"relatedMeaningId"}},{"kind":"Field","name":{"kind":"Name","value":"relationType"}}]}}]}}]}}]} as unknown as DocumentNode<ReviewMeaningMutation, ReviewMeaningMutationVariables>;