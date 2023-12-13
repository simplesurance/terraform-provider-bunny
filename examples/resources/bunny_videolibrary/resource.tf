resource "bunny_videolibrary" "myvl" {
  name   = "testvl"
  replication_regions   = ["NY", "BR"]
}
