## 0.3.0 (Unreleased)

FEATURES:

* **New Resource** `hostname`


IMPROVEMENTS:

* errors: added additional context to error messages of pullzones and edgerules
* resource/pullzone: new attribute: `cache_error_responses`
* resource/pullzone: new block `safehop`

## 0.2.0 (November 12, 2021)

FEATURES:

* **New Resource** `edgerule`

IMPROVEMENTS:

* logging: HTTP-Responses received from the bunny.net API are logged with debug
           log level

BUG FIXES:

* when an unsuccessful HTTP API response with an empty body was received, the
  error message stated wrongly that a JSON parsing error occurred while
  processing the body

## 0.1.0 (November 04, 2021)

FEATURES:

* **New Resource** `pullzone`
