# Docs Verification

Official Originality.ai docs were reachable through the search index but blocked
direct non-browser fetches with JavaScript / verification pages during this run.
The implementation is limited to endpoint shapes visible in search-indexed
official docs snippets and the Microsoft Power Platform Originality connector
Swagger.

Verified public evidence:

- `microsoft/PowerPlatformConnectors` contains
  `independent-publisher-connectors/Originality/apiDefinition.swagger.json`
  with these exact paths:
  - `GET /api/v1/account/credits/balance`
  - `GET /api/v1/account/credits/content_scan_usage`
  - `GET /api/v1/account/credits/payments`
  - `POST /api/v1/scan/ai`
  - `POST /api/v1/scan/url`
- The same Swagger documents the `X-OAI-API-KEY` credential, `content` request
  field for AI detection, `url` request field for URL detection, credit balance
  response field `balance`, usage fields `contentID` / `credits_used` / `date`,
  payment fields `credits` / `price` / `receipt` / `date`, and scan response
  fields including `score.original`, `score.ai`, and `credits_used`.
- `https://docs.originality.ai/scan-results-copy-1` is the official "How to
  Fetch Scan Results With the Originality.ai API" page. Search-indexed content
  shows Version 3 and `GET /scan/{id}`.
- `https://docs.originality.ai/api-v2-0-new` / the docs landing page list the
  Version 3 API sections: `POST Scan`, `POST Batch Scan`, `POST Scan Url`,
  `GET Credit Balance`, and `GET Scan Results`.

Deferred until exact paths/request models are directly verified:

- Batch scans
- Plagiarism, readability, grammar, factuality, and optimization fields beyond
  what the scan response returns
