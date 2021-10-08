---
subcategory: "Getting Started"
page_title: "Getting Started"
description: |-
	Terraform ConfiguratioTerraform Configuration

# Getting Started

## Install the provider

```sh
make install
```

## Requiring Providers

Create a terraform file that specifies the bunnycdn provider as required:

```terraform
terraform {
  required_providers {
    bunny = {
      source = "registry.terraform.io/simplesurance/bunny"
    }
  }
}
```
