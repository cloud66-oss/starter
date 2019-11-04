<img src="http://cdn2-cloud66-com.s3.amazonaws.com/images/oss-sponsorship.png" width=150/>

Cloud 66 Go Library
=======

### Getting Started

    go get github.com/cloud66-oss/cloud66


### Authorization

By default, you can use [Cloud66 Toolbet](http://help.cloud66.com/toolbelt/toolbelt-introduction) token which stores in `~/.cloud66/cx.json`. If the file doesn't exist, you can authorize it yourself

    var (
		tokenFile    string = "YOUR_TOKEN_FILENAME"
		tokenDir     string = "YOUR_TOKEN_DIRECTORY"
	)
    cloud66.Authorize(tokenDir, tokenFile)

Or you can use [Personal Access Token](https://app.cloud66.com/oauth/authorized_applications). Create one on and store it in a file like format below:

    {"AccessToken":"YOUR_TOKEN_GOES_HERE","RefreshToken":"","Expiry":"0001-01-01T00:00:00Z","Extra":null}

### Get Client

	var (
		tokenFile    string = "YOUR_TOKEN_FILENAME"
		tokenDir     string = "YOUR_TOKEN_DIRECTORY"
	)
	client := cloud66.GetClient(tokenDir, tokenFile, "")

### Get Stacks List

	var stacks []cloud66.Stack
	stacks, err := client.StackList()
