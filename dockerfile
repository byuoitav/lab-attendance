# FROM byuoitav/amd64-alpine
# LABEL Brayden Winterton <brayden_winterton@byu.edu>

FROM alpine:3.18

RUN apk update && apk add bash && apk --no-cache add tzdata

ARG NAME
ENV name=${NAME}

#copy binaries
COPY ${name}-bin ${name}-bin 
COPY version.txt version.txt

# copy frontend
COPY autoclave-dist autoclave-dist

ENTRYPOINT [./${name}-bin]