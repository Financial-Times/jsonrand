# jsonrand
A random json testing tool

Generate a random stream of json documents based on a single example document.

# Installation
go get github.com/Financial-Times/jsonrand

# Example

```
jsonrand --template my_example.json --count 10
```

All strings, numbers, arrays and sub-objects will be replaced with random values.  Strings are treated specially if they appear to be UUIDs or RFC3339 dates.  When replacing numbers, it will be in integer if the input was an integer.
