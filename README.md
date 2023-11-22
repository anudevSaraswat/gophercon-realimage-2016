# Real Image Challenge 2016

## Getting Started

1. Clone the repo

```
git clone https://github.com/anudevSaraswat/gophercon-realimage-2016.git
```

2. Run the following command

```
go run .
```

## Notes on using this program

1. While adding a distributor we can add include and exclude locations for it. Currently only one include location is supported however you can add multiple exclude locations. If you need to enter multiple exclude locations input should be hyphen separated. For example - If I am creating a distributor for Gujarat and need to exclude 2 cities Ahmedabad and Vadodara my include input will be `GUJARAT` and exclude input will be `AHMEDABAD-VADODARA`.

2. A sub distributor can also be added for which you have to specify a parent distributor and you can also view list of added distributors.

3. When you view distributors, for each distributor it displays the `Locations Included` and `Locations Excluded` in `COUNTRY-STATE-CITY` format.


https://github.com/RealImage/challenge2016