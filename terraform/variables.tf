variable "domain_name" {
  type = string
}

variable "subdomain_name" {
  type = string
}

variable "discord_public_key" {
  type      = string
  sensitive = true
}

