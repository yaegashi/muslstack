FROM alpine:edge
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
RUN apk add go git hugo musl-dev
RUN go get github.com/yaegashi/muslstack
COPY main.c .
RUN gcc main.c
ARG STACKSIZE=
RUN if test -n "$STACKSIZE"; then \
        muslstack -s "$STACKSIZE" a.out; \
    else \
        muslstack a.out; \
    fi && \
    objdump -p a.out | grep -A1 STACK && \
    ./a.out >log || \
    echo -e "\nExit $?" >>log && tail -5 log && exit 1