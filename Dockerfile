FROM golang:1.10.3-stretch AS app-builder

ENV GOBIN /go/bin

ENV POSTGRES_HOST='postgres://inkster:inkster@db/inkster?sslmode=disable'
ENV INKSTER_PORT=8000
ENV INKSTER_STATIC=static
ENV INKSTER_SERVE_STATIC=1
ENV INKSTER_DEBUG=1
ENV INKSTER_SMTP_HOST=dummy
ENV INKSTER_SMTP_LOGIN=dummy
ENV INKSTER_SMTP_ADDR=dummy
ENV INKSTER_SMTP_PASS=dummy
ENV INKSTER_SMTP_PORT=dummy
ENV INKSTER_SECRET_KEY=plschangeme

RUN mkdir /app
RUN mkdir /go/src/github.com
RUN mkdir /go/src/github.com/dominik-zeglen
RUN mkdir /go/src/github.com/dominik-zeglen/inkster
ADD . /go/src/github.com/dominik-zeglen/inkster
WORKDIR /go/src/github.com/dominik-zeglen/inkster
COPY ./app/graphql.html /app

RUN make schema
RUN dep ensure -vendor-only

RUN CGO_ENABLED=0 go build -o /app/main manage.go


FROM node:8.11.4-alpine AS ui-builder

RUN mkdir /app
run mkdir /src
RUN mkdir /src/app
ADD ./panel /src/app
WORKDIR /src/app

RUN npm install
RUN npm run build

FROM alpine
WORKDIR /app
RUN mkdir /app/app
COPY --from=app-builder /app/main /app
COPY --from=app-builder /app/graphql.html /app/app
COPY --from=app-builder /config.toml /app
COPY --from=ui-builder /src/app/build /app/panel/build

CMD ["/app/main", "runserver"]
