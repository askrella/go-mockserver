# Go MockServer

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Docker](https://github.com/askrella/whatsapp-chatgpt/actions/workflows/docker.yml/badge.svg)
![Docker AMD64](https://img.shields.io/badge/docker-amd64-blue)
![Docker ARM64](https://img.shields.io/badge/docker-arm64-green)
![Build](https://img.shields.io/github/actions/workflow/status/askrella/aws-ses-mock/docker.yml?branch=master)

![Askrella](https://avatars.githubusercontent.com/u/77694724?s=100)

A tiny (< 5mb) implementation of a request capturing proxy written in Golang. Used internally at [Askrella](https://askrella.de/)
for mocking external partners in our CI pipelines and sometimes improving compression and loading times for CSV
datasets.

# Why?

We at [Askrella](https://askrella.de/) sometimes encounter projects with minimal own logic rather than gluing
multiple APIs together and process their data. In these cases we had troubles with the response times these
APIs offer (> 3.5s, rare cases even > 1,4min). Since we want to run extensive tests against these providers in our pipeline, we need
some kind of proxy capturing and mocking these providers in our pipeline.

On our search we encountered [MockServer](https://www.mock-server.com/) written in Java, but the functionality was way to extensive
and the images too big for fast iterations in our scenarios (possibly dozens of different APIs on a small CI/CD pipeline).

Of course we don't offer the same functionality as the big brothers, but we kept it simple and our images
minimal (< 10mb)!
The code itself doesn't need many dependencies and only relies on the Go standard library as well as some
curated compression libraries (e.g. brotli).
Our final container is based on [distroless](https://github.com/GoogleContainerTools/distroless).

# :gear: Getting Started

Run the mockserver using docker compose:
```golang
version: '3'

services:
  mock:
    build: .
    environment:
      MOCK_TARGET: https://google.com 
      # MOCK_PORT: 80 # The port MockServer shall bind to
      # MOCK_HOST: 0.0.0.0 # The interface MockServer shall use
      # CACHE_ENABLED: true # Whether responses to requests shall be cached and re-delivered similar requests
      # RECOMPRESS: true # When enabled the upstream (external server) response body will be decompressed and compressed
                         # using gzip level 9 compression. Depending on the upstream, this can lead to extremely reduced
                         # request times.
    ports:
      - "8080:80/tcp"
```

> Note: Google actually doesn't work (404) because the request currently won't be manipulated accordingly.

# :wave: Contributors

<a href="https://github.com/askrella/aws-ses-mock/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=askrella/aws-ses-mock" />
</a>

* [Askrella Software Agency](askrella.de)
    * [Steve](https://github.com/steve-hb) (Maintainer)

Feel free to open a new pull request with changes or create an issue here on GitHub! :)

# :warning: License
Distributed under the MIT License. See LICENSE.txt for more information.

# Bugs?

Since I coded most of the code at the airport, please let us know if you encounter some kind of bug. :)

# :handshake: Contact Us

In case you need professional support, feel free to <a href="mailto:contact@askrella.de">contact us</a>

