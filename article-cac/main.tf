
terraform {
  required_providers {
    article = {
      source = "article" # This should match the dev_overrides in .terraformrc
    }
  }
}



provider "article" {
  url = "http://localhost:9999"
}

resource "article" "test_article" {
  heading = "My First Article"
  description = "This is an article created by Terraform"
  tags = ["Terraform", "API", "Article"]

}

resource "article" "test2_article" {
  heading = "My Second Article"
  description = "This is an article created by Terraform (second)"
  tags = ["Terraform"]

}
