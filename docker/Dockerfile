FROM alpine
RUN apk update && apk upgrade && \
    apk add --no-cache bash
RUN apk add --no-cache ca-certificates
RUN mkdir -p /sendmailserviceproxy
WORKDIR /sendmailserviceproxy
COPY ./sendmailserviceproxy-server /sendmailserviceproxy/sendmailserviceproxy-server
RUN chmod +x /sendmailserviceproxy/sendmailserviceproxy-server
EXPOSE 80
CMD /sendmailserviceproxy/sendmailserviceproxy-server --port=80 --host=0.0.0.0