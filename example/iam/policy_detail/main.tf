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

data "wangsu_iam_policy_detail" "test-policy" {
  policy_name = "policyName"
}

output "first_policy_name" {
  value = data.wangsu_iam_policy_detail.test-policy
}