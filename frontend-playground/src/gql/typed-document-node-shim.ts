// Shim for @graphql-typed-document-node/core
// This package is types-only and has no runtime code
// We export type aliases that match the expected interface
import type { DocumentNode } from 'graphql';

// TypedDocumentNode is just a type alias for DocumentNode
export type TypedDocumentNode<TResult = any, _TVariables = Record<string, any>> = DocumentNode;

// Additional types exported by the package
export type ResultOf<T> = T extends TypedDocumentNode<infer TResult, any> ? TResult : never;
export type DocumentTypeDecoration<TResult, TVariables> = TypedDocumentNode<TResult, TVariables>;

// Re-export DocumentNode type for type assertions
export type { DocumentNode };

