import markdownit from 'markdown-it';
import hljs from '@highlightjs/cdn-assets/es/highlight.js';
import '@highlightjs/cdn-assets/styles/github-dark.css';

const md = markdownit({
  html: true,
  breaks: true,
  highlight: function (str, lang): string {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return (
          '<pre><code class="hljs">' +
          hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
          '</code></pre>'
        );
      } catch (_) {
        console.log(_);
      }
    }

    return (
      '<pre><code class="hljs">' + md.utils.escapeHtml(str) + '</code></pre>'
    );
  },
});

export const markdownRenderer = () => md;
