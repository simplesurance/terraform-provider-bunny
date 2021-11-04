# Bunny.Net Terraform Provider

This repository provides a [Terraform](https://terraform.io) provider for the
[Bunny.net CDN platform](https://bunny.net/). \
It currently only supports to manage Pull Zones.

## Development

### Using the Local Provider with Terraform

Run:
```sh
make install
```

to compile and install the provider binary into your local
`$HOME/.terraform.d/plugins/` directory.

### Running Integration Tests

To run the integration tests a bunny.net account is needed.
The integration tests will create, modify and delete **real** resources.
Therefore a bunny.net account should be used that does not manage resources
used in production.

To run the integration tests set the `BUNNY_API_KEY` to your bunny.net API
key:

```sh
export BUNNY_API_KEY=MY-TOKEN
```

Then run:

```sh
make testacc
```

#### Sweeper

To cleanup resources that might have been left over by running tests, run:

```sh
make sweep
```

### Generating Documentation

```sh
make docs
```

### Creating a Release

1. Ensure the entry for the version in CHANGELOG.md is uptodate. \
   (Keep the `(Unreleased)` marker.)
2. Run:  

    ```sh
    scripts/create-release.sh VERSION
    ``` 

    To finalize the CHANGELOG.md file, create a signed git tag, build the
    release binaries create a GitHub draft release with the binaries.

3. Publish the draft release on github.


## Known Issues
- When destroying a non-existing pull-zone, the operation fails. It should
  succeed and log a warning instead.
- testcases are not validating computed field values

### Missing Features

- `terraform import` support
- Edge Rules
- Custom Hostnames
- Certificates
- Write support is missing for the following Pull Zone fields:
  - `blocked_referrers`
  - `access_control_origin_header_extensions`
  - all `enable_geo_zone_*` fields
- The following Pull Zone fields are unsupported:
  - `cache_error_response`
  - `enable_query_string_ordering`

## Status

The provider is under initial development and should be considered as
unstable. \
Breaking API changes can happen anytime.
