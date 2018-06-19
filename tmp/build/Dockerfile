FROM alpine:3.6

RUN adduser -D dispatch-operator
USER dispatch-operator

ADD tmp/_output/bin/dispatch-operator /usr/local/bin/dispatch-operator
