echo "Performing go mod tidy"
go mod tidy

echo "Installing gq"
go build

echo "Copying gq to /usr/bin"
sudo cp gq /usr/local/bin/gq
