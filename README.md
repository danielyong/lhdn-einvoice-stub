## Stub for Malaysia LHDN Einvoice 

As of writing this readme, Malaysia's LHDN (entity equivalent to IRS in the US) is implementing E-Invoicing.
However, they cannot seem to get their act together and provide a proper sandbox environment for development.
So to make development easier, this is a webserver that serves dummy data matching LHDN's specifications.
These APIs does not do any validation or authentication or what not, all it does is serving out data for testing.

## Setting up your project

Run docker compose build and then docker compose up... or just go run .

### License
This project is licensed under [GLWTPL](https://github.com/me-shaon/GLWTPL/blob/master/LICENSE)
