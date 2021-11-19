## 0.4.0 (Unreleased)

## 0.3.0 (November 19, 2021)

FEATURES:

* **New Resource** `hostname`

IMPROVEMENTS:

* errors: added additional context to error messages of pullzones and edgerules
* resource/pullzone: new attribute: `cache_error_responses`

BUG FIXES:

* resource/edgerule: the chance that create/update edgerule operations are lost
                     because of concurrency is minimized. The issue is not fixed
                     entirely
                     ([#20](https://github.com/simplesurance/terraform-provider-bunny/issues/20)).

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
