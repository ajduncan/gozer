# gozer #

Gozer aims to be a simple command line utility which searches for keyword in context strings.  Gozer may also be daemonized and run on several machines with a shared key, to perform index searching, possibly useful for log file searching in the cloud.

## Status ##

This concept is not implemented yet, but will be implemented in three distinct phases:

	1. Simple command line search for occurances of a string in files from a start location, such as /var/log.
	2. Daemonized and continuous search for occurances of strings in files from a particular location, using martini to handle api calls.
	3. Daemonized and distributed indexing for search of occurances of strings in files using martini to handle api calls with peer 
	   aggregation of index or search results from different systems using a shared key.

Other projects have started with search and replace functionality similar to what you'd use with grep and sed.  
