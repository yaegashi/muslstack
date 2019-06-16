FROM alpine:edge
WORKDIR /src
COPY . /src
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
RUN apk add go git hugo musl-dev
RUN go get github.com/yaegashi/muslstack
ARG STACKSIZE=
RUN if test -n "$STACKSIZE"; then \
        muslstack -s "$STACKSIZE" /usr/bin/hugo; \
    else \
        muslstack /usr/bin/hugo; \
    fi && \
    objdump -p /usr/bin/hugo | grep -A1 STACK && \
    hugo || exit $?