# gozer #

Gozer aims to be a simple command line utility which searches for keyword in context strings.  Gozer may also be daemonized and run on several machines with a shared key, to perform index searching, possibly useful for log file searching in the cloud.

## Building ##

	$ go get github.com/tools/godep
	$ godep go build

## Running ##

	$ ./gozer -d

Then visit http://localhost:3000/ and search for contents of files in gozer's working directory.

## Status ##

This concept is not fully implemented yet, but may be implemented in three distinct phases:

	1. Simple command line search for occurances of a string in files from a start location, such as /var/log.

		- Currently exploring indexing with index/suffixarray, currently working but needs testing.

	2. Daemonized and continuous search for occurances of strings in files from a particular 
	location, using martini to handle api calls.

		- Somewhat working but needs a lot of clean up work.

	3. Daemonized and distributed indexing for search of occurances of strings in files using 
	martini to handle api calls with peer aggregation of index or search results from different 
	systems using a shared key.

		- Not working on this yet, but will likely test using coreos, docker and vagrant with three 
		machines and some test data.

