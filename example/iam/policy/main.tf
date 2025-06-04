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

resource "wangsu_iam_policy" "policy" {
  policy_name = "tf-example"
  description = "this is a policy test"
  policy_document = jsonencode([
    {
      "effect" : "allow",
      "action" : [
       "productCode:actionCode"
      ],
      "resource" : [
        "*"
      ]
    }
  ])
}