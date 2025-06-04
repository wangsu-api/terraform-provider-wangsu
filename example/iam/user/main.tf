terraform {
  required_providers {
    wangsu = {
      source = "registry.terraform.io/wangsu-api/wangsu"
    }
  }
}

provider "wangsu" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
}

resource "wangsu_iam_user" "test_user" {
  login_name     = "tf-example"
  display_name   = "user_display_name"
  status         = "1"
  email          = "mail@example.com"
  mobile         = "1111111111"
  console_enable = "1"
}