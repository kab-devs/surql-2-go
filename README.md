# Surql 2 Go
This cli tool will turn .surql files into Go structs

Currently the tool only supports simple structs with basic field types. I am planning on adding more features in the future.

## Usage
See the example below on how to use the tool. Currently it requires a -filename flag to specify the .surql file to use.

`surql2go -filename <path-to-.surql-file>`

If all goes well, you should see a generated_structs file in the same directory as the .surql file, containing all the generated Go structs.

