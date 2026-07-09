FROM alpine:3.24
ADD drone-slack woodpecker-slack
CMD ["/woodpecker-slack"]
