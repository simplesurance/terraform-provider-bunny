## 0.7.0 (Unreleased)

BREAKING CHANGES:

* resource/edgerule: renamed action type `ignore_quiery_string` to
                     `ignore_query_string`

IMPROVEMENTS:

* resource/edgerule: support action types `set_status_code` and
                     `bypass_perma_cache`
* resource/edgerule: support trigger types `status_code` and
                     `request_method`

## 0.6.0 (Januar 31, 2022)

IMPROVEMENTS:

* resource/hostname: new block `certificate`, for providing custom SSL
                     certificates

BUG FIXES:

* resource/edgerule: changing the `pull_zone_id` of an edgerule resulted in
                     API errors

## 0.5.0 (Dezember 03, 2021)

IMPROVEMENTS:

* resource/hostname: the `force_ssl` attribute is not a computed field anymore
                     and can be set

BUG FIXES:

* resource/hostname: fix: plan was not empty after hostname creation

## 0.4.0 (November 26, 2021)

BREAKING CHANGES:

* resource/pullzone: All header related attributes were moved to a new block
  called `headers`.
* resource/pullzone: all limits related attributes were moved to the block
  `limits`.
* resource/pullzone: all optimizer related attributes were moved to the block
  `optimizer`.
* resource/pullzone: The type of `access_control_origin_header_extensions`
  changed from string set to a comma-separated string.

IMPROVEMENTS:

* resource/pullzone: new block `safehop`
* resource/pullzone: `access_control_origin_header_extensions` is not a computed
                     field anymore and can be set.
* provider: upgrade terraform-plugin-sdk from version 2.8.0 to 2.9.0

BUG FIXES:

* resource/pullzone: removed raw, unused format string specifiers from some
                     error messages

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
