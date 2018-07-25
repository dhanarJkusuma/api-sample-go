
## Database Design
![alt text](https://user-images.githubusercontent.com/6735183/43212614-a2204eac-905e-11e8-94b3-f2cd556f1f4f.PNG)

## Explanation
From the diagram above, There are 2 table.
`exchange_rate` and `exchange_rate_data`

Table `exchange_rate` it just table for grouping data from `exchange_rate_data` by Exchange Currency.
`exchange_rate` has relation `one-to-many` to `exchange_rate_data`.
so, one record of `exchange_rate` has many record of `exchange_rate_data`.

From the API, we can get the average data of `exchange_rate_data` for the last 7 days, including date selected.
Because `exchange_rate_data` belongs to `exchange_rate`, we can get that average group by `exchange_rate`
