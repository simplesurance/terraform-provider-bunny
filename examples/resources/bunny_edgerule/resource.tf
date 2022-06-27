resource "bunny_pullzone" "mypz" {
  name       = "testpz123aye"
  origin_url = "https://bunny.net"
}

resource "bunny_edgerule" "myer" {
  pull_zone_id          = bunny_pullzone.mypz.id
  action_type           = "block_request"
  trigger_matching_type = "all"
  trigger {
    pattern_matching_type = "any"
    type                  = "random_chance"
    pattern_matches       = ["50"]
  }
}
