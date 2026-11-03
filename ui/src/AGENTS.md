This project uses Mantine v9 for components and tanstack router/query as core libraries.

We have the mantine mcp installed, use that for building components. 

Do not just assume the first answer for components, do research about what patterns are standard and what is considered good ux and good looking uis. we don't want our site to look like AI slop.

Encapsulate meaningful parts of the app in components. 
Stick to the color theme, do not hallucinate random colors. If we're missing something let me know. 
## IMPORTANT: 
We will be using this from our phones, so the design needs to be 100% mobile first.

## Product copy and visual density

- Keep the UI quiet. Do not add filler subtitles, explanatory captions, redundant
  status counts, or helper text that merely describes an obvious control or action.
- Every visible label must help someone identify a field or make a decision.
  Keep accessible names even when a visible label is unnecessary.
- Do not put backend implementation notes in ordinary product flows.
- Keep forms focused on the common task. Reveal optional fields progressively.
- Use short, auto-dismissing toasts for routine success feedback, not persistent banners.
- Prefer Mantine components, layout/style props, theme tokens, and the Styles API.
  Add custom CSS only when these cannot reasonably express the behavior or layout.
- For event editing, prefer compact summaries that expand into controls when clicked.
