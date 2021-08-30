FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY notification-webhook /bin/notification-webhook

CMD ["/bin/notification-webhook", "-c", "/app/notification-webhook.ini"]
