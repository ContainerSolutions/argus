provider "kubernetes" {
  config_path    = "~/.kube/config"
}

resource "kubernetes_manifest" "resource_type3" {
  count=10
  manifest = yamldecode(templatefile("resource.yaml", {name="router-${count.index+1}",class="NetworkElement",parent="rt-platform"}))
}

resource "kubernetes_manifest" "random_provider" {
  manifest = yamldecode(file("provider.yaml"))
  }

