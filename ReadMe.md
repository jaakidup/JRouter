#  JROUTER

Simple http router with GET, POST, DELETE PUT and PATCH wrappers

This is just in the beginning phase and testing some ideas.

Currently supports parameters in the form of http://localhost:8081/person/@id/@other_param/@etc/



	router = &Router{DebugLog: true}
	router.GET("/person/@something", router.LogWrapper(router.getPerson))


##  History
- Added Parameter support
- DELETE + PUT + PATH added 



