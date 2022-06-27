resource "bunny_pullzone" "pullzone-terraform" {
  name       = "pz-terraform"
  origin_url = "https://terraform.io"

  safehop {
    enable                          = true
    origin_connect_timeout          = 10
    origin_response_timeout         = 60
    origin_retries                  = 2
    origin_retry_5xx_response       = true
    origin_retry_connection_timeout = false
    origin_retry_delay              = 3
    origin_retry_response_timeout   = true
  }
}

