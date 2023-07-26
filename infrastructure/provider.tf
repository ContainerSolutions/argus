provider "google" {
  project     = "${var.project_id}"
  region      = "us-central1"
}

variable "project_id" {
    type = string
}
variable "password" {
    type = string
}