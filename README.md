# balGPT

balGPT is a web application that predicts the outcome of soccer matches based on historical data. The application scrapes match results from a website and stores them in a PostgreSQL database. It then uses this data to generate predictions for future matches.

It was written almost entirely with ChatGPT 4.

## Building and Running the Software

1. Clone the repository:

```bash
git clone https://github.com/your-username/balGPT.git
cd balGPT
```

1. Install dependencies:

```bash
# Go dependencies
go mod download
```

1. Run the application:

```bash
source scripts/env.sh
scripts/run-postgres.sh
go run main.go
```

The application will be accessible at http://localhost:8080.

This project was built with the help of:

- The TilburgsAns font: https://www.tilburgsans.nl/
- OpenAI's ChatGPT: https://www.openai.com/chat-gpt/
