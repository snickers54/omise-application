# Setup
Needs Golang 1.11+ for using go modules.  
From the `cmd/song-pah-pa` folder:  

`go mod tidy && go build && ./song-pah-pa ../../assets/fng.1000.csv.rot128`

# Song Pah Pa
I will explain here my train of thoughts while working on this exercise. I think my solution could be improved, but I already spent a lot of time on it, and software engineering is about finding the right balance between time, complexity and solving a problem. 

## ROT 128
The first thing I did was to read the existing code, the existing ROT128 reader felt not right per my experience for multiple reasons:

- I don't need the Writer, only the Reader
- The existing Reader would put everything in-memory, the example files given seems to indicate the number of lines, and I suspect bigger files could be given as input. Therefore be a problem if not enough memory available.
- One of the bonus is to have a minimum memory footprint, therefore using a Scanner and read the file Line by Line would be more efficient than saving everything in-memory and then decipher it. 

For the reasons mentionned above, I decided to rewrite my own Reader (using a Scanner). I had to reimplement the SplitFunc for the Scanner since I'm reading ciphered data and here it's obvious I've to read CR and NL chars that are ROT128. 

## Omise-go Golang library
I used the existing library in an attempt to save time and not reinvent the wheel. That let me focus on the core problem which is concurrency and efficiency.
Not knowing if this library was thread safe or not, I considered it wasn't and created one Client per go routine. 

## Package organization
I followed the [golang-standards](https://github.com/golang-standards/project-layout) way of structuring my project. We can argue it's a good or bad way, but it's a standard and therefore worth using for consistency throughout the community projects. 

## Statistics
Since some currencies can have really big numbers, I'm thinking about THB or VND. I preferred using (Big Integers)[https://golang.org/pkg/math/big/].

Here I made a tradeoff, this Big Integer is a runtime type and therefore add overheads (memory and CPU). But I prefer my software to be slower and consuming more memory than having unexpected behavior if passed really big number of transactions that could cause an integer overflow. 

I used a mutex to make sure the goroutines didn't write at the same time on the Stats object. Avoiding race conditions and inaccurate counting. 

## Dependencies
I used the (go modules)[https://github.com/golang/go/wiki/Modules] which means my project is only compatible with golang 1.11+.

## Concurrency
A naive approach to this, would be to create a go routine per line of the CSV file since go routine are "cheap".

I preferred to go for a consumer / producer (1-n) approach with a pool of limited workers. Because the number of HTTP clients and actual HTTP connections the software can handle has a hard limit because of the number of file descriptor available on the server or laptop running this code. 
Moreover, even if File Descriptors were unlimited, having one go routine per row wouldn't be necessarily faster and could actually be slower than working with a pool of workers. 

## Multi-core CPU
Since 2018 and Golang 1.5+, by default, Go programs run with GOMAXPROCS set to the number of cores available; in prior releases it defaulted to 1.

Therefore, there is no need for my main to contains the famous `runtime.GOMAXPROCS(runtime.NumCPU())`

## Omise.co limit and throttling
I did here a really simple thing, which is telling my go routine to sleep for an arbitrary 5s whenever it gets a `429 Too Many Requests` response. 

This could be obviously done better, but the Omise.co API doesn't provide any insights on how many requests per seconds or minutes per IP or credentials it allows and since it seems to be handled by NGINX and not the API itself, there is no timestamps details to have a dynamic/smart retry. 

## Missing

- More tests
- More statistics
- GoDoc