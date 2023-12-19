FROM amd64/alpine
EXPOSE 3001
WORKDIR /app
COPY auth_service_go .
CMD ["/app/auth_service_go"]
