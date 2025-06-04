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

data "wangsu_iam_users" "user_list" {
  page_size     = 1
  page_number   = 10
}

output "user_list" {
  value = data.wangsu_iam_users.user_list
}