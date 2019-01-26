#  JROUTER

Simple http router with GET, POST, DELETE PUT and PATCH wrappers

This is just in the beginning phase and testing some ideas.

Currently supports parameters in the form of http://localhost:8081/person/@id/@other_param/@etc/



	router = &Router{DebugLog: true}
	router.GET("/person/@firstname/@lastname", router.LogWrapper(router.getPerson))


## RoadMap

- Auth Section, should be plugable


##  History
- Added Parameter count checking
- Added NamedParameter support
- Added Parameter support
- DELETE + PUT + PATH added 



