FROM golang:1.17-alpine as dev
# development build

WORKDIR /build

COPY . .
RUN go mod download

COPY *.go ./

RUN go build -o dn-check

CMD [ "/dn-check" ]

FROM golang:1.17-alpine as prod
# production build

WORKDIR /
COPY --from=dev /build/dn-check /dn-check

ENTRYPOINT [ "dn-check" ]

# EOF
