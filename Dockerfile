FROM alpine:1.1.1

WORKDIR /app
# copy binary
COPY ./bin/app /app/bin/app
# run binary 
CMD [ "./app/bin/app" ]
