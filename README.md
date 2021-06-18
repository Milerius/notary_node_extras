# notary_node_extras

## build notary_mining_taxes

```
git clone https://github.com/Milerius/notary_node_extras
go mod download
cd notary_mining_taxe
# edit config.json
go build
./notary_mining_taxe
```

## config.json

- range -> "month", "year", or specific number of month eg: 6 -> June

## TODO

- Specify explicit month in range in config.json
- Handle fiat dynamically (choose the fiat of the output)