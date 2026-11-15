package main

import (
	//"net"
	"net/http"
	//"sync"
	"time"
	"os"
	"fmt"
	//"bufio"
	"strings"
	//"encoding/csv"
	"encoding/json"
	//"io/ioutil"
	//"log"
	"io"
)

type Variant struct {
	Type string;
	URL string;
	Width int;
	Height int;
	FileExt string;
}

type MediaAsset struct {
	ID int;
	Created_At string;
	Updated_At string;
	MD5 string;
	File_Ext string;
	File_Size int;
	Image_Width int;
	Image_Height int;
	Duration float64;
	Status string;
	File_Key string;
	Is_Public bool;
	Pixel_Hash string;
	Variants []Variant;
}

type Post struct {
	ID int;
	Created_At string;
	Uploader_ID int;
	Score int;
	Source string;
	MD5 string;
	Last_Comment_Bumped_At string;
	Rating string;
	Image_Width int;
	Image_Height int;
	Tag_String string;
	Fav_Count int;
	File_Ext string;
	Last_Noted_At string;
	Parent_ID int;
	Has_Children bool;
	Approver_ID int;
	Tag_Count_General int;
	Tag_Count_Artist int;
	Tag_Count_Character int;
	Tag_Count_Copyright int;
	File_Size int;
	Up_Score int;
	Down_Score int;
	Is_Pending bool;
	Is_Flagged bool;
	Is_Deleted bool;
	Tag_Count int;
	Updated_At string;
	Is_Banned bool;
	Pixiv_ID int;
	Last_Commented_At string;
	Has_Active_Children bool;
	Bit_Flags int;
	Tag_Count_Meta int;
	Has_Large bool;
	Has_Visible_Children bool;
	Media_Asset MediaAsset;
	Tag_String_General string;
	Tag_String_Character string;
	Tag_String_Copyright string;
	Tag_String_Artist string;
	Tag_String_Meta string;
	File_URL string;
	Large_File_URL string;
	Preview_File_URL string;
}

//Do not edit these variables directly; edit credentials.txt instead
var username = "";
var apiKey = "";
var agent = "BatchDLBot/1.0 (user #560996)"; //DO NOT CHANGE THIS UNLESS FORKING OR UPDATING THE BOT

//Do not edit these variables directly; edit arguments.txt instead
var filepath = ""
var tags = "";
var limit = 200; //Posts per page limit; do not change unless you know what you're doing

var client = http.DefaultClient;

func RequestWithAgent(req string, url string) *http.Request {
	//Creates an HTTP request with the custom defined globally
	//using the req and url parameters
	r, err := http.NewRequest(req, url, nil);
	if (err != nil) {
		panic(err);
	}
	r.Header.Set("User-Agent", agent);

	return r;
}

func DownloadEntireSearch(data []Post, page int) bool {
	//This function iterates over an array of posts, named data, and
	//attempts to download every post in the array by
	//getting their URLs. Note that not all posts in the array
	//may have visible URLs.

	//The page number, page, is simply to keep track of how many posts
	//have been iterated through so far and calculate what post
	//we're on accordingly. It has no effect on the actual download.

	//If no posts are in the array, the function returns true to indicate
	//that the end of the search was reached.

	var eos bool; eos = true;
	for i, post := range data{
		if i >= limit {
			break;
		}
		eos = false;

		if (post.Is_Deleted) {
			fmt.Printf("Skipping over deleted post (URL %d)\n", (page-1)*limit+i);
			continue;
		}

		fmt.Printf("Downloading URL %d: %s\n", ((page-1)*limit)+i, post.File_URL);
		//TODO: Format the filename in the same manner as when downloaded from the browser
		filename := fmt.Sprintf("%d.%s", post.ID, post.File_Ext);

		fileUrl := strings.Split(post.File_URL, "//");

		//hacky workaround to get around attempting to iterate over an empty url
		if len(fileUrl) < 2 {
			fmt.Println("No URL found (post may be hidden)");
			continue;
		}

		pre := fmt.Sprintf("https://%s:%s@", username, apiKey);
		dlFlag := "?download=1";
		dlFlag = "";

		rawUrl := fmt.Sprintf("%s%s%s", pre, fileUrl[1], dlFlag);
		//rawUrl := post.File_URL;

		//Sometimes the script activates Cloudflare's protection, resulting in an
		//HTML waiting page to be returned. We try to catch this, wait, and reload
		//until we get the proper page

		var dlSuccess bool;
		dlSuccess = false;

		var bodyCopy []byte;

		for !dlSuccess {
			r := RequestWithAgent("GET", rawUrl);

			content, err := client.Do(r);
			if (err != nil) {
				panic(err);
			}
			defer content.Body.Close();

			bodyCopy, err = io.ReadAll(content.Body);

			header := string(bodyCopy)[0:15];

			//If the file starts with the DOCTYPE html header, then we know we hit
			//the HTML waiting page
			if (header != "<!DOCTYPE html>") {
				dlSuccess = true;
			} else {
				fmt.Println("Download returned an HTML waiting page, trying again shortly");
				time.Sleep(1*time.Second);
			}
		}

		fullpath := fmt.Sprintf("%s%s", filepath, filename);

		out, err := os.Create(fullpath);
		if err != nil {
			fmt.Println("Folder does not exist or another error occurred. Attempting to create folder...");
			err2 := os.MkdirAll(filepath, os.ModePerm);
			if err2 != nil {
				panic(err2);
			}
			out, err2 = os.Create(fullpath);
			if err2 != nil {
				panic(err2);
			}
		}
		defer out.Close();

		_, errr := io.Writer.Write(out, bodyCopy);
		if errr != nil {
			panic(errr);
		}
		//time.Sleep(1*time.Second);
	}
	return eos;
}

func main() {

	credentials, err := os.ReadFile("credentials.txt"); //do not change this
	if err != nil {
		println("credentials.txt file nonexistent or cannot be read. Please make sure it exists and is accessible.")
		panic(1);
	}

	cParams := strings.Split(string(credentials), "\n");
	username = cParams[0];
	apiKey = cParams[1];
	//userID := cParams[2];

	args, err := os.ReadFile("arguments.txt"); //do not change this
	if err != nil {
		println("arguments.txt file nonexistent or cannot be read. Please make sure it exists and is accessible.")
		panic(1);
	}

	argsSeparated := strings.Split(string(args), "\n");
	filepath = argsSeparated[0];
	arrTags := strings.Split(argsSeparated[1], " ");
	i := 0;
	for (i < len(arrTags)-1) {
		tags = fmt.Sprintf("%s%s+", tags, arrTags[i]);
		i++;
	}
	tags = fmt.Sprintf("%s%s", tags, arrTags[i]);

	//agent = fmt.Sprintf(agent, userID);

	url := fmt.Sprintf("https://%s:%s@danbooru.donmai.us/", username, apiKey);

	page := 1;

	req := fmt.Sprintf("posts.json?tags=%s&limit=%d&page=%d", tags, limit, page);

	full := fmt.Sprintf("%s%s", url, req);

	r := RequestWithAgent("GET", full);

	var eos bool; eos = false;
	for !eos {
		req = fmt.Sprintf("posts.json?tags=%s&limit=%d&page=%d", tags, limit, page);
		full = fmt.Sprintf("%s%s", url, req);
		r = RequestWithAgent("GET", full);
		res, err := client.Do(r);
		if (err != nil) {
			panic(err);
		}
		defer res.Body.Close();

		decoder := json.NewDecoder(res.Body);
		var data []Post;
		err = decoder.Decode(&data);
		if err != nil {
			fmt.Printf("%T\n%s\n%#v\n",err, err, err)
			switch v := err.(type){
				case *json.SyntaxError:
					fmt.Println(v);
			}
		}
		//fmt.Printf("ID: %d\n", data.ID);
		//fmt.Printf("Preview URL: %s\n", data.Preview_File_URL);

		//fmt.Printf("%+v\n", data);

		eos = DownloadEntireSearch(data, page);

		page += 1;

		time.Sleep(1*time.Second);
	}


}
