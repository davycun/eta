FROM golang:1.24.11-bookworm

# runtime dependencies
RUN set -eux; \
    mkdir -p -m 0700 ~/.ssh; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
        wget \
        curl \
        unzip \
        tar \
        git \
	; \
	rm -rf /var/lib/apt/lists/*

WORKDIR /eta

EXPOSE 8080

CMD ["tail", "-f", "/dev/null"]
