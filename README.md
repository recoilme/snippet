## Simple library for scrapping title/descriptions

Supported:

 - timeout
 - converting to utf8
 - custom headers
 - meta tags

## Usage
```
	snippet, err := snippet("http://vc.ru", 0, nil)
	fmt.Println("Title:", snippet.Title, "\nDescription:", snippet.Description)

    // use https://github.com/pokatovski/text-cleaner if clean needed
```

## Author

t.me/recoilme