// Runtime shim for @graphql-typed-document-node/core
// This is a placeholder until types are generated
export const gql = (strings: TemplateStringsArray, ...values: any[]) => {
  return strings.reduce((acc, str, i) => acc + str + (values[i] || ''), '');
};


