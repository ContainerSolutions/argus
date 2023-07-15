provider "kubernetes" {
  config_path    = "~/.kube/config"
}

resource "kubernetes_manifest" "resource_type1" {
  count=10
  manifest = yamldecode(templatefile("resource.yaml", {name="app-${count.index+1}",class="Application",parent="app-platform"}))
}

resource "kubernetes_manifest" "resource_type2" {
  count=10
  manifest = yamldecode(templatefile("resource.yaml", {name="vm-${count.index+1}",class="VirtualMachine",parent="vm-platform"}))
}

resource "kubernetes_manifest" "resource_type3" {
  count=10
  manifest = yamldecode(templatefile("resource.yaml", {name="rt-${count.index+1}",class="NetworkElement",parent="rt-platform"}))
}

resource "kubernetes_manifest" "random_provider" {
  manifest = yamldecode(file("provider.yaml"))
  }

