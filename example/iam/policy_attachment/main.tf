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

resource "wangsu_iam_policy_attachment" "policyAttachment" {
  policy_name = ["policyName1","policyName2"]
  login_name = "subAccountLoginName"
}