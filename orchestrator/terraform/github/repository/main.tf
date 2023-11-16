resource "github_repository" "repository" {
  name                   = var.name
  description            = var.description

  auto_init   = true
  
  visibility = "public"
}