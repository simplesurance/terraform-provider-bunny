# Bunny.Net Terraform Provider

[![terraformregistry](https://img.shields.io/badge/terraform-registry-blueviolet)](https://registry.terraform.io/providers/simplesurance/bunny)

This repository provides a [Terraform](https://terraform.io) provider for the
[Bunny.net CDN platform](https://bunny.net/). \
It currently only supports to manage Pull Zones.

## Development

### Using the Local Provider with Terraform

1. Build and install the provider binary via:

    ```sh
    make install
    ```

2. Use the [Development Overrides for Provider
Developers](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers)
feature to enforce using the local `terraform-provider-bunny` binary. \
Run:

    ```sh
    make gen-dev-tftrc
    ```

    to generate a terraform config with `dev_overrides` statement, referencing the
    directory `make install` installed the binary to.

3. Instruct Terraform to use the new config file instead of the default
one by setting the `TF_CLI_CONFIG_FILE` to the path of the generated
`bunny-dev.tftrc` file. For example:

    ```sh
    export TF_CLI_CONFIG_FILE="/home/fho/tf-provider-bunny-dev.tftrc"
    ```

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

### Missing Features

- `terraform import` support
- unsupported Pull Zone features:
  - Certificates
  - `cache_error_response`
  - `enable_query_string_ordering`
- Pull Zone fields with missing write support:
  - `blocked_referrers`
  - `access_control_origin_header_extensions`
  - all `enable_geo_zone_*` fields
- Hostname fields with missing write support:
  - `force_ssl`

## Status

The provider is under initial development and should be considered as
unstable. \
Breaking API changes can happen anytime.
