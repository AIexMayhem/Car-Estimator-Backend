package domain

type Page struct {
    Number int
    Size   int
}

type SortBy int

const (
    SortUnspecified SortBy = iota
    SortDateDesc
    SortDateAsc
    SortPriceDesc
    SortPriceAsc
    SortMileageDesc
    SortMileageAsc
)

type PageResult struct {
    Listings    []Listing
    TotalItems  int64
    TotalPages  int
    CurrentPage int
}