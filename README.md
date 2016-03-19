# Packer Post-Processor to Export Terraform Variables

Building [Amazon Machine Images](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) (AMIs)
using [Packer](https://www.packer.io) with the [Amazon builders](https://www.packer.io/docs/builders/amazon.html) 
is immensely useful, but it can be hard to feed the produced AMI IDs into [Terraform](https://www.terraform.io/), 
or other automation tools. This post-processor exports the producted artifact information into a JSON file
which can be fed into a Terraform plan as a var file, or into other tools. The original artifact is also
passed through for further post-processing.

## Usage

Compile and install the binary `terraform-post-processor-tfvars`, and add the
post-processor to your packer template, something like:

```
{
  "builders": [
    {
      "name": "my-image",
      "type": "amazon-ebs",
      "region": "us-east-1",
      ...
    }
  ],
  "post-processors": [
    {
      "type": "tfvars",
      "output": "packer.tfvars"
    }
  ]
}
```

Run your packer build. It should create a `packer.tfvars` file with JSON
describing the built AMIs. Then you can use this in your terraform config:

```
aws_instance "my-instance" {
  ami = "${var.packer-us-east-1-my-image}"
}
```

## TODO

 - [ ] Remote support, S3, like remote state
 - [ ] Existing files: replace/update/ignore
 - [ ] Customisable variables names and values
 - [ ] HCL support?
