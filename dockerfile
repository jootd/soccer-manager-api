# Build the Go Binary.
FROM golang:1.24.4 as build_sales-api

ENV CGO_ENABLED 0
ARG BUILD_REF

# Create the service directory and the copy the module files first and then
# download the dependencies. If this doesn't change, we won't need to do this
# again in future builds.
RUN mkdir /service
COPY go.* /service/
WORKDIR /service
RUN go mod download

# Copy the source code into the container.
COPY . /service


# Build the service binary.
WORKDIR /service/app/services/sales-api
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# Build sales-admin binary
WORKDIR /service/app/services/tooling/sales-admin
RUN go build -o sales-admin


# Run the Go Binary in Alpine.
FROM alpine:3.17
ARG BUILD_DATE
ARG BUILD_REF
RUN addgroup -g 1000 -S sales && \
    adduser -u 1000 -h /service -G sales -S sales
COPY --from=build_sales-api --chown=sales:sales /service/app/services/tooling/sales-admin/sales-admin /service/sales-admin
COPY --from=build_sales-api --chown=sales:sales /service/app/services/sales-api/sales-api /service/sales-api
WORKDIR /service
USER sales
CMD ["./sales-api"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="sales-api" \
      org.opencontainers.image.authors="jootd " \
      org.opencontainers.image.source="https://github.com/jootd/manager-api/" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="jootd"
