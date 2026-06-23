# Docs Verification

Official Originality.ai docs were reachable through the search index but blocked
direct non-browser fetches with JavaScript / verification pages during this run.
The implementation is limited to endpoint shapes visible in search-indexed
official docs snippets and backed by secondary connector docs.

Verified public evidence:

- `https://docs.originality.ai/scan` is the official "How to Run an
  Originality.ai API Scan" page. Search-indexed content shows Version 3,
  base `api.originality.ai/api/v3`, `POST /scan`, and `X-OAI-API-KEY`.
- `https://docs.originality.ai/scan-results-copy-1` is the official "How to
  Fetch Scan Results With the Originality.ai API" page. Search-indexed content
  shows Version 3 and `GET /scan/{id}`.
- `https://docs.originality.ai/api-v2-0-new` / the docs landing page list the
  Version 3 API sections: `POST Scan`, `POST Batch Scan`, `POST Scan Url`,
  `GET Credit Balance`, and `GET Scan Results`.
- Microsoft Learn's Originality.AI connector page documents
  `X-OAI-API-KEY`, AI detection scan input `content`, score fields
  `score.original` / `score.ai`, and `credits_used`.

Deferred until exact paths/request models are directly verified:

- Batch scans
- URL scans
- Credit balance
- Credit usage/payments
- Plagiarism, readability, grammar, factuality, and optimization fields beyond
  what the scan response returns
