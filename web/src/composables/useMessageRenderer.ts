import markdownit from 'markdown-it';
import hljs from '@highlightjs/cdn-assets/es/highlight.js';
import '@highlightjs/cdn-assets/styles/github-dark.css';

const md = markdownit({
  html: false,
  breaks: true,
  highlight: function (str, lang): string {
    console.log('lang:', hljs.getLanguage(lang));
    if (!(lang && hljs.getLanguage(lang)))
      return `<pre><code class="hljs">${str}</code></pre>`;

    const highlighted = hljs.highlight(lang, str, true);
    return `<pre><code class="hljs">${highlighted.value}</code></pre>`;
  },
});

function sanitizeContent(content: string): string {
  return content
    .replace(/(?:[^\\]|^)((?:\\{2})*)&/g, '$1\\&')
    .replace(/\n(?=\n)/g, '\n&nbsp;');
}

export const useMessageRenderer = () => (content: string | undefined) => {
  if (!content) return '';
  const sanitized = sanitizeContent(content);
  const rendered = md.render(sanitized);
  return rendered;
};
