---
subcategory: "Getting Started"
page_title: "Example Usage"
description: |-
	Example of creating a Pull Zone

# Example Usage

## Basic Pull Zone Blocks

```terraform
resource "bunny_pullzone" "pullzone-terraform" {
  name       = "pz-terraform"
  origin_url = "https://terraform.io"
}
```
