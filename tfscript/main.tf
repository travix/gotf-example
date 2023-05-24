terraform {
  required_providers {
    example = {
      version = "~> 1.0.0"
      source  = "travix.com/providers/example"
    }
  }
}

provider "example" {
  endpoint = "127.0.0.1:50051"
}

resource example_user "test-one" {
  username = "test-one"
  email    = "test-one@example.com"
}

resource example_group "group-one" {
  name  = "group-one"
  users = [
    example_user.test-one.username,
  ]
}

data example_users "all" {
  depends_on = [
    example_user.test-one,
  ]
}

data example_groups "all" {
  depends_on = [
    example_group.group-one,
  ]
}

output users {
  value = data.example_users.all.users
}

output groups {
  value = data.example_groups.all.groups
}
