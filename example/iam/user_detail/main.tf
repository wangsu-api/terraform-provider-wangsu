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

data "wangsu_iam_user_detail" "user1" {
  login_name = "user1"
}

output "show-user1" {
  value = data.wangsu_iam_user_detail.user1
}