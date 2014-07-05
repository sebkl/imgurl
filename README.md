imgurl
======

Package imgurl offers methods to fetch, thumbnail and convert remote images to a base64 encoded data URL.

#Usage

```
resp, err := Urlify("https://raw.githubusercontent.com/sebkl/globejs/master/screenshots/sample_plain.png",100,100)
fmt.Println(err,resp)
```
