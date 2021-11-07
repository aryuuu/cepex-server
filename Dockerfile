FROM alpine:3.14

# copy binary
COPY ./app /
# run binary 
CMD [ "./app" ]
