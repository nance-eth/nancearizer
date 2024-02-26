# Nancearizer

Summarize Nance proposals and their corresponding Discord threads with an LLM.

## Setup

```bash
git clone https://github.com/nance-eth/nancearizer # clone
cd nancearizer # open the directory

cp .example.env .env # set up environment variables in .env

go mod download # download dependencies
go run . # run
```

## Usage

Nancearizer exposes two endpoints:

| Endpoints                    | Description                                                                                                         |
| ---------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| `GET /proposal/{space}/{id}` | Summarize the proposal in the specified space with the specified ID.                                                |
| `GET /thread/{space}/{id}`   | Summarize the Discord discussion thread corresponding to the proposal in the specified space with the specified ID. |

## Docker Build
```docker build --platform=linux/amd64 -f Dockerfile .```
