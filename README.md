# kontena-git-cli

[![Build Status](https://travis-ci.org/jakubknejzlik/kontena-git-cli.svg?branch=master)](https://travis-ci.org/jakubknejzlik/kontena-git-cli)

Manage your kontena cluster in git.

## What's this for

This cli tool helps to organize kontena stacks using git, so you can have all configuration in repository and deploy directly from it.

As it's not recommended to store secrets with project sources this cli assumes you have at least two repositories:

- `grid repository` - store list of stacks, registries, certificates and core configuration (see below)
- `project repository(ies)` - project sources and deployment configuration

### Environment variables

This script communicates with kontena master. Ensure, that you have following environment variables correctly set:

- `KONTENA_MASTER_URL` - api url of kontena master
- `KONTENA_TOKEN` - kontena master token (see `kontena master token create/current`)
- `KONTENA_CLEAR_CERTIFICATES_OFFSET` - number of days after certificates get regenerated (default: `70` - 20 days before expiration)

# Grid repository

Structure of this repository should be:

```
.
|
|--- stacks
|      |--- stack1
|      |      └--- secrets.yml  // file with secrets specific for stack1
|      └--- stack2
|             |--- secrets.yml  // file with secrets specific for stack2
|             └--- kontena.yml  // you can also specify deployment configuration directly
|--- certificates  // add your own certificates from 3rd party providers
|      |--- *.example.com  // full certificate in pem format
|--- certificates.yml // list of certificates managed by kontena (see certificates format)
|--- registries.yml // list of registries (see registry format)
└--- kontena.yml // core stack configuration (see core stack)
```

### Certificates format:

```
www.example.com: #subject name
  alternative_names:
    - test.example.com
    - test2.example.com

doe.com:
  alternative_names:
    - john.doe.com
...
```

_NOTE: If you want to add alternative name/s to already existing certificate, use it as subject name and previous subject move to alternative_names. This will trigger the update (the yaml key has to change as it's used to verify if certificate already exists – for example if you want to create jane.doe.com, add it instead of doe.com and move doe.com to alternative_names)._

### Registry format:

```
registry.example.com:
  username: ...
  email: ... // this is deprecated, but necessary due to older docker version in kontena
  password: ...
```

## Core stack

By default this tool creates stack named `core` with load balancer service for each node which you can link your publicly facing services. Core stack can be modified by providing `kontena.yml` file in repository root. Default configuration:

```
stack: core
services:
  internet_lb:
    image: kontena/lb:latest
    ports:
      - 80:80
      - 443:443
    deploy:
      strategy: daemon
```

**IMPORTANT: this cli tool expects existing service named `core/internet_lb` managing certificates (each certificate is assigned to this load balancer). If you want to use different name for this service/stack, unexpected failures may occur**

## CI setup

Example pipeline configuration file (Gitlab-ci `.gitlab-ci.yml`):

```
image: jakubknejzlik/kontena-git-cli
stages:
  - install
  - cleanup

# this job updates all stacks (secrets, configuration if provided), also generate newly added certificates
install:
  tags:
    - docker
  stage: install
  only:
    - master
  except:
    - schedules
  script:
    - kontena-git grid install production

# this job cleans up old certificates (or expiring soon) and generate new ones
cleanup:
  tags:
    - docker
  stage: cleanup
  only:
    - schedules
  script:
    - kontena-git grid cleanup production
```

_NOTE: install job is triggered on every push to `master`, but for cleanup you have to create scheduled job in gitlab pipelines (for example once a day)_

# Project repository

Project repository should contain at least one yaml file with configuration (named `kontena.yml` by default). You can find more information about kontena stack files here: [Kontena Stack File Reference](https://kontena.io/docs/using-kontena/stack-file.html) (please note, that not all keys are supported, if you are missing some keys, feel free to open new issue).

Stack file example:

```
stack: example
services:
  testpage:
    image: ksdn117/test-page
    environment:
      - KONTENA_LB_VIRTUAL_HOSTS=blah.example.com
    links:
      - core/internet_lb
```

## CI Deployment

Example pipeline configuration file (Gitlab-ci `.gitlab-ci.yml`):

```
stages:
  - deploy

deploy:
  image: jakubknejzlik/kontena-git-cli
  tags:
    - docker
  stage: deploy
  script:
    - kontena-git stack --grid <konten_grid_name> install
    # you can use custom filename like this: kontena-git stack --filename kontena.develop.yml --grid <konten_grid_name> install
```

This pipeline is just for example purposes, there'll be more stages (test, build, migrate) in real life scenarios.

## Executing custom jobs

You can also execute custom code in using already running service.

Example pipeline for manual action (Gitlab-ci `.gitlab-ci.yml`):

```
image: jakubknejzlik/kontena-git-cli
stages:
  - seed

seed-database:
  image: jakubknejzlik/kontena-git-cli
  tags:
    - docker
  stage: seed
  when: manual
  script:
    - kontena-git service --grid <konten_grid_name> exec example/testpage npm run seed-database
```
