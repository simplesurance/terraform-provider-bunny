---
layout: ""
page_title "Provider Bunny.net"
description: |-
	The Bunny.net provider configures Bunny.net resources.
---

# Bunny.net Provider

The Bunny.net provider can be used to configure Bunny.net infrastructure.
Currently it is only supported to configure [Bunny CDN](https://bunny.net/) Pull
Zones.

## Getting Started

Define the bunnycdn provider as a required provider in terraform:

```terraform
terraform {
  required_providers {
    bunny = {
      source = "registry.terraform.io/simplesurance/bunny"
    }
  }
}
```

Create a basic Bunny Pull Zone:

```terraform
resource "bunny_pullzone" "pullzone-terraform" {
  name       = "pz-terraform"
  origin_url = "https://terraform.io"
}
```

## Authentication

The Bunny API Key can be provided via the environment variable `BUNNY_API_KEY`:

```sh
export BUNNY_API_KEY=MY-KEY
terraform plan
```

See [Configuration](guides/configuration.md) for further information.
