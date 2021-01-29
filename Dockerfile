FROM golang AS builder
WORKDIR /go/src/github.com/siegerts/drip
COPY . ./
RUN go build -o drip
RUN go install .

FROM rstudio/plumber  
# get drip binary from previous stage
COPY --from=builder /go/bin/drip /usr/local/bin
RUN apt-get update -y && apt-get install tini

# update to use latest plumber R package
RUN R -e "install.packages('devtools')"
RUN R -e "devtools::install_github('rstudio/plumber')"

ENTRYPOINT ["tini", "--"]
ENV PLUMBER_DIR=/usr/local/lib/R/site-library/plumber/plumber/12-entrypoint
CMD ["bash", "-c", "cd ${PLUMBER_DIR} && drip watch --routes --showHost"]

