export function splitAndCapitalizeCamelCase(word: string) {
  return word
    .replace(/([A-Z])/g, ' $1')
    .replace(/^./, (str) => str.toUpperCase());
}
