# Receipt Processor

## How to Run

Assuming you have Docker installed and set up, run the following command to
build the image specified in the provided `Dockerfile`:

```sh
docker build -t receipt-processor .
```

Then run a container based off of this image with this command:

```sh
docker run -d --rm -p 8080:8080 --name receipt-processor-app receipt-processor
```

**Note** while you can choose any external port (the `8080` before the colon (`:`)),
the internal port (the `8080` after the colon (`:`)) should stay as `8080`.

Finally, start the webservice by running this command:

```sh
docker exec -it receipt-processor-app go run main.go
```

From here, you should be able to send requests to
`localhost:8080/receipts/process` and  `localhost:8080/receipts/{id}/points` (or
replace `8080` with whatever external port you chose).
