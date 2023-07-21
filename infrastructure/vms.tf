resource "google_service_account" "vm" {
  account_id   = "common-vm-sa"
  display_name = "Service Account"
}

resource "google_compute_instance" "one" {
  name         = "vm-server-1"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }


  network_interface {
    network = "default"

    access_config {
      // Ephemeral public IP
    }
  }

   metadata_startup_script = "apt-get update && apt-get install mysql nginx"

  service_account {
    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    email  = google_service_account.vm.email
    scopes = ["cloud-platform"]
  }
}