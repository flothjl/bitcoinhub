build:
	go build -o ./app -v

css:
	npx tailwindcss -i css/input.css -o css/output.css

start: build
	air

