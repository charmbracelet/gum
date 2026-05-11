FROM gcr.io/distroless/static
COPY gum /usr/local/bin/gum
ENTRYPOINT [ "/usr/local/bin/gum" ]
