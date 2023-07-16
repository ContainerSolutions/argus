resource "google_artifact_registry_repository" "this" {
  location      = "us-central1"
  repository_id = "argus"
  description   = "Argus Docker Image Repository"
  format        = "DOCKER"
}

resource "google_service_account" "this" {
  account_id   = "argus-github-actions-sa"
  display_name = "Github Actions Service Account"
}

resource "google_artifact_registry_repository_iam_member" "this" {
  project = google_artifact_registry_repository.this.project
  location = google_artifact_registry_repository.this.location
  repository = google_artifact_registry_repository.this.name
  role = "roles/artifactregistry.writer"
  member = "serviceAccount:${google_service_account.this.email}"
}

resource "google_service_account_key" "key" {
  service_account_id = google_service_account.this.name
}

resource "google_service_account" "default" {
  account_id   = "gke-cluster-sa"
  display_name = "GKE cluster Service Account"
}

resource "google_artifact_registry_repository_iam_member" "default" {
  project = google_artifact_registry_repository.this.project
  location = google_artifact_registry_repository.this.location
  repository = google_artifact_registry_repository.this.name
  role = "roles/artifactregistry.reader"
  member = "serviceAccount:${google_service_account.default.email}"
}

resource "google_container_cluster" "primary" {
  name     = "my-gke-cluster"
  location = "us-central1"
  remove_default_node_pool = true
  initial_node_count       = 1
}

resource "google_container_node_pool" "primary_preemptible_nodes" {
  name       = "gke-nodepool-1"
  location   = "us-central1"
  cluster    = google_container_cluster.primary.name
  node_count = 3

  node_config {
    machine_type = "e2-medium"

    service_account = google_service_account.default.email
    oauth_scopes    = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}