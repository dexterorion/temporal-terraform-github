resource "github_branch" "branch" {
  repository  = var.repository
  branch      = var.branch 
}
