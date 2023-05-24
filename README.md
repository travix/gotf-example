# Example provider

This is an example provider generated from protoc-gen-gotf.

## Requirements

Install [buf.build] or protoc

To use protoc instead of buf, do the following before running `run.sh`
 - Comment buf commands and uncomment protoc commands in [gen.sh]
 - Copy [gotf.proto] to root of repository

## Demo

Execute `run.sh`, it will

- Generate protobuf, gRPC files
- Generate terraform provider Go code
- Compile and install terraform provider
- Start gRPC server
- `init`, `validate`, `plan` and `apply` terraform from [tfscript] directory

At the end of the run, 1 existing user, group and 1 new user, group will be created using terraform.

## Provider and executors

[provider] package contains the provider implementation, it uses executors that communicate with [example-server]

[buf.build]: https://buf.build/docs/installation/
[gen.sh]: gen.sh
[gotf.proto]: https://github.com/travix/protoc-gen-gotf/blob/main/gotf.proto
[tfscript]: ./tfscript
[provider]: ./provider
[example-server]: ./example-server
