# DanbooruBatchDL
Danbooru batch downloader program that mass downloads posts according to what tags are being searched for.

Created by letcreate (ToonLucas22 on GitHub)

This program was born out of intense dissatisfaction with all the existing booru downloader programs I could find when I wanted to do a mass batch download (i.e. none of them worked). With my recently attained knowledge about putting together web clients, I decided to a mass batch downloader myself instead.

## Prerequisites
This program requires Go version 1.24.4 (later versions may also work).

## How does it work?

I should preface this by saying that this program requires having a Danbooru user account. Attempting to use Danbooru's API endpoints anonymously *will* result in a 403 error, which was the primary reason why every other program I could find for this purpose failed. There is no way (that I know of) to perform an automated mass batch download without being logged in.

This program expects the user credentials in a file named credentials.txt in the program's working directory. This is a 2-line text file where you store, in the following order:

* Your Danbooru username
* Your Danbooru API key

If you don't have an API key, you can create one by going to your profile, clicking on "View" next to "API Key", click on the "Add" button, and create your API key from there.

This program *only* uses your API key to authenticate when indexing through the posts. It does *not* remotely grab your API key (ala IP grabbers, etc.). Still, you might want to select "posts:index" as the only granted permission in your API key; it is the only relevant permission required for the program to function.

Additionally, the program expects the arguments to be passed through the arguments.txt file, also in the program's working directory, like such:

* Path to the download directory
* Search tags to apply

It should be noted that it's technically possible to leave the tags field blank, but doing so will result in the program to attempt to download *the entirety of Danbooru's pool,* which is heavily not recommended due to the sheer amount of images and the gargantuan amount of disk space they will occupy.

The tags should be space-separated (i.e. "1girl solo", without quotation marks).

If the specified path doesn't exist, the program will create it at runtime.

# How do I run it?

To run the program, open the terminal and type:

```go run danbooru.go```

The program will then try to mass-download all available posts fitting the search tags to the specified directory. Some posts may not be available to download (i.e. hidden), such as those requiring a Gold account.

# Was this vibecoded?

This program was written entirely with human hands. No generative AI was used whatsoever in the creation of this program. Any bug reports or pull requests suspected of being AI-generated or AI-assisted will be denied outright.

I also ask, as per the LICENSE.md, that this program *not* be fed into any sort of LLM, nor any content downloaded through it be fed into one. Thank you.
