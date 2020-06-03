package main
import "net/http"
import "golang.org/x/net/html"
import "fmt"
import "regexp"
import "os"
import "io/ioutil"
import "strconv"
import "strings"

func getHref(t html.Token) (ok bool, href string) {
	// Iterate over token attributes until we find an "href"
	for _, a := range t.Attr {
		re := regexp.MustCompile("(\\/download\\/)([0-9]){5}(\\/.*)")
	   if a.Key == "href" {
		//Find the links that are download files
		 match:=re.FindStringSubmatch(a.Val)
        if len(match)!=0{
			
			href = match[0]
			href=href[10:]
			ok = true
		}
	}
}
	return
}

func getMusicLinks(url string,ch chan string){
	resp, _ := http.Get(url)  //Get the response body of url
	
	z := html.NewTokenizer(resp.Body)   //Creating a new html token
  for { 
    tt := z.Next()
    
    switch {
    case tt == html.ErrorToken:
    	return
	case tt == html.StartTagToken:
		
		t := z.Token()
		
		isAnchor := t.Data == "a" //Finding links in the page
		
	
        if !isAnchor {

			continue
		}
			ok, url := getHref(t)
			if !ok {
				continue
			}
			
			ch<-url
		}
		
	}
	resp.Body.Close()

}
	
func download(ok string) error{
names:=strings.Split(ok,"/") //splitting the url to get name and id
   resp, err := http.Get("http://mymp3song.guru/files/download/id/"+names[0])

	if err!=nil{
		return err;
	}
	body, _ := ioutil.ReadAll(resp.Body)
	error := ioutil.WriteFile("../../Music/"+names[1]+".mp3", body,0777)
	if error!=nil{
         return error
	}
	fmt.Println(names[1]+".mp3 downloaded")
	return nil
	} 

func main(){
chIds := make(chan string)
_= os.MkdirAll("Music", 0755)
for i:=1; i<=5;i++{	
go getMusicLinks("http://mymp3song.guru/filelist/3376/special_mp3_songs/new2old/"+strconv.Itoa(i),chIds)
}
for {
	select {
	case ok:=<- chIds:
		
		go download(ok)

	}
}
}






