# Continuos deployment example

In order for this to work we provide the token via an environment variable file called `.envrc`, for this to work copy `.envrc.example` to `.envrc` and update your token, then run direnv allow, make sure you have direnv installed. If for some reason you don't like environment variables you can always define your token in a file with extension `tfvars`.

[Go continuous delivery with Terraform and Kubernetes](https://techsquad.rocks/blog/go_continuous_delivery_with_terraform_and_kubernetes/)
